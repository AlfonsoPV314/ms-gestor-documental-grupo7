package store

import (
	"context"
	"errors"
	"github.com/PINGESO-CIVIL/ms-gestor-documental/apps/api/internal/repository"
	"github.com/PINGESO-CIVIL/ms-gestor-documental/apps/api/internal/service"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct{ Pool *pgxpool.Pool }

func uuid(s string) pgtype.UUID { var u pgtype.UUID; _ = u.Scan(s); return u }
func texto(u pgtype.UUID) string {
	v, _ := u.Value()
	if v == nil {
		return ""
	}
	return v.(string)
}
func (p *Postgres) Crear(ctx context.Context, d service.Documento, requestID string) error {
	tx, err := p.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	q := repository.New(tx)
	err = q.CrearDocumento(ctx, repository.CrearDocumentoParams{ID: uuid(d.ID), Propietario: uuid(d.Propietario), Nombre: d.Nombre, Tipo: d.Tipo, Tamano: d.Tamano, Contenido: d.Contenido, Hash: d.Hash, Firma: d.Firma, ClaveID: d.ClaveID, CreadoEn: pgtype.Timestamptz{Time: d.CreadoEn, Valid: true}})
	if err != nil {
		return err
	}
	if err = q.AuditarCarga(ctx, repository.AuditarCargaParams{DocumentoID: uuid(d.ID), Actor: uuid(d.Propietario), RequestID: requestID}); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
func (p *Postgres) Obtener(ctx context.Context, id, propietario string) (service.Documento, error) {
	d, err := repository.New(p.Pool).ObtenerDocumento(ctx, repository.ObtenerDocumentoParams{ID: uuid(id), Propietario: uuid(propietario)})
	if errors.Is(err, pgx.ErrNoRows) {
		return service.Documento{}, service.ErrNoEncontrado
	}
	if err != nil {
		return service.Documento{}, err
	}
	return service.Documento{ID: texto(d.ID), Propietario: texto(d.Propietario), Nombre: d.Nombre, Tipo: d.Tipo, Tamano: d.Tamano, Contenido: d.Contenido, Hash: d.Hash, Firma: d.Firma, ClaveID: d.ClaveID, CreadoEn: d.CreadoEn.Time}, nil
}
func (p *Postgres) Listar(ctx context.Context, propietario string, limite, offset int32) ([]service.Documento, error) {
	filas, err := repository.New(p.Pool).ListarDocumentos(ctx, repository.ListarDocumentosParams{Propietario: uuid(propietario), Limit: limite, Offset: offset})
	if err != nil {
		return nil, err
	}
	docs := make([]service.Documento, 0, len(filas))
	for _, d := range filas {
		docs = append(docs, service.Documento{ID: texto(d.ID), Propietario: texto(d.Propietario), Nombre: d.Nombre, Tipo: d.Tipo, Tamano: d.Tamano, Hash: d.Hash, Firma: d.Firma, ClaveID: d.ClaveID, CreadoEn: d.CreadoEn.Time})
	}
	return docs, nil
}
