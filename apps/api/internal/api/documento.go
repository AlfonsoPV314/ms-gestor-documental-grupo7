package api

import (
	"github.com/PINGESO-CIVIL/ms-gestor-documental/apps/api/internal/service"
	"time"
)

// Documento es el DTO público; el dominio no depende de validaciones HTTP.
type Documento struct {
	ID          string    `json:"id" format:"uuid"`
	Propietario string    `json:"propietario" format:"uuid"`
	Nombre      string    `json:"nombre" minLength:"1" maxLength:"200" example:"informe.pdf"`
	Tipo        string    `json:"tipo" example:"application/pdf"`
	Tamano      int64     `json:"tamano" minimum:"1" maximum:"20971520"`
	Hash        string    `json:"sha256" pattern:"^[a-f0-9]{64}$"`
	Firma       string    `json:"firma"`
	ClaveID     string    `json:"clave_id"`
	CreadoEn    time.Time `json:"creado_en"`
}

func salida(d service.Documento) Documento {
	return Documento{ID: d.ID, Propietario: d.Propietario, Nombre: d.Nombre, Tipo: d.Tipo, Tamano: d.Tamano, Hash: d.Hash, Firma: d.Firma, ClaveID: d.ClaveID, CreadoEn: d.CreadoEn}
}
