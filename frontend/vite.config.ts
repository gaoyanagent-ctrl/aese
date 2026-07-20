import { defineConfig } from 'vitest/config';
import react from '@vitejs/plugin-react';

const iaosProxy = {
  '/api': {
    target: 'http://127.0.0.1:8082',
    changeOrigin: false,
  },
};

export default defineConfig({
  plugins: [react()],
  server: { port: 4173, proxy: iaosProxy },
  preview: { port: 4173, proxy: iaosProxy },
  test: {
    include: ['src/**/*.{test,spec}.{ts,tsx}'],
    environment: 'jsdom',
    setupFiles: './src/test/setup.ts',
    css: true,
    coverage: { reporter: ['text', 'html'] },
  },
});
