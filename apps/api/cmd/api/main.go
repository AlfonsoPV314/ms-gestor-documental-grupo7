package main

import (
	"bufio"
	"context"
	"crypto/ed25519"
	"errors"
	"github.com/PINGESO-CIVIL/ms-gestor-documental/apps/api/internal/api"
	"github.com/PINGESO-CIVIL/ms-gestor-documental/apps/api/internal/config"
	"github.com/PINGESO-CIVIL/ms-gestor-documental/apps/api/internal/service"
	"github.com/PINGESO-CIVIL/ms-gestor-documental/apps/api/internal/store"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime/coverage"
	"syscall"
	"time"
)

var version = "dev"

func main() {
	if os.Getenv("API_ENV") == "test" && os.Getenv("API_COVER_DIR") != "" {
		defer func() {
			if e := coverage.WriteMetaDir(os.Getenv("API_COVER_DIR")); e != nil {
				slog.Error("cobertura meta", "error", e)
			}
			if e := coverage.WriteCountersDir(os.Getenv("API_COVER_DIR")); e != nil {
				slog.Error("cobertura contadores", "error", e)
			}
		}()
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	if err := ejecutar(); err != nil {
		slog.Error("arranque fallido", "error", err)
		os.Exit(1)
	}
}
func ejecutar() error {
	// codegen no necesita red, PostgreSQL ni claves reales.
	if len(os.Args) > 1 && os.Args[1] == "spec" {
		clave := ed25519.NewKeyFromSeed(make([]byte, 32))
		a := api.Nueva(http.NewServeMux(), service.Nuevo(nil, clave), "dev", true, nil)
		b, e := a.OpenAPI().YAML()
		if e != nil {
			return e
		}
		_, e = os.Stdout.Write(b)
		return e
	}
	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		c := &http.Client{Timeout: 2 * time.Second}
		r, e := c.Get("http://127.0.0.1:8080/healthz")
		if e != nil {
			return e
		}
		defer r.Body.Close()
		if r.StatusCode != 200 {
			return errors.New("no saludable")
		}
		return nil
	}
	c, err := config.Leer()
	if err != nil {
		return err
	}
	clave, err := config.LeerClave(c.ClaveRuta)
	if err != nil {
		return err
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if c.Ambiente == "test" && os.Getenv("API_TEST_CONTROL") == "stdin" {
		go func() { _, _ = bufio.NewReader(os.Stdin).ReadString('\n'); stop() }()
	}
	pc, err := pgxpool.ParseConfig(c.Base)
	if err != nil {
		return errors.New("DATABASE_URL inválida")
	}
	pc.ConnConfig.ConnectTimeout = 5 * time.Second
	pc.MaxConns = 10
	pool, err := pgxpool.NewWithConfig(ctx, pc)
	if err != nil {
		return err
	}
	defer pool.Close()
	pctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	err = pool.Ping(pctx)
	cancel()
	if err != nil {
		return errors.New("PostgreSQL no disponible")
	}
	auth, err := api.Autenticacion(ctx, c)
	if err != nil {
		return errors.New("JWKS no disponible; revisar AUTH_JWKS_URL")
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"estado":"vivo"}`))
	})
	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) {
		cctx, ccancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer ccancel()
		w.Header().Set("Content-Type", "application/json")
		if pool.Ping(cctx) != nil {
			w.WriteHeader(503)
			_, _ = w.Write([]byte(`{"estado":"no_disponible"}`))
			return
		}
		_, _ = w.Write([]byte(`{"estado":"listo"}`))
	})
	api.Nueva(mux, service.Nuevo(&store.Postgres{Pool: pool}, clave), c.Ambiente, c.Mocks, auth)
	servidor := &http.Server{Addr: ":" + c.Puerto, Handler: api.Observar(mux, version), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 30 * time.Second, WriteTimeout: 30 * time.Second, IdleTimeout: 60 * time.Second, MaxHeaderBytes: 16 * 1024}
	go func() {
		<-ctx.Done()
		cctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = servidor.Shutdown(cctx)
	}()
	slog.Info("API iniciada", "version", version, "mocks", c.Mocks)
	err = servidor.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}
