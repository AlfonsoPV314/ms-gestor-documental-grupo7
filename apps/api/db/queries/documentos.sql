-- name: CrearDocumento :exec
INSERT INTO documentos(id,propietario,nombre,tipo,tamano,contenido,hash,firma,clave_id,creado_en)
VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10);

-- name: AuditarCarga :exec
INSERT INTO auditoria(documento_id,actor,accion,request_id) VALUES($1,$2,'documento.creado',$3);

-- name: ObtenerDocumento :one
SELECT * FROM documentos WHERE id=$1 AND propietario=$2;

-- name: ListarDocumentos :many
SELECT id,propietario,nombre,tipo,tamano,hash,firma,clave_id,creado_en FROM documentos
WHERE propietario=$1 ORDER BY creado_en DESC,id LIMIT $2 OFFSET $3;
