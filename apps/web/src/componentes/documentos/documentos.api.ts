import createClient from 'openapi-fetch'
import type { paths, components } from '../../contrato'
export type Documento = components['schemas']['Documento']
function cliente(baseUrl: string) {
  return createClient<paths>({ baseUrl, credentials: 'include' })
}
export async function listar(baseUrl: string, offset: number) {
  const { data, error } = await cliente(baseUrl).GET('/api/v1/documentos', {
    params: { query: { limite: 20, offset } },
  })
  if (error) throw Error(error.detail ?? 'No se pudieron consultar los documentos')
  return data.documentos ?? []
}
export async function verificar(baseUrl: string, id: string) {
  const { data, error } = await cliente(baseUrl).GET(
    '/api/v1/documentos/{id}/verificacion',
    { params: { path: { id } } },
  )
  if (error) throw Error(error.detail ?? 'No se pudo verificar')
  return data
}
export async function descargar(baseUrl: string, d: Documento) {
  const r = await fetch(`${baseUrl}/api/v1/documentos/${d.id}/contenido`, {
    credentials: 'include',
  })
  if (!r.ok) throw Error('No se pudo descargar el documento')
  const u = URL.createObjectURL(await r.blob())
  const a = document.createElement('a')
  a.href = u
  a.download = d.nombre
  a.click()
  setTimeout(() => URL.revokeObjectURL(u), 1000)
}
