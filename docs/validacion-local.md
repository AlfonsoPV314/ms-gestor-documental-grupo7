# Validación local — 23 de noviembre de 2026

Implementación nueva en este repositorio. La implementación anterior de la carpeta superior no se reutilizó ni se modificó. No se realizó push, publicación, despliegue cloud ni creación de tags oficiales.

| Verificación                                           | Resultado observado                                                                                                                          |
| ------------------------------------------------------ | -------------------------------------------------------------------------------------------------------------------------------------------- |
| Build API, frontend, migraciones e identidad de prueba | Correcto con Docker Linux                                                                                                                    |
| Herramientas locales                                   | Node 22.22.0, pnpm 10.13.1, Go 1.26.11 y Gitleaks 8.30.1                                                                                     |
| Hooks                                                  | Instalados; pre-commit del commit inicial pasó formato, arquitectura, Go y secretos                                                          |
| Prueba de secreto sintético                            | Hook rechazó el commit en un repositorio temporal aislado                                                                                    |
| Contrato                                               | Huma genera OpenAPI 3.1 y tipos; regeneración sin diferencias                                                                                |
| Swagger UI                                             | Recursos locales; siete operaciones visibles y consulta autenticada 200 comprobadas en Chromium; rutas deshabilitadas en producción probadas |
| sqlc                                                   | Consultas generadas reproducibles                                                                                                            |
| Modo Hito 1                                            | Cinco endpoints comprobados con fixture simulado en contenedor temporal                                                                      |
| E2E API                                                | Binario real y PostgreSQL real; golden aprobado en UTC y Pacific/Auckland                                                                    |
| Escenarios                                             | JWT, roles, aislamiento, tamaño, firma, manipulación, concurrencia, rollback y reinicio                                                      |
| Mutación                                               | Retirar filtro de propietario provoca el fallo esperado; fuente restaurada                                                                   |
| Pruebas frontend                                       | 13 aprobadas; 93,93 % de líneas en microcomponentes                                                                                          |
| Chromium                                               | Flujo de carga, firma, descarga, lector y usuario ajeno aprobado                                                                             |
| Revisión visual                                        | Captura completa de la aplicación inspeccionada                                                                                              |
| SonarQube local                                        | Quality Gate PASSED, cobertura nueva 66,6 %, bugs nuevos 0, vulnerabilidades nuevas 0, duplicación nueva 0 %                                 |
| Hotspots                                               | No se detectaron hotspots pendientes en el análisis local                                                                                    |
| pnpm audit                                             | Cero vulnerabilidades de todas las severidades                                                                                               |
| govulncheck                                            | No vulnerabilities found con Go 1.26.11                                                                                                      |
| Trivy repositorio e imágenes                           | Cero HIGH/CRITICAL corregibles; sin excepciones                                                                                              |
| Hadolint                                               | Sin errores al umbral requerido; observaciones informativas no bloqueantes                                                                   |
| Dockle API, web y migrate                              | Sin WARN/FATAL; informa ausencia de HEALTHCHECK en imágenes (Compose sí comprueba salud) y content trust                                     |
| SBOM                                                   | CycloneDX de API, web y migrate generado localmente                                                                                          |
| Paquete UI                                             | Tarball solo con librería, estilos, declaraciones y contrato; sin la app de demostración                                                     |

Los reportes reproducibles se conservan en `.artifacts/`, ignorado por Git. La cobertura Go sale del binario instrumentado, no de un porcentaje inventado. SonarQube analiza también los comandos del backend; no se omitieron para elevar la cobertura.

El código funcional que obtuvo el gate verde está en el commit local inicial `576e348`; cambios posteriores de documentación, empaquetado o evidencias deben conservar los controles y repetir el análisis para su propio commit de entrega.

## Pendientes para la entrega oficial

- Repositorio público de PINGESO-CIVIL, permisos, protección de main y tags.
- Proyecto/instancia SonarQube accesible desde CI y secrets cargados directamente allí.
- Contrato real del grupo 6 y verificación de sus claims, roles, cookie y JWKS. El proveedor local es un fixture explícito.
- Revisión del curso de los parches de seguridad documentados respecto de las versiones antiguas del DAS.
- Runs remotos, evidencia roja remota, tags hito-1/hito-2, releases e issues oficiales.
- Publicación de imágenes y paquete, firma Cosign y despliegues posteriores. No se habilitó un workflow con permisos de publicación.

La evidencia local no se presenta como aprobación del ayudante ni como cumplimiento de dependencias externas aún no verificadas.
