import { copyFileSync } from 'node:fs'
copyFileSync(
  new URL('../apps/web/src/contrato.d.ts', import.meta.url),
  new URL('../apps/web/dist/contrato.d.ts', import.meta.url),
)
