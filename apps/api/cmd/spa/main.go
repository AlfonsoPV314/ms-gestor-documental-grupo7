package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(200) })
	for ruta, variable := range map[string]string{"/api/": "BACKEND_URL", "/identidad/": "IDENTIDAD_URL"} {
		u, err := url.Parse(os.Getenv(variable))
		if err != nil || u.Host == "" {
			panic("falta " + variable)
		}
		p := httputil.NewSingleHostReverseProxy(u)
		p.Transport = &http.Transport{ResponseHeaderTimeout: 10 * time.Second}
		mux.Handle(ruta, p)
	}
	archivos := http.FileServer(http.Dir("/srv"))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "same-origin")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; connect-src 'self'; frame-ancestors 'none'")
		if strings.HasPrefix(r.URL.Path, "/assets/") {
			w.Header().Set("Cache-Control", "public,max-age=31536000,immutable")
		} else {
			w.Header().Set("Cache-Control", "no-cache")
		}
		limpio := filepath.Clean("/" + r.URL.Path)
		if _, err := os.Stat(filepath.Join("/srv", limpio)); err != nil && !strings.HasPrefix(r.URL.Path, "/assets/") {
			r.URL.Path = "/"
		}
		archivos.ServeHTTP(w, r)
	})
	s := &http.Server{Addr: ":8080", Handler: mux, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 30 * time.Second, WriteTimeout: 30 * time.Second}
	if err := s.ListenAndServe(); err != nil {
		panic(err)
	}
}
