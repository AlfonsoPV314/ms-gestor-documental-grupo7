// Proveedor aislado para desarrollo. No implementa ni sustituye al grupo 6.
package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"time"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		cliente := &http.Client{Timeout: 2 * time.Second}
		respuesta, err := cliente.Get("http://127.0.0.1:8080/jwks")
		if err != nil {
			os.Exit(1)
		}
		_ = respuesta.Body.Close()
		if respuesta.StatusCode != http.StatusOK {
			os.Exit(1)
		}
		return
	}
	if os.Getenv("API_ENV") != "dev" && os.Getenv("API_ENV") != "test" {
		panic("identidad-dev solo dev/test")
	}
	publica, privada, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /jwks", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"keys": []any{map[string]any{"kty": "OKP", "crv": "Ed25519", "alg": "EdDSA", "use": "sig", "kid": "dev", "x": base64.RawURLEncoding.EncodeToString(publica)}}})
	})
	mux.HandleFunc("POST /identidad/sesion", func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" && origin != os.Getenv("FRONTEND_URL") {
			http.Error(w, "origen inválido", 403)
			return
		}
		var in struct {
			Usuario string `json:"usuario"`
		}
		r.Body = http.MaxBytesReader(w, r.Body, 1024)
		if json.NewDecoder(r.Body).Decode(&in) != nil {
			http.Error(w, "JSON inválido", 400)
			return
		}
		id := "22222222-2222-4222-8222-222222222222"
		permisos := []string{"documentos:leer", "documentos:operar"}
		switch in.Usuario {
		case "operador":
		case "lector":
			permisos = []string{"documentos:leer"}
		case "ajeno":
			id = "33333333-3333-4333-8333-333333333333"
		case "sin-permisos":
			permisos = []string{}
		default:
			http.Error(w, "usuario inválido", 422)
			return
		}
		t := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.MapClaims{"sub": id, "iss": os.Getenv("AUTH_ISSUER"), "aud": "documental", "exp": time.Now().Add(time.Hour).Unix(), "iat": time.Now().Unix(), "permisos": permisos, "rol": in.Usuario})
		t.Header["kid"] = "dev"
		token, e := t.SignedString(privada)
		if e != nil {
			http.Error(w, "firma no disponible", 500)
			return
		}
		http.SetCookie(w, &http.Cookie{Name: "sesion", Value: token, Path: "/", HttpOnly: true, Secure: true, SameSite: http.SameSiteLaxMode, MaxAge: 3600})
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store")
		_ = json.NewEncoder(w).Encode(map[string]any{"token": token, "usuario": in.Usuario})
	})
	s := &http.Server{Addr: ":8080", Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	if err = s.ListenAndServe(); err != nil {
		panic(err)
	}
}
