import { sveltekit } from '@sveltejs/kit/vite';
import type { UserConfig } from 'vite';

export default ({
  plugins: [sveltekit()],
  server: {
    port: 8000,
    strictPort: true,
    proxy: {
      '/auth': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
}) satisfies UserConfig;
