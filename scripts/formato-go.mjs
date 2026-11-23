import { spawnSync } from 'node:child_process'
const cwd = new URL('../apps/api/', import.meta.url)
const formato = spawnSync('gofmt', ['-l', '.'], { cwd, encoding: 'utf8' })
if (formato.error) throw formato.error
if (formato.status !== 0 || formato.stdout.trim()) {
  console.error(formato.stdout, formato.stderr)
  process.exit(1)
}
const vet = spawnSync('go', ['vet', './...'], {
  cwd,
  stdio: 'inherit',
  env: { ...process.env, GOTOOLCHAIN: 'go1.26.11' },
})
if (vet.error) throw vet.error
process.exit(vet.status ?? 1)
