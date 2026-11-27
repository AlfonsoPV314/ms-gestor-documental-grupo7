import assert from 'node:assert/strict'
import { readFileSync, mkdirSync, writeFileSync } from 'node:fs'
import { execFileSync } from 'node:child_process'
import { createHash } from 'node:crypto'

const raiz = new URL('../', import.meta.url)
const mocks = process.argv.includes('--mocks')
const grupo6 = process.argv.includes('--grupo6')
assert(!(mocks && grupo6), 'Separar demo de mocks y prueba con grupo 6')
const modo = grupo6 ? 'grupo6' : mocks ? 'mocks' : 'real'
const base = process.env.DEMO_API_URL ?? 'http://localhost:5171'
const proveedor = process.env.DEMO_IDENTIDAD_URL ?? 'http://localhost:5161'
const archivo = readFileSync(new URL('ejemplos/documento.txt', raiz))
const resultados = []
async function llamada(ruta, esperado, opciones = {}, origen = base) {
  const respuesta = await fetch(new URL(ruta, origen), {
    ...opciones,
    signal: AbortSignal.timeout(15000),
  })
  assert.equal(
    respuesta.status,
    esperado,
    `${ruta}: HTTP ${respuesta.status}, esperado ${esperado}`,
  )
  resultados.push({ metodo: opciones.method ?? 'GET', ruta, estado: respuesta.status })
  return respuesta
}
async function token(usuario, permisos, id) {
  if (grupo6) {
    const respuesta = await llamada(
      '/api/dev/tokens',
      200,
      {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          roles: ['operador'],
          usuario_id: id,
          permisos_extra: permisos,
          duracion: '15m',
        }),
      },
      proveedor,
    )
    return (await respuesta.json()).token
  }
  const respuesta = await llamada(
    '/identidad/sesion',
    200,
    {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ usuario }),
    },
    process.env.DEMO_FIXTURE_URL ?? 'http://localhost:5172',
  )
  return (await respuesta.json()).token
}
const propietario = '22222222-2222-4222-8222-222222222222'
const ajenoID = '33333333-3333-4333-8333-333333333333'
const operador = await token(
  'operador',
  ['documentos:leer', 'documentos:operar'],
  propietario,
)
const lector = await token('lector', ['documentos:leer'], propietario)
const ajeno = await token('ajeno', ['documentos:leer', 'documentos:operar'], ajenoID)
const headers = (jwt) => ({ Authorization: `Bearer ${jwt}` })
const formulario = () => {
  const form = new FormData()
  form.append('archivo', new Blob([archivo], { type: 'text/plain' }), 'documento.txt')
  return form
}
await llamada('/healthz', 200)
await llamada('/readyz', 200)
await llamada('/api/docs', 200)
const spec = await (await llamada('/api/openapi.json', 200)).json()
assert.equal(
  Object.values(spec.paths).flatMap((ruta) =>
    Object.keys(ruta).filter((metodo) => ['get', 'post'].includes(metodo)),
  ).length,
  7,
)
await llamada('/api/v1/documentos', 401)
const documento = await (
  await llamada('/api/v1/documentos', 201, {
    method: 'POST',
    headers: headers(operador),
    body: formulario(),
  })
).json()
assert.match(documento.id, /^[0-9a-f-]{36}$/)
const listado = await (
  await llamada('/api/v1/documentos?limite=100', 200, { headers: headers(lector) })
).json()
assert(
  listado.documentos.some((d) => d.id === documento.id),
  'Documento visible para su propietario',
)
const ruta = `/api/v1/documentos/${documento.id}`
const metadatos = await (await llamada(ruta, 200, { headers: headers(lector) })).json()
assert.equal(metadatos.id, documento.id)
const contenido = Buffer.from(
  await (
    await llamada(`${ruta}/contenido`, 200, { headers: headers(lector) })
  ).arrayBuffer(),
)
const verificacion = await (
  await llamada(`${ruta}/verificacion`, 200, { headers: headers(lector) })
).json()
assert.equal(verificacion.valida, true)
if (!mocks) {
  assert.deepEqual(contenido, archivo, 'Descarga idéntica al original')
  assert.equal(documento.sha256, createHash('sha256').update(archivo).digest('hex'))
  await llamada(ruta, 404, { headers: headers(ajeno) })
}
await llamada('/api/v1/documentos', 403, {
  method: 'POST',
  headers: headers(lector),
  body: formulario(),
})
mkdirSync(new URL('.artifacts/', raiz), { recursive: true })
const reporte = {
  commit: execFileSync('git', ['rev-parse', 'HEAD'], {
    cwd: raiz,
    encoding: 'utf8',
  }).trim(),
  cambios_pendientes:
    execFileSync('git', ['status', '--porcelain'], {
      cwd: raiz,
      encoding: 'utf8',
    }).trim() !== '',
  modo,
  fecha: new Date().toISOString(),
  documento_id: documento.id,
  archivo_sha256: createHash('sha256').update(archivo).digest('hex'),
  resultados,
  alcance: grupo6
    ? 'JWT/JWKS reales del grupo 6; permisos_extra es exclusivo de desarrollo'
    : mocks
      ? 'Contrato con respuestas simuladas; no demuestra persistencia'
      : 'Lógica documental y PostgreSQL reales; identidad fixture local',
}
writeFileSync(
  new URL(`.artifacts/demo-${modo}.json`, raiz),
  JSON.stringify(reporte, null, 2) + '\n',
)
console.log(
  `Demo ${modo}: ${resultados.length} respuestas esperadas; carga, consulta, descarga y firma correctas. Reporte sin tokens: .artifacts/demo-${modo}.json`,
)
