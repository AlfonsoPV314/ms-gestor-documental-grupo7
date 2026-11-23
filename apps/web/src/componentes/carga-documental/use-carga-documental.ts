import { useState } from 'react'
import { cargarArchivo, type Documento } from './carga-documental.api'
export function useCargaDocumental(baseUrl: string, alCargar?: (d: Documento) => void) {
  const [ocupado, setOcupado] = useState(false),
    [mensaje, setMensaje] = useState(''),
    [error, setError] = useState('')
  async function cargar(archivo?: File) {
    if (!archivo) {
      setError('Selecciona un archivo')
      return
    }
    if (archivo.size > 20 * 1024 * 1024) {
      setError('El archivo supera los 20 MiB')
      return
    }
    setOcupado(true)
    setError('')
    setMensaje('')
    try {
      const d = await cargarArchivo(baseUrl, archivo)
      setMensaje(`«${d.nombre}» se guardó y firmó correctamente`)
      alCargar?.(d)
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Error de carga')
    } finally {
      setOcupado(false)
    }
  }
  return { ocupado, mensaje, error, cargar }
}
