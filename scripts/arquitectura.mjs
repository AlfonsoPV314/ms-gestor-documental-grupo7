import { readdirSync, readFileSync } from 'node:fs'
function recorrer(dir) {
  return readdirSync(dir, { withFileTypes: true }).flatMap((e) =>
    e.isDirectory() ? recorrer(`${dir}/${e.name}`) : [`${dir}/${e.name}`],
  )
}
const fallas = []
for (const f of recorrer('apps/web/src').filter((f) => /\.tsx?$/.test(f))) {
  const s = readFileSync(f, 'utf8')
  if (/\b(?:localStorage|sessionStorage)\b/.test(s))
    fallas.push(`${f}: almacenamiento de sesión prohibido`)
  if (/\bfetch\s*\(/.test(s) && !f.endsWith('.api.ts'))
    fallas.push(`${f}: fetch fuera de cliente API`)
  if (
    f.includes('/componentes/') &&
    /from ['"]\.\.\/[^.]/.test(s) &&
    !f.endsWith('.api.ts')
  )
    fallas.push(`${f}: dependencia entre microcomponentes`)
}
for (const f of recorrer('apps/api/internal/api').filter((f) => f.endsWith('.go'))) {
  if (/internal\/(store|repository)/.test(readFileSync(f, 'utf8')))
    fallas.push(`${f}: handler accede a persistencia`)
}
if (fallas.length) throw Error(fallas.join('\n'))
console.log('Límites de arquitectura verificados')
