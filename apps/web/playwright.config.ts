import { defineConfig } from '@playwright/test'
export default defineConfig({
  testDir: 'tests/navegador',
  use: { baseURL: process.env.UI_URL ?? 'http://localhost:5170', headless: true },
  reporter: [['list'], ['html', { open: 'never' }]],
})
