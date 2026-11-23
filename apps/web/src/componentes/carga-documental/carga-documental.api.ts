import type { components } from '../../contrato'
export type Documento = components['schemas']['Documento']
export async function cargarArchivo(baseUrl: string, archivo: File): Promise<Documento> {
  const form = new FormData()
  form.append('archivo', archivo)
  const r = await fetch(`${baseUrl}/api/v1/documentos`, {
    method: 'POST',
    body: form,
    credentials: 'include',
  })
  if (!r.ok) {
    const p = await r.json().catch(() => ({}))
    throw Error(p.detail ?? 'No se pudo cargar el archivo')
  }
  return r.json()
}
