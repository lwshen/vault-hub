import path from 'path';
import tailwindcss from '@tailwindcss/vite';
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:3000',
        changeOrigin: true,
      },
    },
  },
  build: {
    rollupOptions: {
      output: {
        manualChunks(id) {
          if (/\/(@radix-ui|framer-motion|lucide-react)/.test(id)) {
            return 'ui-libs';
          }
          if (/\/react(?:-dom)?/.test(id)) {
            return 'vendor';
          }
          if (/@lwshen\/vault-hub-ts-fetch-client/.test(id)) {
            return 'api';
          }
          return null;
        },
      },
    },
  },
});
