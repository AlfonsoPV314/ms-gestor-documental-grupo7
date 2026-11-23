import { useId, useState } from 'react'
import { useCargaDocumental } from './use-carga-documental'
import type { Documento } from './carga-documental.api'
export interface CargaDocumentalProps {
  baseUrl: string
  puedeCargar: boolean
  alCargar?: (d: Documento) => void
}
export function CargaDocumental({
  baseUrl,
  puedeCargar,
  alCargar,
}: CargaDocumentalProps) {
  const [archivo, setArchivo] = useState<File>()
  const id = useId()
  const { ocupado, mensaje, error, cargar } = useCargaDocumental(baseUrl, alCargar)
  return (
    <section className="documental-ui panel carga" aria-labelledby={`${id}-titulo`}>
      <div className="etiqueta">INGRESO SEGURO</div>
      <h2 id={`${id}-titulo`}>Un nuevo documento</h2>
      <p>Conserva el original y comprueba su integridad cuando lo necesites.</p>
      <form
        onSubmit={(e) => {
          e.preventDefault()
          void cargar(archivo)
        }}
      >
        <label className="zona-archivo" htmlFor={`${id}-archivo`}>
          <span className="icono-carga">↑</span>
          <strong>{archivo?.name ?? 'Selecciona tu archivo'}</strong>
          <span>PDF, PNG, JPEG o texto · hasta 20 MiB</span>
          <input
            id={`${id}-archivo`}
            aria-label="Archivo para cargar"
            type="file"
            accept=".pdf,.png,.jpg,.jpeg,.txt"
            disabled={!puedeCargar || ocupado}
            onChange={(e) => setArchivo(e.target.files?.[0])}
          />
        </label>
        <button type="submit" disabled={!puedeCargar || ocupado || !archivo}>
          {ocupado ? 'Guardando…' : 'Guardar y firmar'}
        </button>
      </form>
      {!puedeCargar && (
        <p className="nota">Tu rol permite consultar y descargar documentos.</p>
      )}
      {error && (
        <p role="alert" className="error">
          {error}
        </p>
      )}
      {mensaje && (
        <p role="status" className="exito">
          {mensaje}
        </p>
      )}
    </section>
  )
}
