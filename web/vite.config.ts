import { purgeCss } from 'vite-plugin-tailwind-purgecss';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
  plugins: [sveltekit(), purgeCss()],
  define: {
    'process.env': {
      NODE_ENV: JSON.stringify(process.env.NODE_ENV)
    },
    'process.platform': JSON.stringify(process.platform),
    'process.cwd': JSON.stringify('/'),
    'process.browser': true,
    'process': {
      cwd: () => ('/')
    }
  },
  resolve: {
    alias: {
      process: 'process/browser'
    }
  },
  server: {
    fs: {
      allow: ['..']  // allows importing from the parent directory
    },
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, '')
      }
    },
    watch: {
      usePolling: true,
      interval: 100,
      ignored: ['**/node_modules/**', '**/dist/**', '**/.git/**', '**/.svelte-kit/**']
    }
  }
});
