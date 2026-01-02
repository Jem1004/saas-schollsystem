import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  server: {
    port: 5173,
    strictPort: true,
    proxy: {
      // WebSocket proxy - must be before /api to match first
      '/api/v1/ws': {
        target: 'http://localhost:3000',
        changeOrigin: true,
        ws: true,
      },
      // REST API proxy
      '/api': {
        target: 'http://localhost:3000',
        changeOrigin: true,
        ws: true, // Enable WebSocket proxy for other /api paths
      },
    },
  },
})
