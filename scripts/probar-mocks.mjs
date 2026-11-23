import assert from 'node:assert/strict'
const base = process.env.MOCK_API_URL ?? 'http://localhost:5174'
const identidad = process.env.DEV_IDENTITY_URL ?? 'http://localhost:5172'
const sesion = await fetch(`${identidad}/identidad/sesion`, {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ usuario: 'operador' }),
}).then((r) => r.json())
const headers = { Authorization: `Bearer ${sesion.token}` }
const id = '11111111-1111-4111-8111-111111111111'
const lista = await fetch(`${base}/api/v1/documentos`, { headers })
assert.equal(lista.status, 200)
assert.equal((await lista.json()).documentos[0].id, id)
const doc = await fetch(`${base}/api/v1/documentos/${id}`, { headers })
assert.equal(doc.status, 200)
assert.equal((await doc.json()).nombre, 'ejemplo.txt')
const archivo = await fetch(`${base}/api/v1/documentos/${id}/contenido`, { headers })
assert.equal(archivo.status, 200)
assert.equal(await archivo.text(), 'ejemplo')
const firma = await fetch(`${base}/api/v1/documentos/${id}/verificacion`, { headers })
assert.equal(firma.status, 200)
assert.equal((await firma.json()).valida, true)
const form = new FormData()
form.append('archivo', new Blob(['ejemplo'], { type: 'text/plain' }), 'ejemplo.txt')
const carga = await fetch(`${base}/api/v1/documentos`, {
  method: 'POST',
  headers,
  body: form,
})
assert.equal(carga.status, 201)
assert.equal((await carga.json()).id, id)
console.log('Hito 1: los cinco endpoints responden con el fixture del contrato')
