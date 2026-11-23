import { spawnSync } from 'node:child_process'
import { mkdirSync, writeFileSync } from 'node:fs'
const raiz = new URL('../', import.meta.url)
const r = spawnSync('go', ['run', './cmd/api', 'spec'], {
  cwd: new URL('apps/api/', raiz),
  encoding: 'utf8',
  env: { ...process.env, GOTOOLCHAIN: 'go1.26.11' },
})
if (r.status !== 0) throw Error(r.stderr)
mkdirSync(new URL('contracts/', raiz), { recursive: true })
writeFileSync(new URL('contracts/openapi.yaml', raiz), r.stdout.replaceAll('\r\n', '\n'))
const t = spawnSync(
  process.execPath,
  [
    'apps/web/node_modules/openapi-typescript/bin/cli.js',
    'contracts/openapi.yaml',
    '-o',
    'apps/web/src/contrato.d.ts',
  ],
  { cwd: raiz, stdio: 'inherit' },
)
if (t.status !== 0) process.exit(t.status ?? 1)
