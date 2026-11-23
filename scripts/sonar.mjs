import { readFileSync } from 'node:fs'
import { spawnSync } from 'node:child_process'
import { fileURLToPath } from 'node:url'
const secreto = JSON.parse(
  readFileSync(new URL('../.secrets/sonar.json', import.meta.url), 'utf8'),
)
const raiz = fileURLToPath(new URL('../', import.meta.url))
const env = {
  ...process.env,
  SONAR_TOKEN: secreto.token,
  SONAR_HOST_URL: 'http://host.docker.internal:9000',
}
const r = spawnSync(
  'docker',
  [
    'run',
    '--rm',
    '--add-host',
    'host.docker.internal:host-gateway',
    '-e',
    'SONAR_TOKEN',
    '-e',
    'SONAR_HOST_URL',
    '-v',
    `${raiz}:/usr/src`,
    '-v',
    'documental-sonar-cache:/opt/sonar-scanner/.sonar',
    'sonarsource/sonar-scanner-cli:12.5.0.4351_8.5.0@sha256:8bb84373f00346f9a579adb7ccb2999f8cd050ab7167fa9ae08a16069de03c84',
    '-Dsonar.qualitygate.wait=true',
  ],
  { cwd: raiz, env, stdio: 'inherit' },
)
if (r.error) throw r.error
process.exit(r.status ?? 1)
