import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/api': 'http://127.0.0.1:9091',
      '/products': 'http://127.0.0.1:9091',
      '/customers': 'http://127.0.0.1:9091',
      '/vendors': 'http://127.0.0.1:9091',
      '/orders': 'http://127.0.0.1:9091',
      '/quotes': 'http://127.0.0.1:9091',
      '/invoices': 'http://127.0.0.1:9091',
      '/locations': 'http://127.0.0.1:9091',
      '/health': 'http://127.0.0.1:9091',
      '/activities': 'http://127.0.0.1:9091',
      '/contacts': 'http://127.0.0.1:9091',
      '/documents': 'http://127.0.0.1:9091',
      '/gl': 'http://127.0.0.1:9091',
      '/parsing': 'http://127.0.0.1:9091',
      '/price_levels': 'http://127.0.0.1:9091',
      '/pricing': 'http://127.0.0.1:9091',
      '/purchase-orders': 'http://127.0.0.1:9091'
    }
  }
})
