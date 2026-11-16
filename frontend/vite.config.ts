import { fileURLToPath, URL } from 'node:url'
import vue from '@vitejs/plugin-vue'
import { defineConfig, loadEnv } from 'vite'

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')
  const base = (env.VITE_BASE_DOMAIN || 'localhost').toLowerCase()
  const allowedHosts = ['.localhost', 'localhost']
  if (base !== 'localhost') allowedHosts.push(base, `.${base}`)
  return {
    plugins: [vue()],
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url)),
      },
    },
    server: {
      host: true,
      port: 5173,
      allowedHosts,
      proxy: {
        '/api': {
          target: 'http://localhost:8080',
          changeOrigin: false,
        },
      },
    },
    preview: {
      host: true,
      port: 4173,
      allowedHosts,
    },
  }
})
