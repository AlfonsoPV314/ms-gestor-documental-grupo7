# Preparación de entrega — referencia 4 de diciembre de 2026

Ambos hitos se entregan en el repositorio público de PINGESO-CIVIL. El correo
reciente del curso aclara que el 4 de diciembre es una fecha de referencia.
Ver [checklist actual de cierre](cierre-entrega.md) y [guía de demo](demo.md).

## Hito 1

Revisar las capas y los modos de falla. Ejecutar con `API_MOCKS=true`, verificar todos los endpoints del contrato y compartirlo con grupo 5. Preservar este estado mediante un commit revisado en main antes de finalizar la lógica real. Requiere pipeline verde de ese commit, tag protegido hito-1, release con OpenAPI adjunto e issue Entrega Hito 1. Las plantillas de issue están incluidas.

El modo mock no pretende crear retrospectivamente una entrega: aún faltan repositorio oficial, revisión y evidencia remota. No crear tags oficiales antes de tener el run correcto. Al usar el mismo contrato para ambos hitos se conserva compatibilidad.

## Hito 2

Ejecutar con API_MOCKS=false y proveedor real del grupo 6. Requiere PostgreSQL real, permisos por rol, auditoría, pruebas golden y cobertura E2E. Conservar evidencia del run fallido de mutación; la fuente protegida debe estar restaurada y el pipeline de la entrega debe quedar verde.

Quality Gate: cero bugs/vulnerabilidades nuevas; 100% de hotspots revisados; duplicación menor a 3%; cobertura de código nuevo al menos 60%. Trivy sin HIGH/CRITICAL corregibles y Gitleaks sin hallazgos. No se han aprobado excepciones. Revisar las actualizaciones de seguridad del toolchain con el ayudante.

Sobre el commit final de main: tag hito-2, release, OpenAPI y SBOM CycloneDX adjuntos, resultado oasdiff contra hito-1, enlaces a pipeline, SonarQube, golden y run rojo, e issue de entrega.

## Publicación

El repositorio público ya está en `PINGESO-CIVIL/gestor-documental` y el remoto
local fue actualizado. Los demás integrantes deben ejecutar
`git remote set-url origin https://github.com/PINGESO-CIVIL/gestor-documental.git`.
Los ayudantes configuran la protección de main y tags; verificar PR obligatorio,
aprobación y checks verdes. El historial se escaneó con Gitleaks sin hallazgos en
la validación del commit `eac48e0`; repetir el control para el commit de entrega.

SONAR_TOKEN y SONAR_HOST_URL van directamente en Secrets del repositorio. El runner de GitHub no puede acceder a localhost de este PC. SonarQube Community local no proporciona análisis nativo de PR/ramas: la instancia del curso debe soportar la modalidad requerida antes de considerar CI completo. No se omite el scan cuando faltan secrets.

Publicar solo imágenes que pasaron el pipeline, con SHA del commit, en GHCR: api y migrate. El frontend de desarrollo se construye y escanea, pero el entregable reutilizable es el paquete npm. Adjuntar SBOM y firma Cosign v3.1.4 a la publicación. El workflow de publicación se habilita al configurar organización, permisos y repositorio; no se ha publicado nada desde esta carpeta.
