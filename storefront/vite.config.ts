import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'
import tailwindcss from '@tailwindcss/vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue(), tailwindcss()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  server: {
    port: 5175, // Storefront on 5175
    host: true, // Listen on all addresses (0.0.0.0)
    strictPort: false,
    // Allow subdomain access like sophiaboutique.localhost:5175
    allowedHosts: [
      'localhost',
      '.localhost', // Allow all subdomains of localhost
      '127.0.0.1'
    ],
    proxy: {
      // Proxy API calls to backend on port 8080
      '/shops': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        secure: false,
      },
      '/auth': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        secure: false,
      },
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        secure: false,
      },
    },
  }
})
