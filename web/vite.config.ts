import { reactRouter } from "@react-router/dev/vite";
import tailwindcss from "@tailwindcss/vite";
import { defineConfig } from "vite";
import tsconfigPaths from "vite-tsconfig-paths";

export default defineConfig(({ isSsrBuild }) => ({
  build: {
    rollupOptions: isSsrBuild
      ? { input: "./server/app.ts" }
      : undefined,
  },
  ssr: {
    noExternal: ["zod"],
  },
  plugins: [tailwindcss(), reactRouter(), tsconfigPaths()],
  server: {
    port: 8000,
    proxy: {
      "/auth": "http://localhost:8080",
      "/api": "http://localhost:8080",
    },
  },
}));
