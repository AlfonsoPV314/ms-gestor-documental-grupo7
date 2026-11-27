# Cierre de entrega — Grupo 7

Estado revisado el 27 de noviembre de 2026. El correo más reciente del curso
trata el 4 de diciembre como fecha de referencia. Repositorio público transferido
a `PINGESO-CIVIL/gestor-documental`. No confundir validación local, run
remoto aprobado y entrega oficial.

## Trabajo preparado

- Implementación documental con PostgreSQL, firma, auditoría y autorización.
- Contrato y tipos generados; Swagger local y modo de respuestas simuladas.
- Rama `fix/separar-pruebas-e2e` y PR #15 con las correcciones iniciales del CI.
- [Demo reproducible](demo.md), archivo ficticio y reportes separados para mocks,
  lógica real e integración JWT/JWKS del grupo 6.
- Contrato externo versionado con commit de origen en
  [contracts/externos](../contracts/externos/README.md).
- Pipeline que conserva cobertura, reportes de demo, escaneo Trivy y SBOM.
- Plantillas de issues y borradores de notas de release. Son borradores, no
  certificaciones de cumplimiento.

## Antes de crear cualquier tag oficial

- [ ] Revisar y fusionar el PR #15 una vez validados sus cambios.
- [ ] Configurar `SONAR_TOKEN` y `SONAR_HOST_URL` directamente en GitHub. Debe
      existir un servidor accesible desde el runner y compatible con el análisis.
      La instancia local no sustituye esa configuración.
- [ ] Comprobar `calidad`, `seguridad` e `imagen` verdes en el commit de `main`.
      No omitir SonarQube ni bajar umbrales para convertirlo en verde.
- [ ] Escanear el historial completo con Gitleaks y revisar datos confidenciales
      del cliente antes de transferir/publicar. El escaneo no sustituye esa revisión.
- [x] Repositorio público transferido a `PINGESO-CIVIL/gestor-documental`, verificado
      mediante la API de GitHub. El PR #15 se conservó abierto.
- [x] Remoto local actualizado a `https://github.com/PINGESO-CIVIL/gestor-documental.git`.
- [ ] Cada integrante actualiza su propio remoto a la URL nueva.
- [ ] Confirmar con los ayudantes protección de main y tags. No crear el tag
      sobre un PR sin fusionar.

## Hito 1

- [ ] Compartir la versión del contrato con el consumidor, grupo 5, y conservar
      el enlace o confirmación. No se envían correos automáticamente.
- [ ] Ejecutar la demo con mocks sobre el commit a entregar.
- [ ] Revisar arquitectura, capas, contratos externos y modos de falla.
- [ ] Tag `hito-1` sobre main con su run verde.
- [ ] Release con `openapi.yaml` adjunto, cambios y enlaces verificables.
- [ ] Issue Entrega Hito 1 con checklist marcada solo donde haya evidencia.

## Hito 2

- [ ] Acordar con grupo 6 los permisos `documentos:leer` y `documentos:operar`
      de los roles reales; probar sesión del host cliente. Los permisos extra del
      endpoint de desarrollo no sustituyen esa asignación real.
- [ ] Revisar el estado de persistencia del proveedor: el commit consultado de
      usuarios y roles aún usa memoria. No atribuirle una integración final aprobada.
- [ ] Demostración real, PostgreSQL, auditoría e integración de identidad.
- [ ] Cobertura Go desde E2E, golden, concurrencia/rechazos y ambas zonas horarias.
- [ ] Ejecución manual roja que falle exactamente al permitir acceso ajeno;
      conservar el enlace, y entregar la protección restaurada.
- [ ] Quality Gate del commit final: cobertura nueva >=60%, duplicación <3%,
      cero bugs/vulnerabilidades nuevas, hotspots revisados al 100%.
- [ ] Gitleaks sin hallazgos; Trivy sin HIGH/CRITICAL corregibles.
- [ ] oasdiff contra `hito-1` sin cambios incompatibles.
- [ ] Tag `hito-2`, release con OpenAPI y SBOM CycloneDX e issue de entrega.
- [ ] Revisar con el ayudante las desviaciones de versiones por seguridad.

## Evidencia que debe enlazarse

| Evidencia     | Dónde se obtiene                                                   |
| ------------- | ------------------------------------------------------------------ |
| Commit        | SHA del tag y `git status` limpio                                  |
| Pipeline      | Run del mismo SHA; no uno anterior ni de otra rama                 |
| Contrato      | `contracts/openapi.yaml`, regeneración sin diferencias y oasdiff   |
| Demo          | `.artifacts/demo-mocks.json`, `demo-real.json`, `demo-grupo6.json` |
| Cobertura     | `apps/api/coverage.out` y `apps/web/coverage/lcov.info`            |
| Calidad       | Proyecto SonarQube y análisis del SHA entregado                    |
| Seguridad     | Job seguridad y artifact `evidencia-seguridad`                     |
| Imágenes/SBOM | Job imagen y artifact `evidencia-imagenes`                         |
| Mutación      | Run rojo esperado y `.artifacts/mutacion.log`                      |

El CI sube artifacts aun si una prueba falla cuando es posible; tener un artifact
no significa haber aprobado el control. No publicar claves, tokens, cookies,
`.env` ni `.secrets`. Los reportes locales anteriores no prueban automáticamente
el estado de un nuevo commit.
