package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDocumentacionSoloFueraDeProduccion(t *testing.T) {
	for _, ambiente := range []string{"dev", "test", "prod"} {
		t.Run(ambiente, func(t *testing.T) {
			mux := http.NewServeMux()
			Nueva(mux, nil, ambiente, false, nil)
			for _, ruta := range []string{"/api/docs", "/api/docs/", "/api/docs/iniciar.js", "/api/openapi.json"} {
				respuesta := httptest.NewRecorder()
				mux.ServeHTTP(respuesta, httptest.NewRequest("GET", ruta, nil))
				esperado := http.StatusOK
				if ambiente == "prod" {
					esperado = http.StatusNotFound
				}
				if respuesta.Code != esperado {
					t.Fatalf("%s: estado %d; esperado %d", ruta, respuesta.Code, esperado)
				}
				if ruta == "/api/docs" && ambiente != "prod" && (!strings.Contains(respuesta.Body.String(), "Swagger UI") || strings.Contains(respuesta.Body.String(), "https://")) {
					t.Fatal("la documentación debe usar Swagger UI y recursos locales")
				}
			}
		})
	}
}
