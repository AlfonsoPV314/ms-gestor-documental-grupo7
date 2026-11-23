# Arquitectura y decisiones

La entrada HTTP y los structs del contrato viven en `internal/api`; las reglas de acceso, validación de archivos y firma en `internal/service`; PostgreSQL en `internal/store`. El store utiliza exclusivamente consultas generadas con sqlc. Los handlers no importan store ni repository. Los datos documentales pertenecen solo al grupo 7.

## Flujo de carga

1. Validar JWT con algoritmo permitido, JWKS, issuer, audience, expiración y subject UUID. Leer permisos del token firmado.
2. Exigir `documentos:operar`; para cookies mutantes exigir Origin igual a FRONTEND_URL.
3. Huma procesa multipart. Limitar el request y leer como máximo 20 MiB + 1; validar nombre, tamaño y tipo detectado del contenido.
4. Generar UUID y SHA-256. Firmar la representación `documental:v1`, id, propietario, nombre, tipo, tamaño y hash, separados por saltos de línea. El nombre no admite controles ni saltos de línea.
5. En una transacción PostgreSQL, guardar archivo y auditoría de creación. Confirmar solo después del commit.

Los archivos son inmutables. Listado, consulta, descarga y verificación exigen `documentos:leer` y filtran por propietario en SQL. Un documento ajeno responde 404. El lector y el operador de demostración comparten subject para mostrar permisos de rol sobre el mismo propietario; el usuario ajeno tiene otro subject.

## Persistencia y firmas

Original en bytea; metadatos e integridad en la misma fila. Fechas timestamptz y tamaños bigint. Esta elección evita inconsistencias entre un filesystem y la base y es suficiente para archivos de 20 MiB en el alcance académico. No representa una solución para archivos masivos.

Ed25519 firma desde el servicio; no se afirma que sea firma personal ni firma electrónica avanzada. La clave privada se monta de solo lectura desde `.secrets/firma.key`. El id de clave deriva del hash de su clave pública. Las firmas sobreviven reinicios mientras se conserve esa clave. La primera versión no ofrece rotación automática: reemplazar la clave rompe la verificación de documentos anteriores. Respaldar base y clave juntas; validar restauración antes de cambios operativos.

## Microcomponentes

`CargaDocumental` y `Documentos` son componentes React independientes. Reciben baseUrl y contexto por props, sin importar entre sí ni exigir router, layout o estado global. Sus clientes se generan/tipan desde OpenAPI y cada fetch vive en un archivo `.api.ts`. Los estilos de la librería se limitan a `.documental-ui`; la app de demostración tiene estilos propios.

No se implementan carpetas, eliminación, versiones, compartición ni firma de personas. Ampliaciones requieren revisar contrato y modos de falla.
