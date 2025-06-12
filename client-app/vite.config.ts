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
  },
  server: {
    host: true,
    port: 5173,
    allowedHosts: ['dengovie.ingress'],
    cors: {
      origin: ['http://dengovie.ingress', 'http://localhost:5173'],
      methods: ['POST', 'OPTIONS', 'GET', 'PUT', 'DELETE'],
      allowedHeaders: ['Content-Type', 'Cookie', 'Content-Length', 'Accept-Encoding', 'X-CSRF-Token', 'Authorization', 'Pragma', 'accept', 'origin', 'Cache-Control', 'X-Requested-With']
    },
  },
  preview: {
    host: true,
    port: 5173
  }
})
