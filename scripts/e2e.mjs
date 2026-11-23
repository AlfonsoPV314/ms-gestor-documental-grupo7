import { spawnSync } from 'node:child_process'
import { mkdirSync, rmSync } from 'node:fs'
import { fileURLToPath } from 'node:url'
const raiz = new URL('../', import.meta.url),
  api = new URL('apps/api/', raiz)
const cobertura = new URL('.artifacts/gocover/', raiz)
if (!process.env.TEST_DATABASE_URL)
  throw Error('Define TEST_DATABASE_URL: PostgreSQL real requerido')
if (!process.env.MUTACION) {
  rmSync(cobertura, { recursive: true, force: true })
}
mkdirSync(cobertura, { recursive: true })
const bin = fileURLToPath(
  new URL(`.artifacts/api-e2e${process.platform === 'win32' ? '.exe' : ''}`, raiz),
)
const env = {
  ...process.env,
  GOTOOLCHAIN: 'go1.26.11',
  API_BINARY: bin,
  GOCOVERDIR: fileURLToPath(cobertura),
  API_COVER_DIR: fileURLToPath(cobertura),
}
function go(args) {
  const r = spawnSync('go', args, { cwd: api, env, stdio: 'inherit' })
  if (r.error) throw r.error
  if (r.status !== 0) process.exit(r.status ?? 1)
}
go(['build', '-cover', '-covermode=atomic', '-coverpkg=./...', '-o', bin, './cmd/api'])
for (const tz of ['UTC', 'Pacific/Auckland']) {
  env.TZ = tz
  go(['test', './e2e', '-count=1', '-v'])
}
if (!process.env.MUTACION)
  go(['tool', 'covdata', 'textfmt', '-i', fileURLToPath(cobertura), '-o', 'coverage.out'])
