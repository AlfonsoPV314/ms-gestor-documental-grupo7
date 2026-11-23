import { useCallback, useEffect, useRef, useState } from 'react'
import { listar, verificar, descargar, type Documento } from './documentos.api'
export function useDocumentos(baseUrl: string, revision: number) {
  const [documentos, setDocumentos] = useState<Documento[]>([]),
    [error, setError] = useState(''),
    [cargando, setCargando] = useState(true),
    [offset, setOffset] = useState(0),
    [seleccionado, setSeleccionado] = useState<Documento>(),
    [resultado, setResultado] = useState<string>(),
    [ocupado, setOcupado] = useState(false)
  const generacion = useRef(0)
  const recargar = useCallback(async () => {
    const actual = ++generacion.current
    setCargando(true)
    setError('')
    setSeleccionado(undefined)
    setResultado(undefined)
    try {
      const ds = await listar(baseUrl, offset)
      if (actual === generacion.current) setDocumentos(ds)
    } catch (e) {
      if (actual === generacion.current) {
        setDocumentos([])
        setError(e instanceof Error ? e.message : 'Error de consulta')
      }
    } finally {
      if (actual === generacion.current) setCargando(false)
    }
  }, [baseUrl, offset])
  useEffect(() => {
    void recargar()
    return () => {
      generacion.current++
    }
  }, [recargar, revision])
  async function comprobar(d: Documento) {
    setOcupado(true)
    setError('')
    try {
      const r = await verificar(baseUrl, d.id)
      setResultado(
        r.valida
          ? 'Firma válida: el contenido conserva su integridad'
          : 'Firma inválida: el documento requiere revisión',
      )
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Error de verificación')
    } finally {
      setOcupado(false)
    }
  }
  async function bajar(d: Documento) {
    setOcupado(true)
    setError('')
    try {
      await descargar(baseUrl, d)
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Error de descarga')
    } finally {
      setOcupado(false)
    }
  }
  return {
    documentos,
    error,
    cargando,
    offset,
    setOffset,
    seleccionado,
    seleccionar: (d: Documento) => {
      setSeleccionado(d)
      setResultado(undefined)
    },
    resultado,
    ocupado,
    comprobar,
    bajar,
    recargar,
  }
}
