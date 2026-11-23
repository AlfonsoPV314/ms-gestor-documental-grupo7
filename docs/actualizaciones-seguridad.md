# Actualizaciones de seguridad respecto del DAS

Preparadas el 23 de noviembre de 2026. Deben presentarse al ayudante/jefe de proyecto en una revisión de actualización, separada de nuevas funcionalidades. No son excepciones de Trivy ni una aprobación docente ya obtenida.

| Componente                           | DAS                              | Implementación                                           | Motivo observado                                                                                                                                                                 |
| ------------------------------------ | -------------------------------- | -------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Toolchain Go e imagen de compilación | 1.26.5                           | 1.26.11                                                  | Trivy detectó 12 HIGH corregibles en stdlib 1.26.5; incluye CVE-2026-111943, corregido en 1.26.10                                                                                |
| golang.org/x/text, transitiva        | 0.29.0 resuelta por dependencias | 0.41.0                                                   | CVE-2026-56851 y CVE-2026-56852                                                                                                                                                  |
| golang.org/x/sync, transitiva        | 0.17.0 resuelta                  | 0.22.0                                                   | Requerida por x/text 0.41.0                                                                                                                                                      |
| Ejecutable migrate                   | CLI Goose con todos sus drivers  | Binario propio con la librería Goose 3.24.1 y pgx 5.10.0 | La CLI original incorpora drivers innecesarios con vulnerabilidades en jwt/v4, grpc, x/crypto y x/net. Se conserva Goose y su formato de migraciones; solo se compila PostgreSQL |

Vite se actualiza de 8.1.5 a 8.1.6 por un HIGH corregible detectado por pnpm audit y avisos adicionales de su servidor de desarrollo. Vitest y coverage-v8 se actualizan juntos de 4.1.10 a 4.1.12 para corregir avisos de lectura de archivos e inyección en el mecanismo de mocks.

Los reportes iniciales y posteriores se conservan localmente en `.artifacts/`. El digest de golang:1.26.11-alpine3.23 se consultó al registro oficial. No se cambia la versión principal del stack.

Govulncheck detectó además GO-2026-6859 alcanzable en net/http de Go 1.26.10, corregido en 1.26.11. El toolchain final y todos los builds usan 1.26.11.
