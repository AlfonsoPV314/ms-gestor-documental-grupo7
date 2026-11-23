import { expect, test } from '@playwright/test'
test('carga, verifica y limita al lector', async ({ page }) => {
  await page.goto('/')
  await page.getByLabel('Identidad de prueba').selectOption('operador')
  await expect(page.getByRole('heading', { name: 'Biblioteca documental' })).toBeVisible()
  const nombre = `prueba-${Date.now()}.txt`
  await page.getByLabel('Archivo para cargar').setInputFiles({
    name: nombre,
    mimeType: 'text/plain',
    buffer: Buffer.from('contenido original desde navegador'),
  })
  await page.getByRole('button', { name: 'Guardar y firmar' }).click()
  await expect(
    page.getByRole('status').filter({ hasText: 'se guardó y firmó' }),
  ).toBeVisible()
  await page.getByRole('button').filter({ hasText: nombre }).click()
  await page.getByRole('button', { name: 'Verificar firma' }).click()
  await expect(page.getByRole('status').filter({ hasText: 'Firma válida' })).toBeVisible()
  const download = page.waitForEvent('download')
  await page.getByRole('button', { name: 'Descargar original' }).click()
  expect((await download).suggestedFilename()).toBe(nombre)
  await page.getByLabel('Identidad de prueba').selectOption('lector')
  await expect(page.getByRole('button', { name: 'Guardar y firmar' })).toBeDisabled()
  await page.getByLabel('Identidad de prueba').selectOption('ajeno')
  await expect(page.getByRole('button').filter({ hasText: nombre })).toHaveCount(0)
  await page.screenshot({ path: 'test-results/boveda.png', fullPage: true })
})
