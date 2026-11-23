# MS Gestor Documental · Grupo 7

Bóveda documental de la base común LPS: recibe PDF, PNG, JPEG o texto de hasta 20 MiB, conserva el original, firma su integridad y expone consulta, descarga y verificación. Go/Huma, React/TypeScript y PostgreSQL. Proyecto nuevo e independiente de la implementación previa de la carpeta superior.

## Inicio local en Windows

Requiere Docker Desktop iniciado con motor Linux y Go instalado. Los scripts preparan herramientas locales y no sustituyen versiones globales.

```powershell
./scripts/herramientas.ps1
. ./scripts/activar.ps1
pnpm install --frozen-lockfile
pnpm exec lefthook install
./scripts/iniciar.ps1
```

- Aplicación: <http://localhost:5170>
- Swagger UI: <http://localhost:5170/api/docs>
- API directa: <http://localhost:5171>
- OpenAPI servido: <http://localhost:5171/api/openapi.yaml>
- Salud: <http://localhost:5171/healthz> y <http://localhost:5171/readyz>

Seleccionar una identidad de prueba y cargar un archivo. `operador` carga y consulta; `lector` consulta los documentos del mismo propietario; `ajeno` tiene otra identidad y no puede verlos. La app de desarrollo no se publica como producto.

Huma genera el contrato OpenAPI y Swagger UI lo presenta para explorar y ejecutar los endpoints. Abrir Swagger desde el puerto 5170 permite usar la sesión iniciada en la aplicación. También está disponible en el puerto 5171, donde se puede proporcionar un JWT con `Authorize`. Sus recursos se sirven localmente, sin CDN; la documentación está deshabilitada en producción.

`iniciar.ps1` genera `.env` y `.secrets/firma.key` solo si faltan. Son privados y se ignoran en Git. Conservar la clave y la base juntas: la firma de documentos anteriores depende de esa clave. `docker compose stop` conserva datos. No usar `down -v` para detener el trabajo cotidiano.

## SonarQube local

```powershell
./scripts/iniciar.ps1 -Calidad
# Esperar a que http://localhost:9000/api/system/status indique UP.
node scripts/configurar-sonar.mjs
./scripts/probar.ps1
pnpm test
pnpm sonar
```

Panel: <http://localhost:9000/dashboard?id=ms-gestor-documental>. El configurador cambia admin/admin por una contraseña aleatoria, crea un token limitado al proyecto y configura el Quality Gate DAS. Credenciales en `.secrets/sonar.json`, nunca en el README, Git ni el chat. La instancia usa PostgreSQL propio y volúmenes independientes.

El primer análisis establece una línea base; un gate inicial sin condiciones de código nuevo no demuestra cobertura del código futuro. Revisar también la cobertura total y repetir el análisis tras cada cambio. No se disminuyen umbrales ni se excluyen archivos para aprobar.

## Validación

```powershell
pnpm verify
pnpm generate
pnpm codegen
./scripts/probar.ps1
./scripts/probar.ps1 -Mutacion
pnpm --filter @pingeso-civil/ms-gestor-documental-ui exec playwright install chromium
pnpm test:e2e:web
```

Los E2E crean y eliminan exclusivamente una base `documental_test_<id>`. Ejecutan la API real con PostgreSQL, comparan un golden, prueban UTC y Pacific/Auckland y generan `apps/api/coverage.out`. La mutación temporal se restaura incluso al fallar; su log queda en `.artifacts/mutacion.log`. El navegador crea documentos de prueba en la base de demostración.

## Contrato e integración

Contrato generado y versionado en `contracts/openapi.yaml`; tipos en `apps/web/src/contrato.d.ts`. La API usa UUID, errores RFC 9457 y permisos `documentos:leer` / `documentos:operar`. Ver [arquitectura](docs/arquitectura.md), [integración](docs/integracion.md), [modos de falla](docs/modos-de-falla.md) y [entrega](docs/entrega.md).

`API_MOCKS=true` habilita respuestas simuladas para preparar Hito 1; no se permite en producción. Por defecto la lógica documental es real. El proveedor de identidad local es una dependencia simulada explícita y debe reemplazarse por el grupo 6 para entregar Hito 2.

## Estado de publicación

Trabajo local. El repositorio remoto, las invitaciones de PINGESO-CIVIL, la protección de main, SonarQube accesible desde CI, la integración real del grupo 6 y las releases oficiales requieren los datos y accesos del curso. No se han creado tags hito-* ni se presenta un pipeline remoto como aprobado.

Las correcciones mínimas de seguridad respecto de las versiones originales del DAS están registradas en [actualizaciones de seguridad](docs/actualizaciones-seguridad.md) para revisión del curso.
