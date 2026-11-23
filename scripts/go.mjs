import { spawnSync } from 'node:child_process'
import { fileURLToPath } from 'node:url'
if (process.platform === 'win32' && process.argv[2] === 'test') {
  const raiz = fileURLToPath(new URL('../', import.meta.url))
  const build = spawnSync(
    'docker',
    [
      'build',
      '--quiet',
      '-f',
      'apps/api/Dockerfile.test',
      '-t',
      'documental-tests:local',
      '.',
    ],
    { cwd: raiz, stdio: 'inherit' },
  )
  if (build.status !== 0) process.exit(build.status ?? 1)
  console.log('Pruebas Go en Linux mediante Docker')
  const prueba = spawnSync(
    'docker',
    [
      'run',
      '--rm',
      '--user',
      '0',
      '-v',
      `${raiz}:/src`,
      '-v',
      'documental-go-mod:/go/pkg/mod',
      '-v',
      'documental-go-cache:/root/.cache/go-build',
      '-w',
      '/src/apps/api',
      'documental-tests:local',
      'go',
      ...process.argv.slice(2),
    ],
    { cwd: raiz, stdio: 'inherit' },
  )
  if (prueba.error) throw prueba.error
  process.exit(prueba.status ?? 1)
}
const r = spawnSync('go', process.argv.slice(2), {
  cwd: new URL('../apps/api/', import.meta.url),
  stdio: 'inherit',
  env: { ...process.env, GOTOOLCHAIN: 'go1.26.11' },
})
if (r.error) throw r.error
process.exit(r.status ?? 1)
