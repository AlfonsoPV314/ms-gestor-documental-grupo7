# Integración entre grupos

## Grupo 6: proveedor de identidad pendiente

El contrato real no fue entregado. No se ha inventado ni publicado un contrato externo atribuido al grupo 6. `identidad-dev` es un fixture explícito para local y pruebas; su código no forma parte de la imagen API ni del paquete UI.

Solicitar al grupo 6: contrato OpenAPI, URL interna, JWKS, issuer, audience, nombre de cookie, algoritmos, claims de permisos y tokens de prueba. Copiar su contrato a `contracts/externos/identidad.yaml`, generar cliente Go mediante oapi-codegen v2.8.2 si sus operaciones son necesarias y revisar la adaptación de claims antes de entregar Hito 2. Validar JWKS no equivale a necesitar una llamada de introspección por cada request.

Contrato provisional esperado del JWT:

```json
{
  "sub": "UUID del usuario",
  "iss": "emisor acordado",
  "aud": "documental",
  "exp": 1796000000,
  "permisos": ["documentos:leer", "documentos:operar"],
  "rol": "operador"
}
```

La autorización se decide por permisos firmados, no por el rol enviado desde la interfaz. El campo rol se conserva como contexto provisional. La matriz probada está en acceso_test.go y los endpoints se prueban además en E2E. La sesión usa cookie HttpOnly, Secure y SameSite=Lax; no se guarda en almacenamiento web. Clientes de API pueden usar Bearer JWT.

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
