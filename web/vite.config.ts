/// <reference types="vitest" />

import path from "path";
import react from "@vitejs/plugin-react-swc";
import { defineConfig as defineViteConfig, mergeConfig } from "vite";
import { defineConfig as defineVitestConfig } from "vitest/config";

const vitestConfig = defineVitestConfig({
  test: {
    globals: true,
    environment: "jsdom",
    setupFiles: "./tests/setup.ts",
  },
});

const viteConfig = defineViteConfig({
  plugins: [react()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  server: {
    host: true, // needed for the DC port mapping to work
    strictPort: true,
    port: 8090,
  },
});

export default mergeConfig(viteConfig, vitestConfig);
