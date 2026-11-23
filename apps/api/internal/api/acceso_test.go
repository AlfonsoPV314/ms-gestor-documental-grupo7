package api

import (
	"github.com/PINGESO-CIVIL/ms-gestor-documental/apps/api/internal/service"
	"testing"
)

// Matriz provisional: los permisos vienen del JWT firmado; el rol no se confía desde el navegador.
func TestMatrizAcceso(t *testing.T) {
	for _, caso := range []struct {
		rol          string
		permisos     []string
		leer, operar bool
	}{
		{"lector", []string{"documentos:leer"}, true, false},
		{"operador", []string{"documentos:leer", "documentos:operar"}, true, true},
		{"sin-permisos", nil, false, false},
	} {
		t.Run(caso.rol, func(t *testing.T) {
			i := service.Identidad{Permisos: caso.permisos}
			if i.Permite("documentos:leer") != caso.leer || i.Permite("documentos:operar") != caso.operar {
				t.Fatal("matriz de permisos incorrecta")
			}
		})
	}
}
