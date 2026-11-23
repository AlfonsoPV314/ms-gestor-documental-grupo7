import { expect, test } from '@playwright/test'

test('Swagger UI carga el contrato y prueba la API con la sesión del mismo origen', async ({
  page,
}) => {
  await page.goto('/')
  await page.getByLabel('Identidad de prueba').selectOption('operador')
  await expect(page.getByRole('heading', { name: 'Biblioteca documental' })).toBeVisible()
  const externos: string[] = []
  page.on('request', (request) => {
    if (new URL(request.url()).origin !== new URL(page.url()).origin)
      externos.push(request.url())
  })
  await page.goto('/api/docs')
  await expect(page).toHaveTitle('Swagger UI · Gestor Documental')
  await expect(page.locator('.swagger-ui .opblock')).toHaveCount(7)
  const listado = page.locator('#operations-Documentos-listar-documentos')
  await listado.locator('.opblock-summary').click()
  await listado.getByRole('button', { name: 'Try it out' }).click()
  await listado.getByRole('button', { name: 'Execute', exact: true }).click()
  await expect(
    listado.locator('.live-responses-table .response-col_status:not(.col_header)'),
  ).toContainText('200')
  expect(externos).toEqual([])
})
