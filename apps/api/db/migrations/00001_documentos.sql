-- +goose Up
CREATE TABLE documentos (
 id uuid PRIMARY KEY,
 propietario uuid NOT NULL,
 nombre text NOT NULL CHECK (length(nombre) BETWEEN 1 AND 200),
 tipo text NOT NULL,
 tamano bigint NOT NULL CHECK (tamano BETWEEN 1 AND 20971520),
 contenido bytea NOT NULL,
 hash text NOT NULL CHECK (length(hash)=64),
 firma text NOT NULL,
 clave_id text NOT NULL,
 creado_en timestamptz NOT NULL DEFAULT now(),
 CHECK (octet_length(contenido)=tamano)
);
CREATE INDEX documentos_propietario_fecha ON documentos(propietario,creado_en DESC,id);
CREATE TABLE auditoria (
 id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
 documento_id uuid NOT NULL REFERENCES documentos(id),
 actor uuid NOT NULL,
 accion text NOT NULL,
 request_id text NOT NULL,
 creado_en timestamptz NOT NULL DEFAULT now()
);
-- +goose Down
DROP TABLE auditoria;
DROP TABLE documentos;
