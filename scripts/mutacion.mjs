import { readFileSync, writeFileSync, mkdirSync } from 'node:fs'
import { spawnSync } from 'node:child_process'
const fuente = new URL(
  '../apps/api/internal/repository/documentos.sql.go',
  import.meta.url,
)
const original = readFileSync(fuente, 'utf8')
const mutado = original.replace(
  'WHERE id=$1 AND propietario=$2',
  'WHERE id=$1 AND $2::uuid IS NOT NULL',
)
if (mutado === original)
  throw Error('No se encontró la protección; revisar prueba de mutación')
mkdirSync(new URL('../.artifacts/', import.meta.url), { recursive: true })
try {
  writeFileSync(fuente, mutado)
  const r = spawnSync(
    process.execPath,
    [new URL('./e2e.mjs', import.meta.url).pathname.replace(/^\/([A-Z]:)/, '$1')],
    {
      cwd: new URL('../', import.meta.url),
      encoding: 'utf8',
      env: { ...process.env, MUTACION: '1' },
    },
  )
  writeFileSync(
    new URL('../.artifacts/mutacion.log', import.meta.url),
    (r.stdout ?? '') + (r.stderr ?? ''),
  )
  if (
    r.status === 0 ||
    !((r.stdout ?? '') + (r.stderr ?? '')).includes('ajeno: estado 200, esperado 404')
  )
    throw Error(
      'La mutación no fue detectada por la aserción esperada; revisar .artifacts/mutacion.log',
    )
  console.log(
    'Mutación detectada: acceso ajeno produjo 200; evidencia en .artifacts/mutacion.log',
  )
  if (process.argv.includes('--evidencia-roja'))
    console.log((r.stdout ?? '') + (r.stderr ?? ''))
} finally {
  writeFileSync(fuente, original)
}
if (process.argv.includes('--evidencia-roja')) process.exit(1)
