# Guía breve de demostración

Ejecutar desde la raíz del repositorio con PowerShell 7 y Docker Desktop Linux.
Usar exclusivamente archivos ficticios como `ejemplos/documento.txt`. Cada demo
real crea un documento y su auditoría; no elimina datos existentes.

## Preparación

```powershell
./scripts/herramientas.ps1
. ./scripts/activar.ps1
pnpm install --frozen-lockfile
./scripts/iniciar.ps1
```

App: <http://localhost:5170>. Swagger: <http://localhost:5170/api/docs>.
Seleccionar `operador` en la app antes de abrir Swagger para reutilizar la cookie.
Los tokens nunca se copian a archivos de evidencia ni a capturas.

## Hito 1: contrato y mocks

```powershell
$env:API_MOCKS='true'
docker compose up -d --wait backend frontend
node scripts/demo.mjs --mocks
Remove-Item Env:API_MOCKS
```

El script comprueba salud y disponibilidad, Swagger y OpenAPI, POST de archivo
con 201, listado, metadatos, descarga y verificación con 200. Sin sesión espera
401 y un lector no puede cargar (403). Los datos y la firma son simulados:
esta demo no demuestra persistencia, integridad criptográfica ni aislamiento.
El reporte sin tokens queda en `.artifacts/demo-mocks.json`.

Para mostrarlo manualmente: en Swagger ejecutar POST `/api/v1/documentos` con el
archivo de ejemplo, copiar el UUID de respuesta y consultar las tres operaciones
por ID, además del listado. Revisar [arquitectura](arquitectura.md) y
[modos de falla](modos-de-falla.md) durante la presentación.

## Hito 2: documentos reales

```powershell
$env:API_MOCKS='false'
docker compose up -d --wait backend frontend
node scripts/demo.mjs
Remove-Item Env:API_MOCKS
```

Ahora la descarga debe ser idéntica al original, su SHA-256 debe coincidir y la
firma debe ser válida. El lector consulta pero no carga (403), y otra identidad
no puede obtener el documento (404). Reporte: `.artifacts/demo-real.json`.
Esta variante usa identidad fixture, por lo que no cierra la integración externa.

Para mostrar persistencia, cargar desde la app, anotar el UUID, ejecutar
`docker compose restart backend postgres`, esperar `/readyz` 200, volver a
consultar ese UUID y verificarlo. No usar `down -v`: borraría los datos.
La prueba E2E automatiza el reinicio de la API y otros casos adversos.

## Integración con el Grupo 6

La API documental ya acepta RS256. No se reemplaza la validación por un nombre
de rol ni se consulta la base del proveedor. La prueba usa su implementación
real con tokens de desarrollo, no nuestro fixture.

```powershell
# Solo si todavía no existe la copia local:
git clone https://github.com/PINGESO-CIVIL/usuarios-roles.git .artifacts/usuarios-roles
git -C .artifacts/usuarios-roles checkout --detach 2d5ecccfd0fd27415ffbd4d6dd017ed895d94fd1
docker compose -f docker-compose.yaml -f compose.integracion.yaml up -d --build identidad-grupo6-db identidad-grupo6
# Esperar hasta que este comando responda 200:
(Invoke-WebRequest http://localhost:5161/.well-known/jwks.json).StatusCode
docker compose -f docker-compose.yaml -f compose.integracion.yaml up -d --wait backend frontend
node scripts/demo.mjs --grupo6
```

El reporte `.artifacts/demo-grupo6.json` comprueba carga, lectura, descarga,
firma, rechazo sin sesión, permisos insuficientes y aislamiento entre dos UUID.
El proveedor tiene su propia base en este Compose. Ningún componente documental
accede a ella. La versión fijada del grupo 6 aún guarda usuarios en memoria y sus
roles seed no incluyen permisos documentales: se usa `permisos_extra` solamente
en `/api/dev/tokens`. Falta acordar la asignación real de esos permisos y probar
la sesión del host cliente; el selector de identidades local no es su login.

Restaurar el proveedor fixture para las pruebas del navegador:

```powershell
docker compose up -d --wait backend frontend
docker compose -f docker-compose.yaml -f compose.integracion.yaml stop identidad-grupo6 identidad-grupo6-db
```

## Pruebas y evidencias de calidad

```powershell
pnpm verify
pnpm generate
pnpm codegen
git diff --exit-code -- contracts/openapi.yaml apps/web/src/contrato.d.ts apps/api/internal/repository
./scripts/probar.ps1
./scripts/probar.ps1 -Mutacion
pnpm test:e2e:web
pnpm audit --audit-level high
node scripts/seguridad.mjs
./scripts/iniciar.ps1 -Calidad
# Esperar SonarQube UP; configurar solo si todavía no existe el proyecto local:
# node scripts/configurar-sonar.mjs
pnpm sonar
```

En una máquina nueva, instalar Chromium con
`pnpm --filter @pingeso-civil/ms-gestor-documental-ui exec playwright install chromium`.
Go genera `apps/api/coverage.out` desde la API instrumentada y PostgreSQL real;
Vitest genera `apps/web/coverage/lcov.info`. La mutación guarda
`.artifacts/mutacion.log` y restaura el código incluso si falla. En CI, conservar
los artifacts `evidencia-pruebas`, `evidencia-seguridad` y `evidencia-imagenes`.
El workflow manual de evidencia roja debe fallar por la aserción de acceso ajeno.

Los reportes de demo incluyen commit y si había cambios pendientes. Una prueba
hecha con cambios sin commit sirve para desarrollo, no como evidencia final del
tag. Repetir los controles sobre el commit definitivo y enlazar su run.

## Detener todo lo levantado por esta guía

```powershell
docker compose -f docker-compose.yaml -f compose.integracion.yaml stop
docker compose -f compose.quality.yaml stop
```

No se borran volúmenes ni claves. No adjuntar `.env`, `.secrets`, tokens ni cookies
a releases o issues. Ver [checklist de cierre](cierre-entrega.md).
