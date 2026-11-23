import { mkdirSync, readFileSync, writeFileSync } from 'node:fs'
import { spawnSync } from 'node:child_process'
import { fileURLToPath } from 'node:url'
const raiz = fileURLToPath(new URL('../', import.meta.url))
mkdirSync(new URL('../.artifacts/', import.meta.url), { recursive: true })
const trivy =
  'aquasec/trivy:0.76.1@sha256:790fd66d1697395d8693d1acaa1ea702406a20f109b6fc7bac8db623eca6c06c'
function docker(args, archivo) {
  const r = spawnSync('docker', args, { cwd: raiz, encoding: 'utf8' })
  if (archivo)
    writeFileSync(
      new URL(`../.artifacts/${archivo}`, import.meta.url),
      (r.stdout ?? '') + (r.stderr ?? ''),
    )
  if (r.error) throw r.error
  if (r.status !== 0)
    throw Error(`Falló ${archivo ?? args.join(' ')}; revisar .artifacts`)
}
const comun = [
  'run',
  '--rm',
  '-v',
  `${raiz}:/src`,
  '-v',
  'documental-trivy-cache:/root/.cache/trivy',
]
docker(
  [
    ...comun,
    trivy,
    'fs',
    '--quiet',
    '--scanners',
    'vuln,misconfig,secret',
    '--severity',
    'HIGH,CRITICAL',
    '--ignore-unfixed',
    '--exit-code',
    '1',
    '--skip-dirs',
    '/src/.secrets,/src/.artifacts,/src/.tools,/src/node_modules,/src/.git',
    '--format',
    'json',
    '--output',
    '/src/.artifacts/trivy-fs.json',
    '/src',
  ],
  'trivy-fs.log',
)
for (const [nombre, imagen] of [
  ['api', 'ms-gestor-documental-backend:local'],
  ['web', 'ms-gestor-documental-frontend:local'],
  ['migrate', 'ms-gestor-documental-migrate:local'],
]) {
  docker(
    [
      ...comun,
      '-v',
      '/var/run/docker.sock:/var/run/docker.sock',
      trivy,
      'image',
      '--quiet',
      '--scanners',
      'vuln,secret',
      '--severity',
      'HIGH,CRITICAL',
      '--ignore-unfixed',
      '--exit-code',
      '1',
      '--format',
      'json',
      '--output',
      `/src/.artifacts/trivy-${nombre}.json`,
      imagen,
    ],
    `trivy-${nombre}.log`,
  )
  docker(
    [
      ...comun,
      '-v',
      '/var/run/docker.sock:/var/run/docker.sock',
      trivy,
      'image',
      '--quiet',
      '--format',
      'cyclonedx',
      '--output',
      `/src/.artifacts/sbom-${nombre}.cdx.json`,
      imagen,
    ],
    `sbom-${nombre}.log`,
  )
  const reporte = JSON.parse(
    readFileSync(new URL(`../.artifacts/trivy-${nombre}.json`, import.meta.url), 'utf8'),
  )
  console.log(
    `${nombre}: ${(reporte.Results ?? []).flatMap((r) => r.Vulnerabilities ?? []).length} HIGH/CRITICAL corregibles; SBOM generado`,
  )
}
console.log('Trivy repositorio e imágenes: sin hallazgos bloqueantes')
