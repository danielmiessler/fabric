import { purgeCss } from 'vite-plugin-tailwind-purgecss';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import { mdsvex } from 'mdsvex';

export default defineConfig({
	plugins: [sveltekit(), mdsvex(), purgeCss()],
	server: {
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
			ignored: ['**/node_modules/**', '**/dist/**', '**/.git/**']
		}
	},
	optimizeDeps: {}
});