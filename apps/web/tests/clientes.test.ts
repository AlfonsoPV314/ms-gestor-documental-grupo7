import { afterAll, afterEach, beforeAll, describe, expect, it } from 'vitest'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import { listar, verificar } from '../src/componentes/documentos/documentos.api'
import { cargarArchivo } from '../src/componentes/carga-documental/carga-documental.api'
const base = 'http://documental.test'
const server = setupServer()
beforeAll(() => server.listen({ onUnhandledRequest: 'error' }))
afterEach(() => server.resetHandlers())
afterAll(() => server.close())
describe('clientes del contrato', () => {
  it('lista documentos y envía paginación', async () => {
    server.use(
      http.get(`${base}/api/v1/documentos`, ({ request }) => {
        expect(new URL(request.url).searchParams.get('offset')).toBe('20')
        return HttpResponse.json({ documentos: [] })
      }),
    )
    expect(await listar(base, 20)).toEqual([])
  })
  it('expone rechazo de permisos', async () => {
    server.use(
      http.get(`${base}/api/v1/documentos`, () =>
        HttpResponse.json({ detail: 'Permiso insuficiente' }, { status: 403 }),
      ),
    )
    await expect(listar(base, 0)).rejects.toThrow('Permiso insuficiente')
  })
  it('verifica un documento', async () => {
    server.use(
      http.get(`${base}/api/v1/documentos/:id/verificacion`, () =>
        HttpResponse.json({ valida: true, algoritmo: 'Ed25519', clave_id: 'prueba' }),
      ),
    )
    expect((await verificar(base, 'id')).valida).toBe(true)
  })
  it('sube el binario en multipart', async () => {
    server.use(
      http.post(`${base}/api/v1/documentos`, async ({ request }) => {
        const form = await request.formData()
        expect((form.get('archivo') as File).name).toBe('informe.txt')
        return HttpResponse.json({ id: 'prueba' }, { status: 201 })
      }),
    )
    expect(
      (
        await cargarArchivo(
          base,
          new File(['texto'], 'informe.txt', { type: 'text/plain' }),
        )
      ).id,
    ).toBe('prueba')
  })
  it('expone un rechazo de carga', async () => {
    server.use(
      http.post(`${base}/api/v1/documentos`, () =>
        HttpResponse.json({ detail: 'Máximo 20 MiB' }, { status: 413 }),
      ),
    )
    await expect(cargarArchivo(base, new File(['x'], 'informe.txt'))).rejects.toThrow(
      'Máximo 20 MiB',
    )
  })
})
