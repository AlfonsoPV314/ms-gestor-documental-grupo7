package api

import (
	"context"
	"errors"
	"github.com/PINGESO-CIVIL/ms-gestor-documental/apps/api/internal/service"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"io"
	"mime"
	"net/http"
	"time"
)

type EntradaID struct {
	ID string `path:"id" format:"uuid" example:"11111111-1111-4111-8111-111111111111"`
}
type EntradaLista struct {
	Limite int32 `query:"limite" minimum:"1" maximum:"100" default:"20"`
	Offset int32 `query:"offset" minimum:"0" maximum:"100000" default:"0"`
}
type SalidaDocumento struct{ Body Documento }
type SalidaLista struct {
	Body struct {
		Documentos []Documento `json:"documentos"`
	}
}
type Formulario struct {
	Archivo huma.FormFile `form:"archivo" required:"true" contentType:"application/pdf,image/png,image/jpeg,text/plain"`
}
type EntradaCarga struct {
	RawBody huma.MultipartFormFiles[Formulario]
}
type SalidaVerificacion struct {
	Body struct {
		Valida    bool   `json:"valida"`
		Algoritmo string `json:"algoritmo"`
		ClaveID   string `json:"clave_id"`
	}
}
type SalidaContenido struct {
	Tipo        string `header:"Content-Type"`
	Disposicion string `header:"Content-Disposition"`
	Cache       string `header:"Cache-Control"`
	Body        []byte
}

func traducir(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, service.ErrNoEncontrado):
		return huma.Error404NotFound("Documento no encontrado")
	case errors.Is(err, service.ErrProhibido):
		return huma.Error403Forbidden("Permiso insuficiente")
	case errors.Is(err, service.ErrEntrada):
		return huma.Error422UnprocessableEntity("Archivo inválido; se admiten PDF, PNG, JPEG y texto")
	case errors.Is(err, service.ErrGrande):
		return huma.NewError(413, "Máximo 20 MiB")
	default:
		return huma.Error503ServiceUnavailable("Almacenamiento temporalmente no disponible")
	}
}
func Nueva(mux *http.ServeMux, s *service.Servicio, ambiente string, mocks bool, auth func(huma.API)) huma.API {
	cfg := huma.DefaultConfig("MS Gestor Documental — Grupo 7", "1.0.0")
	cfg.DocsPath = "/api/docs"
	cfg.OpenAPIPath = "/api/openapi"
	cfg.SchemasPath = "/api/schemas"
	if ambiente == "prod" {
		cfg.DocsPath = ""
		cfg.OpenAPIPath = ""
		cfg.SchemasPath = ""
	}
	cfg.Components.SecuritySchemes = map[string]*huma.SecurityScheme{"sesion": {Type: "apiKey", In: "cookie", Name: "sesion"}, "bearer": {Type: "http", Scheme: "bearer", BearerFormat: "JWT"}}
	a := humago.New(mux, cfg)
	if auth != nil {
		auth(a)
	}
	op := func(id, metodo, ruta, resumen, permiso string) huma.Operation {
		return huma.Operation{OperationID: id, Method: metodo, Path: ruta, Summary: resumen, Tags: []string{"Documentos"}, Security: []map[string][]string{{"sesion": {}}, {"bearer": {}}}, Extensions: map[string]any{"x-permiso": permiso}, MaxBodyBytes: service.MaxArchivo + 1024*64, Errors: []int{401, 403, 404, 413, 422, 503}}
	}
	ejemplo := service.Documento{ID: "11111111-1111-4111-8111-111111111111", Propietario: "22222222-2222-4222-8222-222222222222", Nombre: "ejemplo.txt", Tipo: "text/plain", Tamano: 7, Hash: "0000000000000000000000000000000000000000000000000000000000000000", Firma: "simulada", ClaveID: "simulada", CreadoEn: time.Date(2026, 11, 23, 0, 0, 0, 0, time.UTC), Contenido: []byte("ejemplo")}
	carga := op("cargar-documento", "POST", "/api/v1/documentos", "Cargar y firmar documento", "documentos:operar")
	carga.DefaultStatus = 201
	huma.Register(a, carga, func(ctx context.Context, in *EntradaCarga) (*SalidaDocumento, error) {
		if in.RawBody.Form != nil {
			defer func() { _ = in.RawBody.Form.RemoveAll() }()
		}
		f := in.RawBody.Data().Archivo
		if f.File == nil {
			return nil, huma.Error422UnprocessableEntity("Falta archivo")
		}
		defer f.Close()
		b, err := io.ReadAll(io.LimitReader(f.File, service.MaxArchivo+1))
		if err != nil {
			return nil, huma.Error400BadRequest("No se pudo leer el archivo")
		}
		if mocks {
			return &SalidaDocumento{Body: salida(ejemplo)}, nil
		}
		rid, _ := ctx.Value(peticionContexto).(string)
		d, err := s.Cargar(ctx, identidad(ctx), f.Filename, b, rid)
		if err != nil {
			return nil, traducir(err)
		}
		return &SalidaDocumento{Body: salida(d)}, nil
	})
	huma.Register(a, op("listar-documentos", "GET", "/api/v1/documentos", "Listar documentos propios", "documentos:leer"), func(ctx context.Context, in *EntradaLista) (*SalidaLista, error) {
		out := &SalidaLista{}
		if mocks {
			out.Body.Documentos = []Documento{salida(ejemplo)}
			return out, nil
		}
		ds, err := s.Listar(ctx, identidad(ctx), in.Limite, in.Offset)
		out.Body.Documentos = make([]Documento, 0, len(ds))
		for _, d := range ds {
			out.Body.Documentos = append(out.Body.Documentos, salida(d))
		}
		return out, traducir(err)
	})
	obtener := func(ctx context.Context, id string) (service.Documento, error) {
		if mocks {
			return ejemplo, nil
		}
		d, e := s.Obtener(ctx, identidad(ctx), id)
		return d, traducir(e)
	}
	huma.Register(a, op("consultar-documento", "GET", "/api/v1/documentos/{id}", "Consultar metadatos", "documentos:leer"), func(ctx context.Context, in *EntradaID) (*SalidaDocumento, error) {
		d, e := obtener(ctx, in.ID)
		return &SalidaDocumento{Body: salida(d)}, e
	})
	huma.Register(a, op("descargar-documento", "GET", "/api/v1/documentos/{id}/contenido", "Descargar original", "documentos:leer"), func(ctx context.Context, in *EntradaID) (*SalidaContenido, error) {
		d, e := obtener(ctx, in.ID)
		if e != nil {
			return nil, e
		}
		return &SalidaContenido{Tipo: d.Tipo, Disposicion: mime.FormatMediaType("attachment", map[string]string{"filename": d.Nombre}), Cache: "no-store", Body: d.Contenido}, nil
	})
	huma.Register(a, op("verificar-documento", "GET", "/api/v1/documentos/{id}/verificacion", "Verificar firma e integridad", "documentos:leer"), func(ctx context.Context, in *EntradaID) (*SalidaVerificacion, error) {
		d, e := obtener(ctx, in.ID)
		if e != nil {
			return nil, e
		}
		out := &SalidaVerificacion{}
		out.Body.Valida = mocks || s.Verificar(d)
		out.Body.Algoritmo = "Ed25519 + SHA-256 / documental:v1"
		out.Body.ClaveID = d.ClaveID
		return out, nil
	})

	for _, ruta := range []string{"/healthz", "/readyz"} {
		a.OpenAPI().Paths[ruta] = &huma.PathItem{Get: &huma.Operation{OperationID: map[string]string{"/healthz": "salud-proceso", "/readyz": "salud-postgresql"}[ruta], Summary: "Estado del servicio", Tags: []string{"Salud"}, Responses: map[string]*huma.Response{"200": {Description: "Servicio disponible"}, "503": {Description: "Dependencia no disponible"}}}}
	}
	return a
}
