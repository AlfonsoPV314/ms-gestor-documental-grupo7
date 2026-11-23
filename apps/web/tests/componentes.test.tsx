// @vitest-environment jsdom
import { afterAll, afterEach, beforeAll, expect, it, vi } from 'vitest'
import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import { CargaDocumental } from '../src/componentes/carga-documental/carga-documental'
import { Documentos } from '../src/componentes/documentos/documentos'
import { cargarArchivo } from '../src/componentes/carga-documental/carga-documental.api'
vi.mock(
  '../src/componentes/carga-documental/carga-documental.api',
  async (importOriginal) => ({
    ...(await importOriginal<object>()),
    cargarArchivo: vi.fn(),
  }),
)
const base = 'http://documental.test'
const documento = {
  id: '11111111-1111-4111-8111-111111111111',
  propietario: '22222222-2222-4222-8222-222222222222',
  nombre: 'informe.txt',
  tipo: 'text/plain',
  tamano: 10,
  sha256: 'a'.repeat(64),
  firma: 'prueba',
  clave_id: 'prueba',
  creado_en: '2026-11-23T00:00:00Z',
}
const server = setupServer()
beforeAll(() => server.listen({ onUnhandledRequest: 'error' }))
afterEach(() => {
  cleanup()
  server.resetHandlers()
  vi.restoreAllMocks()
})
afterAll(() => server.close())
it('un lector ve la carga deshabilitada', () => {
  render(<CargaDocumental baseUrl={base} puedeCargar={false} />)
  expect(
    (screen.getByRole('button', { name: 'Guardar y firmar' }) as HTMLButtonElement)
      .disabled,
  ).toBe(true)
  expect(
    screen.getByText('Tu rol permite consultar y descargar documentos.'),
  ).toBeTruthy()
})
it('carga un archivo y notifica al contenedor', async () => {
  vi.mocked(cargarArchivo).mockResolvedValue(documento)
  const alCargar = vi.fn()
  render(<CargaDocumental baseUrl={base} puedeCargar alCargar={alCargar} />)
  fireEvent.change(screen.getByLabelText('Archivo para cargar'), {
    target: { files: [new File(['texto'], 'informe.txt', { type: 'text/plain' })] },
  })
  fireEvent.click(screen.getByRole('button', { name: 'Guardar y firmar' }))
  await screen.findByText(/se guardó y firmó correctamente/)
  expect(alCargar).toHaveBeenCalledOnce()
})
it('rechaza tamaño excesivo antes de enviar', async () => {
  render(<CargaDocumental baseUrl={base} puedeCargar />)
  const archivo = new File(['x'], 'grande.txt')
  Object.defineProperty(archivo, 'size', { value: 21 * 1024 * 1024 })
  fireEvent.change(screen.getByLabelText('Archivo para cargar'), {
    target: { files: [archivo] },
  })
  fireEvent.click(screen.getByRole('button', { name: 'Guardar y firmar' }))
  expect(await screen.findByText('El archivo supera los 20 MiB')).toBeTruthy()
})
it('muestra un rechazo de la API', async () => {
  vi.mocked(cargarArchivo).mockRejectedValue(Error('Sesión vencida'))
  render(<CargaDocumental baseUrl={base} puedeCargar />)
  fireEvent.change(screen.getByLabelText('Archivo para cargar'), {
    target: { files: [new File(['texto'], 'informe.txt')] },
  })
  fireEvent.click(screen.getByRole('button', { name: 'Guardar y firmar' }))
  await screen.findByRole('alert')
  expect(screen.getByText('Sesión vencida')).toBeTruthy()
})
it('muestra biblioteca vacía y permite actualizar', async () => {
  server.use(
    http.get(`${base}/api/v1/documentos`, () => HttpResponse.json({ documentos: [] })),
  )
  render(<Documentos baseUrl={base} />)
  expect(await screen.findByText('Tu biblioteca comienza aquí')).toBeTruthy()
  fireEvent.click(screen.getByRole('button', { name: 'Actualizar' }))
  expect(await screen.findByText('Tu biblioteca comienza aquí')).toBeTruthy()
})
it('consulta detalle, verifica firma y descarga', async () => {
  server.use(
    http.get(`${base}/api/v1/documentos`, () =>
      HttpResponse.json({ documentos: [documento] }),
    ),
    http.get(`${base}/api/v1/documentos/:id/verificacion`, () =>
      HttpResponse.json({ valida: true, algoritmo: 'Ed25519', clave_id: 'prueba' }),
    ),
    http.get(`${base}/api/v1/documentos/:id/contenido`, () => new HttpResponse('texto')),
  )
  Object.defineProperty(URL, 'createObjectURL', {
    configurable: true,
    value: vi.fn(() => 'blob:prueba'),
  })
  Object.defineProperty(URL, 'revokeObjectURL', { configurable: true, value: vi.fn() })
  vi.spyOn(HTMLAnchorElement.prototype, 'click').mockImplementation(() => {})
  render(<Documentos baseUrl={base} />)
  fireEvent.click(await screen.findByRole('button', { name: /informe.txt/ }))
  fireEvent.click(screen.getByRole('button', { name: 'Verificar firma' }))
  await screen.findByText('Firma válida: el contenido conserva su integridad')
  fireEvent.click(screen.getByRole('button', { name: 'Descargar original' }))
  await waitFor(() => expect(URL.createObjectURL).toHaveBeenCalled())
})
it('muestra fallos de consulta', async () => {
  server.use(
    http.get(`${base}/api/v1/documentos`, () =>
      HttpResponse.json({ detail: 'Sin acceso' }, { status: 403 }),
    ),
  )
  render(<Documentos baseUrl={base} />)
  expect(await screen.findByText('Sin acceso')).toBeTruthy()
})
it('muestra una firma inválida', async () => {
  server.use(
    http.get(`${base}/api/v1/documentos`, () =>
      HttpResponse.json({ documentos: [documento] }),
    ),
    http.get(`${base}/api/v1/documentos/:id/verificacion`, () =>
      HttpResponse.json({ valida: false, algoritmo: 'Ed25519', clave_id: 'prueba' }),
    ),
  )
  render(<Documentos baseUrl={base} />)
  fireEvent.click(await screen.findByRole('button', { name: /informe.txt/ }))
  fireEvent.click(screen.getByRole('button', { name: 'Verificar firma' }))
  expect(
    await screen.findByText('Firma inválida: el documento requiere revisión'),
  ).toBeTruthy()
})
