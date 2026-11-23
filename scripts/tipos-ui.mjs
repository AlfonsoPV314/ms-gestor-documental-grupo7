import { copyFileSync, readFileSync, writeFileSync } from 'node:fs'
copyFileSync(
  new URL('../apps/web/src/contrato.d.ts', import.meta.url),
  new URL('../apps/web/dist/contrato.d.ts', import.meta.url),
)
const indice = new URL('../apps/web/dist/index.d.ts', import.meta.url)
writeFileSync(
  indice,
  readFileSync(indice, 'utf8').replace(/^import ['"].*\.css['"];?\r?\n/gm, ''),
)
