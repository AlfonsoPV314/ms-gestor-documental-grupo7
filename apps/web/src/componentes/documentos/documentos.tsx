import { useDocumentos } from './use-documentos'
import { useId } from 'react'
export interface DocumentosProps {
  baseUrl: string
  revision?: number
}
export function Documentos({ baseUrl, revision = 0 }: DocumentosProps) {
  const s = useDocumentos(baseUrl, revision)
  const id = useId()
  return (
    <section className="documental-ui panel biblioteca" aria-labelledby={`${id}-titulo`}>
      <div className="cabecera-seccion">
        <div>
          <div className="etiqueta">TU ESPACIO</div>
          <h2 id={`${id}-titulo`}>Biblioteca documental</h2>
        </div>
        <button
          className="secundario"
          onClick={() => void s.recargar()}
          disabled={s.cargando}
        >
          Actualizar
        </button>
      </div>
      {s.error && (
        <p className="error" role="alert">
          {s.error}
        </p>
      )}
      {s.cargando ? (
        <p role="status">Consultando documentos…</p>
      ) : s.documentos.length === 0 ? (
        <div className="vacio">
          <span>▤</span>
          <h3>Tu biblioteca comienza aquí</h3>
          <p>Los archivos que guardes aparecerán en este espacio.</p>
        </div>
      ) : (
        <ul className="lista-documentos">
          {s.documentos.map((d) => (
            <li key={d.id}>
              <button className="fila-documento" onClick={() => s.seleccionar(d)}>
                <span className="tipo-archivo">
                  {d.tipo === 'application/pdf'
                    ? 'PDF'
                    : d.tipo.startsWith('image/')
                      ? 'IMG'
                      : 'TXT'}
                </span>
                <span>
                  <strong>{d.nombre}</strong>
                  <small>
                    {new Intl.DateTimeFormat('es-CL', { dateStyle: 'medium' }).format(
                      new Date(d.creado_en),
                    )}{' '}
                    · {(d.tamano / 1024).toFixed(1)} KiB
                  </small>
                </span>
                <span className="flecha">↗</span>
              </button>
            </li>
          ))}
        </ul>
      )}
      <nav className="paginacion" aria-label="Paginación">
        <button
          className="secundario"
          disabled={s.offset === 0 || s.cargando}
          onClick={() => s.setOffset(Math.max(0, s.offset - 20))}
        >
          Anterior
        </button>
        <span>Página {s.offset / 20 + 1}</span>
        <button
          className="secundario"
          disabled={s.documentos.length < 20 || s.cargando}
          onClick={() => s.setOffset(s.offset + 20)}
        >
          Siguiente
        </button>
      </nav>
      {s.seleccionado && (
        <aside className="detalle" aria-label="Detalle del documento">
          <h3>{s.seleccionado.nombre}</h3>
          <dl>
            <dt>Identificador</dt>
            <dd>{s.seleccionado.id}</dd>
            <dt>Huella SHA-256</dt>
            <dd>{s.seleccionado.sha256}</dd>
            <dt>Clave de firma</dt>
            <dd>{s.seleccionado.clave_id}</dd>
          </dl>
          <div className="acciones">
            <button onClick={() => void s.bajar(s.seleccionado!)} disabled={s.ocupado}>
              Descargar original
            </button>
            <button
              className="secundario"
              onClick={() => void s.comprobar(s.seleccionado!)}
              disabled={s.ocupado}
            >
              Verificar firma
            </button>
          </div>
          {s.resultado && (
            <p
              role="status"
              className={s.resultado.startsWith('Firma válida') ? 'exito' : 'error'}
            >
              {s.resultado}
            </p>
          )}
        </aside>
      )}
    </section>
  )
}
