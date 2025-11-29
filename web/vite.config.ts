import type { UserConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'


export default ({
  plugins: [svelte()],
}) satisfies UserConfig
