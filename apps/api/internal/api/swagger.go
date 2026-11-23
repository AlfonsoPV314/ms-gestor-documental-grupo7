package api

import (
	"net/http"
	"os"
	"path/filepath"
)

const paginaSwagger = `<!doctype html>
<html lang="es">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Swagger UI · Gestor Documental</title>
  <link rel="stylesheet" href="/api/docs/swagger-ui.css">
</head>
<body>
  <p>API Gestor Documental · Grupo 7. Para usar tu sesión local, abre esta página desde el mismo origen de la aplicación. También puedes ingresar un JWT en Authorize.</p>
  <div id="swagger-ui"></div>
  <script src="/api/docs/swagger-ui-bundle.js"></script>
  <script src="/api/docs/iniciar.js"></script>
</body>
</html>`

const iniciarSwagger = `window.onload = function () {
  SwaggerUIBundle({
    url: '/api/openapi.json',
    dom_id: '#swagger-ui',
    deepLinking: true,
    displayRequestDuration: true,
    persistAuthorization: false,
    withCredentials: true,
    validatorUrl: null
  });
};`

func registrarSwagger(mux *http.ServeMux) {
	pagina := func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; connect-src 'self'; frame-ancestors 'none'")
		_, _ = w.Write([]byte(paginaSwagger))
	}
	mux.HandleFunc("GET /api/docs", pagina)
	mux.HandleFunc("GET /api/docs/{$}", pagina)
	mux.HandleFunc("GET /api/docs/iniciar.js", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		_, _ = w.Write([]byte(iniciarSwagger))
	})
	for _, archivo := range []string{"swagger-ui.css", "swagger-ui-bundle.js"} {
		mux.HandleFunc("GET /api/docs/"+archivo, func(w http.ResponseWriter, r *http.Request) {
			raiz := os.Getenv("SWAGGER_ASSETS_DIR")
			if raiz == "" {
				raiz = "../../node_modules/swagger-ui-dist"
			}
			http.ServeFile(w, r, filepath.Join(raiz, archivo))
		})
	}
}
