package api

import (
	"context"
	"github.com/PINGESO-CIVIL/ms-gestor-documental/apps/api/internal/service"
	"log/slog"
	"net/http"
	"regexp"
	"time"
)

type respuesta struct {
	http.ResponseWriter
	estado int
}

func (r *respuesta) WriteHeader(c int) {
	if r.estado == 0 {
		r.estado = c
		r.ResponseWriter.WriteHeader(c)
	}
}
func (r *respuesta) Write(b []byte) (int, error) {
	if r.estado == 0 {
		r.WriteHeader(200)
	}
	return r.ResponseWriter.Write(b)
}

var correlacion = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,80}$`)

func Observar(next http.Handler, version string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		inicio := time.Now()
		id := r.Header.Get("X-Request-Id")
		if !correlacion.MatchString(id) {
			id = service.UUID()
		}
		w.Header().Set("X-Request-Id", id)
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Cache-Control", "no-store")
		ctx, cancel := context.WithTimeout(context.WithValue(r.Context(), peticionContexto, id), 10*time.Second)
		defer cancel()
		rw := &respuesta{ResponseWriter: w}
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
		slog.Info("http", "servicio", "ms-gestor-documental", "version", version, "request_id", id, "metodo", r.Method, "ruta", r.Pattern, "estado", rw.estado, "duracion_ms", time.Since(inicio).Milliseconds())
	})
}
