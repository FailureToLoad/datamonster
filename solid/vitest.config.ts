// vitest.config.ts
import { defineConfig, mergeConfig } from "vitest/config";
import viteConfig from "./vite.config";

export default mergeConfig(
  viteConfig,
  defineConfig({
    resolve: {
      alias: {
        "~test": "/test",
      },
      dedupe: ["solid-js", "solid-js/web"],
    },
    test: {
      watch: false,
      globals: false,
      environment: "happy-dom",
      setupFiles: ["./test/setup.ts"],
      server: {
        deps: {
          inline: [
            "@solidjs/testing-library",
            "@solidjs/router",
            "lucide-solid",
            "@tanstack/solid-form",
            "@tanstack/solid-store",
          ],
        },
      },
    },
  })
);
