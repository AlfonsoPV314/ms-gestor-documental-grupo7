# Integración entre grupos

## Grupo 6: contrato recibido, permisos y sesión final pendientes

El contrato real está copiado en `contracts/externos/identidad.yaml`, desde el
repositorio <https://github.com/PINGESO-CIVIL/usuarios-roles>, commit
`2d5ecccfd0fd27415ffbd4d6dd017ed895d94fd1`. `identidad-dev` sigue siendo un
fixture explícito de desarrollo; no se atribuye al grupo 6.

Su implementación usa RS256, JWKS en `/.well-known/jwks.json`, issuer
`pingeso-identidad`, audience `pingeso`, cookie `sesion` y permisos firmados.
La API documental ya soporta esos valores mediante configuración. La prueba de
[demo](demo.md) usa su código real; no consulta su base ni necesita introspección
por request. No se genera un cliente CRUD de usuarios porque no consumimos esas
operaciones: keyfunc implementa el cliente JWKS con timeout de cinco segundos.

Claims relevantes del JWT recibido:

```json
{
  "sub": "UUID del usuario",
  "iss": "pingeso-identidad",
  "aud": ["pingeso"],
  "exp": 1796000000,
  "permisos": ["documentos:leer", "documentos:operar"],
  "roles": ["operador"]
}
```

La autorización se decide por permisos firmados, no por roles enviados desde la
interfaz. El arreglo `roles` no se usa para conceder acceso. La matriz está en
acceso_test.go y los endpoints se prueban además en E2E. La sesión usa cookie
HttpOnly, Secure y SameSite=Lax; no se guarda en almacenamiento web. Clientes de
API pueden usar Bearer JWT.

Pendiente: acordar qué roles reales reciben `documentos:leer` y
`documentos:operar`, las URLs del despliegue y la sesión del host cliente. Los
roles seed del grupo 6 no tienen permisos documentales; `permisos_extra` y el
almacenamiento en memoria de esa versión son herramientas de desarrollo, no
evidencia de cierre productivo de Hito 2. El selector de identidad de la app
local sigue siendo el fixture; no se presenta como un login real del grupo 6.

## Grupo 5: consumidor

Compartir `contracts/openapi.yaml` de una versión publicada. Cargar con POST multipart `/api/v1/documentos`, campo `archivo`; reenviar el JWT del usuario y X-Request-Id. La respuesta 201 contiene id, propietario, nombre, tipo, tamaño, sha256, firma, clave_id y creado_en. No acceder a PostgreSQL del grupo 7.

Consulta y descarga solo son posibles para el propietario con permiso de lectura. El grupo 5 debe reenviar la identidad de ese propietario. Los errores usan application/problem+json; 401 sesión, 403 permiso, 404 documento inaccesible, 413 tamaño, 422 validación, 503 dependencia no disponible.

## Microcomponentes para grupos cliente

```tsx
import {CargaDocumental,Documentos} from '@pingeso-civil/ms-gestor-documental-ui'
import '@pingeso-civil/ms-gestor-documental-ui/style.css'
// Montarlos por separado; el contenedor coordina revision con alCargar.
<CargaDocumental baseUrl="" puedeCargar={permisos.includes('documentos:operar')} alCargar={recargar}/>
<Documentos baseUrl="" revision={revision}/>
```

El host debe servir la API por el mismo origen y gestionar su sesión real. Para un origen distinto se necesita acordar CORS, cookies y CSRF; la versión inicial no habilita CORS indiscriminado. El paquete declara contrato API 1.0.0. La app local y el selector de identidades no se exportan.

## Auditoría

Cada creación registra actor, acción, documento, timestamp y request_id en PostgreSQL. No se inventa un contrato de envío al grupo 8; la integración central requiere su contrato y acuerdo de alcance. Los logs JSON de HTTP no contienen archivos, cookies ni tokens.
