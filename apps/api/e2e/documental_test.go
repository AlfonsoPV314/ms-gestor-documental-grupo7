package e2e

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestFlujoDocumental(t *testing.T) {
	adminURL := os.Getenv("TEST_DATABASE_URL")
	if adminURL == "" {
		t.Skip("requiere TEST_DATABASE_URL; no sustituye PostgreSQL con mocks")
	}
	ctx := context.Background()
	admin, err := pgx.Connect(ctx, adminURL)
	if err != nil {
		t.Fatal(err)
	}
	defer admin.Close(ctx)
	nombre := fmt.Sprintf("documental_test_%d", time.Now().UnixNano())
	if _, err = admin.Exec(ctx, "CREATE DATABASE "+pgx.Identifier{nombre}.Sanitize()); err != nil {
		t.Fatal(err)
	}
	defer func() { _, _ = admin.Exec(ctx, "DROP DATABASE "+pgx.Identifier{nombre}.Sanitize()+" WITH (FORCE)") }()
	dsn, err := url.Parse(adminURL)
	if err != nil {
		t.Fatal(err)
	}
	dsn.Path = "/" + nombre
	cfg, err := pgx.ParseConfig(dsn.String())
	if err != nil {
		t.Fatal(err)
	}
	db, err := pgx.ConnectConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close(ctx)
	migration, err := os.ReadFile("../db/migrations/00001_documentos.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = db.Exec(ctx, strings.Split(string(migration), "-- +goose Down")[0]); err != nil {
		t.Fatal(err)
	}
	publica, privada, _ := ed25519.GenerateKey(rand.Reader)
	jwks := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{"keys": []any{map[string]any{"kty": "OKP", "crv": "Ed25519", "alg": "EdDSA", "kid": "prueba", "x": base64.RawURLEncoding.EncodeToString(publica)}}})
	}))
	defer jwks.Close()
	token := func(sub string, permisos []string, exp time.Time, alg string) string {
		var metodo jwt.SigningMethod = jwt.SigningMethodEdDSA
		var clave any = privada
		if alg == "HS256" {
			metodo = jwt.SigningMethodHS256
			clave = []byte("clave-sintetica-solo-pruebas")
		}
		j := jwt.NewWithClaims(metodo, jwt.MapClaims{"sub": sub, "iss": "e2e", "aud": "documental", "exp": exp.Unix(), "permisos": permisos})
		j.Header["kid"] = "prueba"
		s, e := j.SignedString(clave)
		if e != nil {
			t.Fatal(e)
		}
		return s
	}
	propietario := "22222222-2222-4222-8222-222222222222"
	operador := token(propietario, []string{"documentos:leer", "documentos:operar"}, time.Now().Add(time.Hour), "")
	lector := token(propietario, []string{"documentos:leer"}, time.Now().Add(time.Hour), "")
	ajeno := token("33333333-3333-4333-8333-333333333333", []string{"documentos:leer", "documentos:operar"}, time.Now().Add(time.Hour), "")
	dir := t.TempDir()
	claveRuta := filepath.Join(dir, "firma.key")
	if err = os.WriteFile(claveRuta, []byte(base64.StdEncoding.EncodeToString(privada)), 0600); err != nil {
		t.Fatal(err)
	}
	bin := os.Getenv("API_BINARY")
	if bin == "" {
		t.Fatal("API_BINARY requerido: ejecutar pnpm test:e2e:api")
	}
	logFile, err := os.Create(filepath.Join(dir, "api.log"))
	if err != nil {
		t.Fatal(err)
	}
	defer logFile.Close()
	var entrada io.WriteCloser
	iniciar := func() *exec.Cmd {
		cmd := exec.Command(bin)
		cmd.Env = append(os.Environ(), "API_PORT=8080", "API_ENV=test", "API_TEST_CONTROL=stdin", "API_MOCKS=false", "DATABASE_URL="+cfg.ConnString(), "AUTH_JWKS_URL="+jwks.URL, "AUTH_ISSUER=e2e", "AUTH_AUDIENCE=documental", "AUTH_COOKIE_NAME=sesion", "FRONTEND_URL=http://localhost:5170", "SIGNING_KEY_FILE="+claveRuta, "TZ="+os.Getenv("TZ"))
		cmd.Stdout = logFile
		cmd.Stderr = logFile
		entrada, err = cmd.StdinPipe()
		if err != nil {
			t.Fatal(err)
		}
		if err := cmd.Start(); err != nil {
			t.Fatal(err)
		}
		return cmd
	}
	detener := func(cmd *exec.Cmd) {
		_ = entrada.Close()
		done := make(chan error, 1)
		go func() { done <- cmd.Wait() }()
		select {
		case <-done:
		case <-time.After(10 * time.Second):
			_ = cmd.Process.Kill()
			<-done
		}
		if b, e := os.ReadFile(logFile.Name()); e == nil {
			_ = os.MkdirAll("../../../.artifacts", 0700)
			_ = os.WriteFile("../../../.artifacts/e2e-api.log", b, 0600)
		}
	}
	cmd := iniciar()
	defer func() {
		if cmd != nil {
			detener(cmd)
		}
	}()
	client := &http.Client{Timeout: 15 * time.Second}
	request := func(metodo, ruta, token string, cuerpo io.Reader, tipo string) (int, []byte, http.Header) {
		req, e := http.NewRequest(metodo, "http://127.0.0.1:8080"+ruta, cuerpo)
		if e != nil {
			t.Fatal(e)
		}
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		if tipo != "" {
			req.Header.Set("Content-Type", tipo)
		}
		req.Header.Set("X-Request-Id", "e2e-correlacion")
		r, e := client.Do(req)
		if e != nil {
			t.Fatal(e)
		}
		defer r.Body.Close()
		b, e := io.ReadAll(r.Body)
		if e != nil {
			t.Fatal(e)
		}
		return r.StatusCode, b, r.Header
	}
	esperar := func() {
		t.Helper()
		for i := 0; i < 100; i++ {
			r, e := client.Get("http://127.0.0.1:8080/readyz")
			if e == nil {
				r.Body.Close()
				if r.StatusCode == 200 {
					return
				}
			}
			time.Sleep(100 * time.Millisecond)
		}
		b, _ := os.ReadFile(logFile.Name())
		t.Fatalf("API no inició: %s", b)
	}
	esperar()
	var reporte strings.Builder
	exigir := func(nombre string, actual, esperado int) {
		t.Helper()
		if actual != esperado {
			t.Fatalf("%s: estado %d, esperado %d", nombre, actual, esperado)
		}
		fmt.Fprintf(&reporte, "- %s: %d\n", nombre, actual)
	}
	cargar := func(nombre string, b []byte, token string) (int, []byte) {
		var buf bytes.Buffer
		w := multipart.NewWriter(&buf)
		p, e := w.CreateFormFile("archivo", nombre)
		if e != nil {
			t.Fatal(e)
		}
		_, _ = p.Write(b)
		_ = w.Close()
		status, res, _ := request("POST", "/api/v1/documentos", token, &buf, w.FormDataContentType())
		return status, res
	}
	status, _, _ := request("GET", "/api/v1/documentos", "", nil, "")
	exigir("sin sesión", status, 401)
	status, _, _ = request("GET", "/api/v1/documentos", token(propietario, []string{"documentos:leer"}, time.Now().Add(-time.Hour), ""), nil, "")
	exigir("token vencido", status, 401)
	status, _, _ = request("GET", "/api/v1/documentos", token(propietario, []string{"documentos:leer"}, time.Now().Add(time.Hour), "HS256"), nil, "")
	exigir("algoritmo ajeno", status, 401)
	status, _, _ = request("GET", "/api/v1/documentos", token(propietario, nil, time.Now().Add(time.Hour), ""), nil, "")
	exigir("sin permiso", status, 403)
	status, _ = cargar("lector.txt", []byte("texto"), lector)
	exigir("lector no carga", status, 403)
	status, _ = cargar("vacio.txt", nil, operador)
	exigir("archivo vacío", status, 422)
	status, _ = cargar("..", []byte("texto"), operador)
	exigir("nombre inseguro", status, 422)
	status, _ = cargar("grande.txt", bytes.Repeat([]byte("a"), 20*1024*1024+1), operador)
	exigir("archivo grande", status, 413)
	contenido := []byte("contenido original para firma\n")
	status, b := cargar("informe.txt", contenido, operador)
	if status != 201 {
		t.Logf("respuesta carga: %s", b)
		rows, e := db.Query(ctx, "SELECT name,setting FROM pg_settings WHERE name='transaction_read_only'")
		if e == nil {
			for rows.Next() {
				var n, v string
				_ = rows.Scan(&n, &v)
				t.Log(n, v)
			}
			rows.Close()
		}
	}
	exigir("carga original", status, 201)
	var doc struct {
		ID   string `json:"id"`
		Hash string `json:"sha256"`
	}
	if err = json.Unmarshal(b, &doc); err != nil || doc.ID == "" {
		t.Fatalf("respuesta carga: %s", b)
	}
	for _, sufijo := range []string{"", "/contenido", "/verificacion"} {
		status, _, _ = request("GET", "/api/v1/documentos/"+doc.ID+sufijo, ajeno, nil, "")
		exigir("ajeno"+sufijo, status, 404)
	}
	status, b, h := request("GET", "/api/v1/documentos/"+doc.ID+"/contenido", lector, nil, "")
	exigir("descarga lector", status, 200)
	if !bytes.Equal(b, contenido) || !strings.Contains(h.Get("Content-Disposition"), "attachment") {
		t.Fatal("descarga distinta al original")
	}
	status, b, _ = request("GET", "/api/v1/documentos/"+doc.ID+"/verificacion", lector, nil, "")
	exigir("verificación original", status, 200)
	if !bytes.Contains(b, []byte(`"valida":true`)) {
		t.Fatalf("firma inválida: %s", b)
	}
	if _, err = db.Exec(ctx, "UPDATE documentos SET contenido=$1 WHERE id=$2", bytes.Repeat([]byte("x"), len(contenido)), doc.ID); err != nil {
		t.Fatal(err)
	}
	_, b, _ = request("GET", "/api/v1/documentos/"+doc.ID+"/verificacion", lector, nil, "")
	if !bytes.Contains(b, []byte(`"valida":false`)) {
		t.Fatal("no detectó manipulación")
	}
	reporte.WriteString("- contenido alterado: detectado\n")
	_, _ = db.Exec(ctx, "UPDATE documentos SET contenido=$1 WHERE id=$2", contenido, doc.ID)
	var wg sync.WaitGroup
	estados := make(chan int, 6)
	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s, _ := cargar("concurrente.txt", []byte("concurrencia"), operador)
			estados <- s
		}()
	}
	wg.Wait()
	close(estados)
	for s := range estados {
		if s != 201 {
			t.Fatalf("carga concurrente: %d", s)
		}
	}
	var docs, audits int
	_ = db.QueryRow(ctx, "SELECT count(*) FROM documentos").Scan(&docs)
	_ = db.QueryRow(ctx, "SELECT count(*) FROM auditoria").Scan(&audits)
	if docs != 7 || audits != 7 {
		t.Fatalf("inconsistencia: %d documentos %d auditorías", docs, audits)
	}
	reporte.WriteString("- 6 cargas concurrentes: documentos y auditorías consistentes\n")
	_, err = db.Exec(ctx, `CREATE FUNCTION rechazar_auditoria() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RAISE EXCEPTION 'fallo de prueba'; END $$; CREATE TRIGGER fallo BEFORE INSERT ON auditoria FOR EACH ROW EXECUTE FUNCTION rechazar_auditoria();`)
	if err != nil {
		t.Fatal(err)
	}
	status, _ = cargar("rollback.txt", []byte("rollback"), operador)
	exigir("fallo auditoría", status, 503)
	_ = db.QueryRow(ctx, "SELECT count(*) FROM documentos").Scan(&docs)
	if docs != 7 {
		t.Fatal("quedó documento parcial")
	}
	reporte.WriteString("- fallo auditoría: rollback completo\n")
	_, _ = db.Exec(ctx, "DROP TRIGGER fallo ON auditoria; DROP FUNCTION rechazar_auditoria()")
	// Cookie mutante exige Origin correcto.
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8080/api/v1/documentos", nil)
	req.AddCookie(&http.Cookie{Name: "sesion", Value: operador})
	req.Header.Set("Origin", "https://ajeno.invalid")
	r, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	r.Body.Close()
	exigir("origen ajeno", r.StatusCode, 403)
	detener(cmd)
	cmd = iniciar()
	esperar()
	_, b, _ = request("GET", "/api/v1/documentos/"+doc.ID+"/verificacion", lector, nil, "")
	if !bytes.Contains(b, []byte(`"valida":true`)) {
		t.Fatal("firma no sobrevive reinicio")
	}
	reporte.WriteString("- reinicio: firma y documento conservados\n")
	golden, err := os.ReadFile("testdata/documental.golden.md")
	if err != nil {
		t.Fatal(err)
	}
	actual := "# Flujo documental\n\n" + reporte.String()
	if strings.ReplaceAll(string(golden), "\r\n", "\n") != actual {
		t.Fatalf("golden distinto\n%s", actual)
	}
}
