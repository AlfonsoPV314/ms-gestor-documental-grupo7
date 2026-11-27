# Contrato de identidad — Grupo 6

`identidad.yaml` es una copia sin modificaciones del contrato publicado en
<https://github.com/PINGESO-CIVIL/usuarios-roles>, commit
`2d5ecccfd0fd27415ffbd4d6dd017ed895d94fd1`, consultado el 27 de noviembre de 2026.

Se conserva el archivo original para poder comparar futuras versiones. Nuestro
servicio consume el JWKS con keyfunc; no necesita consultar usuarios, roles ni
sesiones por cada documento, ni acceder a la base de identidad. Los endpoints
de desarrollo se utilizan solamente en la prueba de integración.

Valores del proveedor: algoritmo RS256, issuer `pingeso-identidad`, audience
`pingeso`, cookie `sesion`, `sub` UUID y `permisos` como arreglo de strings. El
campo `roles` también es un arreglo; la autorización documental se decide por
los permisos firmados. Los roles seed no tienen permisos documentales: hay que
acordar su asignación real; `permisos_extra` solo sirve para desarrollo.

La versión consultada almacena usuarios en memoria. Probar sus JWT reales no
acredita una identidad persistente de producción. Ver [guía de integración](../../docs/demo.md).
