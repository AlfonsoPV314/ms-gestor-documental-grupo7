# Modos de falla

| Falla                                              | Protección                                                        | Verificación                                                             |
| -------------------------------------------------- | ----------------------------------------------------------------- | ------------------------------------------------------------------------ |
| Acceso a documento ajeno                           | Permiso de rol y filtro por propietario en SQL                    | E2E con dos identidades: consulta, descarga y verificación responden 404 |
| Token ausente, vencido, firma o audiencia inválida | JWT con algoritmo permitido, issuer, audience, expiración y JWKS  | acceso_test.go y E2E con tokens firmados localmente                      |
| Solicitud mutante desde otro origen                | Validación de Origin contra FRONTEND_URL para sesión cookie       | E2E de rechazo 403                                                       |
| Exceso de tamaño                                   | Límite HTTP y lectura de archivo a 20 MiB + 1                     | E2E: 413, sin documento ni auditoría parcial                             |
| Nombre peligroso o metadatos inválidos             | Nombre normalizado y validación de entrada                        | E2E de ruta y tipo no admitidos                                          |
| Escritura parcial                                  | Documento y auditoría en la misma transacción                     | E2E con fallo de auditoría y rollback                                    |
| Contenido alterado                                 | SHA-256 y firma Ed25519 ligados a id y propietario                | E2E modifica bytes y detecta verificación inválida                       |
| Clave perdida                                      | Volumen privado persistente; arranque falla si falta              | Reinicio y verificación del documento previo                             |
| PostgreSQL no disponible                           | Contextos y timeout; readiness 503                                | E2E y Compose                                                            |
| Migración fallida                                  | API depende de migrate exitoso                                    | Validación del orden de Compose                                          |
| Dependencia de identidad caída                     | Caché JWKS; sin llave válida se deniega acceso                    | Prueba de token con kid desconocido                                      |
| Datos sensibles en logs                            | Solo método, patrón de ruta, estado, duración e id de correlación | Inspección de logs E2E                                                   |
| Cambio incompatible del contrato                   | codegen reproducible y oasdiff                                    | CI contra main y contra hito-1                                           |

La evidencia de mutación se obtiene con scripts/mutacion.mjs: retira temporalmente el filtro de propietario, exige que el E2E falle y restaura la fuente en finally. Nunca se publica una imagen con esa modificación.
