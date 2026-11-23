import { existsSync, mkdirSync, readFileSync, writeFileSync } from 'node:fs'
import { randomBytes } from 'node:crypto'
const base = process.env.SONAR_HOST_URL ?? 'http://localhost:9000'
const archivo = new URL('../.secrets/sonar.json', import.meta.url)
let secreto = existsSync(archivo)
  ? JSON.parse(readFileSync(archivo, 'utf8'))
  : { usuario: 'admin', password: 'G7!aA1' + randomBytes(24).toString('hex') }
async function api(ruta, datos, pass = secreto.password) {
  const r = await fetch(`${base}/api/${ruta}`, {
    method: datos ? 'POST' : 'GET',
    headers: {
      Authorization: `Basic ${Buffer.from(`admin:${pass}`).toString('base64')}`,
      ...(datos ? { 'Content-Type': 'application/x-www-form-urlencoded' } : {}),
    },
    body: datos ? new URLSearchParams(datos) : undefined,
  })
  const text = await r.text()
  if (!r.ok) throw Error(`${ruta}: HTTP ${r.status}: ${text}`)
  return text ? JSON.parse(text) : {}
}
const estado = await fetch(`${base}/api/system/status`).then((r) => r.json())
if (estado.status !== 'UP') throw Error(`SonarQube aún no está UP: ${estado.status}`)
if (!existsSync(archivo)) {
  await api(
    'users/change_password',
    { login: 'admin', previousPassword: 'admin', password: secreto.password },
    'admin',
  )
  mkdirSync(new URL('../.secrets/', import.meta.url), { recursive: true })
  writeFileSync(archivo, JSON.stringify(secreto), { mode: 0o600 })
}
const proyectos = await api('projects/search?projects=ms-gestor-documental')
if (!proyectos.components.length)
  await api('projects/create', {
    project: 'ms-gestor-documental',
    name: 'MS Gestor Documental - Grupo 7',
    visibility: 'private',
  })
const gates = await api('qualitygates/list')
let gate = gates.qualitygates.find((g) => g.name === 'DAS Grupo 7')
if (!gate) gate = await api('qualitygates/create', { name: 'DAS Grupo 7' })
const detalle = await api(`qualitygates/show?name=${encodeURIComponent('DAS Grupo 7')}`)
for (const [metric, op, error] of [
  ['new_bugs', 'GT', '0'],
  ['new_vulnerabilities', 'GT', '0'],
  ['new_security_hotspots_reviewed', 'LT', '100'],
  ['new_duplicated_lines_density', 'GT', '2.99'],
  ['new_coverage', 'LT', '60'],
]) {
  const condicion = detalle.conditions?.find((c) => c.metric === metric)
  if (condicion)
    await api('qualitygates/update_condition', {
      id: String(condicion.id),
      metric,
      op,
      error,
    })
  else
    await api('qualitygates/create_condition', {
      gateName: 'DAS Grupo 7',
      metric,
      op,
      error,
    })
}
await api('qualitygates/select', {
  projectKey: 'ms-gestor-documental',
  gateName: 'DAS Grupo 7',
})
await api('new_code_periods/set', {
  project: 'ms-gestor-documental',
  type: 'PREVIOUS_VERSION',
})
await api('settings/set', {
  component: 'ms-gestor-documental',
  key: 'sonar.qualitygate.ignoreSmallChanges',
  value: 'false',
})
if (!secreto.token) {
  const t = await api('user_tokens/generate', {
    name: `documental-local-${Date.now()}`,
    type: 'PROJECT_ANALYSIS_TOKEN',
    projectKey: 'ms-gestor-documental',
  })
  secreto.token = t.token
  writeFileSync(archivo, JSON.stringify(secreto), { mode: 0o600 })
}
console.log(
  'Proyecto, código nuevo y Quality Gate DAS configurados. Credenciales únicamente en .secrets/sonar.json',
)
