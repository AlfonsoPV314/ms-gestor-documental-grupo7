package api

import (
	"context"
	"github.com/MicahParks/keyfunc/v3"
	"github.com/PINGESO-CIVIL/ms-gestor-documental/apps/api/internal/config"
	"github.com/PINGESO-CIVIL/ms-gestor-documental/apps/api/internal/service"
	"github.com/danielgtaylor/huma/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"net/http"
	"strings"
	"time"
)

type claveContexto string

const identidadContexto claveContexto = "identidad"
const peticionContexto claveContexto = "request_id"

func identidad(ctx context.Context) service.Identidad {
	i, _ := ctx.Value(identidadContexto).(service.Identidad)
	return i
}

type Reclamaciones struct {
	jwt.RegisteredClaims
	Permisos []string `json:"permisos"`
	Rol      string   `json:"rol"`
}

func Autenticacion(ctx context.Context, c config.Config) (func(huma.API), error) {
	jwks, err := keyfunc.NewDefaultOverrideCtx(ctx, []string{c.JWKS}, keyfunc.Override{Client: &http.Client{Timeout: 5 * time.Second}, HTTPTimeout: 5 * time.Second, RateLimitWaitMax: 5 * time.Second})
	if err != nil {
		return nil, err
	}
	return func(a huma.API) {
		a.UseMiddleware(func(hctx huma.Context, next func(huma.Context)) {
			if hctx.Operation().Security == nil {
				next(hctx)
				return
			}
			token := strings.TrimPrefix(hctx.Header("Authorization"), "Bearer ")
			usaCookie := false
			if token == "" {
				req := &http.Request{Header: http.Header{"Cookie": []string{hctx.Header("Cookie")}}}
				if cookie, e := req.Cookie(c.Cookie); e == nil {
					token = cookie.Value
					usaCookie = true
				}
			}
			if usaCookie && hctx.Method() != "GET" && hctx.Method() != "HEAD" && hctx.Header("Origin") != c.Frontend {
				_ = huma.WriteErr(a, hctx, 403, "Origen no permitido")
				return
			}
			claims := &Reclamaciones{}
			_, err := jwt.ParseWithClaims(token, claims, jwks.KeyfuncCtx(hctx.Context()), jwt.WithValidMethods([]string{"EdDSA", "RS256"}), jwt.WithIssuer(c.Emisor), jwt.WithAudience(c.Audiencia), jwt.WithExpirationRequired())
			var id pgtype.UUID
			if err != nil || claims.Subject == "" || id.Scan(claims.Subject) != nil {
				_ = huma.WriteErr(a, hctx, 401, "Sesión inválida o vencida")
				return
			}
			permiso, _ := hctx.Operation().Extensions["x-permiso"].(string)
			i := service.Identidad{ID: claims.Subject, Permisos: claims.Permisos}
			if !i.Permite(permiso) {
				_ = huma.WriteErr(a, hctx, 403, "Permiso insuficiente")
				return
			}
			next(huma.WithContext(hctx, context.WithValue(hctx.Context(), identidadContexto, i)))
		})
	}, nil
}
