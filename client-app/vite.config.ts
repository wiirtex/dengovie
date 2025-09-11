/// <reference types="vitest/config" />
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
  ],
  test: {
    environment: 'jsdom',
    setupFiles: ['src/setupTests.ts'], // Подключаем setup файл
    reporters: [
      'default', // Vitest's default reporter so that terminal output is still visibles
    ],
    coverage: {
      provider: 'istanbul',
      reporter: [
        'text',
        ["lcov", { outputFile: "lcov.info", silent: false }],
      ],
      exclude: [
        '**/node_modules/**',
        '**/dist/**',
        '**/coverage/**',
        '**/tests/**',
        '**/*.test.{js,ts,jsx,tsx}',
        '**/*.spec.{js,ts,jsx,tsx}',
        '**/__tests__/**',
        '**/__mocks__/**',
        '**/types/**',
        '**/*.{js,ts}', // игнорировать файлы index
        'src/App.vue',
      ]
    }
  },
  server: {
    host: true,
    port: 5173
  },
  preview: {
    host: true,
    port: 5173
  }
})
