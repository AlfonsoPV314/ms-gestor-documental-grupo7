package service

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"
	"unicode"
)

const MaxArchivo = 20 * 1024 * 1024

var ErrNoEncontrado = errors.New("documento no encontrado")
var ErrProhibido = errors.New("permiso insuficiente")
var ErrEntrada = errors.New("archivo inválido")
var ErrGrande = errors.New("archivo demasiado grande")

type Documento struct {
	ID          string
	Propietario string
	Nombre      string
	Tipo        string
	Tamano      int64
	Hash        string
	Firma       string
	ClaveID     string
	CreadoEn    time.Time
	Contenido   []byte
}

type Identidad struct {
	ID       string
	Permisos []string
}
type Almacen interface {
	Crear(context.Context, Documento, string) error
	Obtener(context.Context, string, string) (Documento, error)
	Listar(context.Context, string, int32, int32) ([]Documento, error)
}
type Servicio struct {
	almacen Almacen
	privada ed25519.PrivateKey
	claveID string
}

func Nuevo(a Almacen, clave ed25519.PrivateKey) *Servicio {
	hash := sha256.Sum256(clave.Public().(ed25519.PublicKey))
	return &Servicio{a, clave, hex.EncodeToString(hash[:8])}
}
func (i Identidad) Permite(permiso string) bool {
	for _, p := range i.Permisos {
		if p == permiso {
			return true
		}
	}
	return false
}
func UUID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	b[6] = (b[6] & 15) | 64
	b[8] = (b[8] & 63) | 128
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[:4], b[4:6], b[6:8], b[8:10], b[10:])
}
func mensaje(d Documento) []byte {
	return []byte("documental:v1\n" + d.ID + "\n" + d.Propietario + "\n" + d.Nombre + "\n" + d.Tipo + "\n" + fmt.Sprint(d.Tamano) + "\n" + d.Hash)
}
func (s *Servicio) Cargar(ctx context.Context, i Identidad, nombre string, contenido []byte, requestID string) (Documento, error) {
	if !i.Permite("documentos:operar") {
		return Documento{}, ErrProhibido
	}
	if len(contenido) > MaxArchivo {
		return Documento{}, ErrGrande
	}
	nombre = strings.TrimSpace(nombre)
	if len(contenido) == 0 || len([]rune(nombre)) > 200 || nombre == "" || nombre == "." || nombre == ".." || path.Base(nombre) != nombre || strings.ContainsAny(nombre, "\\\r\n") || strings.ContainsFunc(nombre, unicode.IsControl) {
		return Documento{}, ErrEntrada
	}
	tipo := strings.Split(http.DetectContentType(contenido), ";")[0]
	switch tipo {
	case "application/pdf", "image/png", "image/jpeg", "text/plain":
	default:
		return Documento{}, ErrEntrada
	}
	hash := sha256.Sum256(contenido)
	d := Documento{ID: UUID(), Propietario: i.ID, Nombre: nombre, Tipo: tipo, Tamano: int64(len(contenido)), Hash: hex.EncodeToString(hash[:]), ClaveID: s.claveID, CreadoEn: time.Now().UTC(), Contenido: contenido}
	d.Firma = base64.StdEncoding.EncodeToString(ed25519.Sign(s.privada, mensaje(d)))
	if err := s.almacen.Crear(ctx, d, requestID); err != nil {
		return Documento{}, err
	}
	return d, nil
}
func (s *Servicio) Obtener(ctx context.Context, i Identidad, id string) (Documento, error) {
	if !i.Permite("documentos:leer") {
		return Documento{}, ErrProhibido
	}
	return s.almacen.Obtener(ctx, id, i.ID)
}
func (s *Servicio) Listar(ctx context.Context, i Identidad, limite, offset int32) ([]Documento, error) {
	if !i.Permite("documentos:leer") {
		return nil, ErrProhibido
	}
	return s.almacen.Listar(ctx, i.ID, limite, offset)
}
func (s *Servicio) Verificar(d Documento) bool {
	hash := sha256.Sum256(d.Contenido)
	firma, err := base64.StdEncoding.DecodeString(d.Firma)
	return err == nil && d.ClaveID == s.claveID && hex.EncodeToString(hash[:]) == d.Hash && ed25519.Verify(s.privada.Public().(ed25519.PublicKey), mensaje(d), firma)
}
func (s *Servicio) ClavePublica() (string, string) {
	return s.claveID, base64.StdEncoding.EncodeToString(s.privada.Public().(ed25519.PublicKey))
}
