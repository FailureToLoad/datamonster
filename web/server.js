import { randomBytes } from "node:crypto";
import compression from "compression";
import express from "express";
import morgan from "morgan";
import { createProxyServer } from "http-proxy-3";

const BUILD_PATH = "./build/server/index.js";
const API_URL = process.env.API_URL;
const AUTH_DOMAIN = process.env.AUTH_DOMAIN;
const DEVELOPMENT = process.env.NODE_ENV === "development";
const PORT = Number.parseInt(process.env.PORT);

const app = express();
const proxy = createProxyServer();

app.use("/auth", (req, res) => {
  proxy.web(req, res, { target: `${API_URL}/auth` });
});

app.use(compression());
app.disable("x-powered-by");

app.use((_req, res, next) => {
  res.locals.nonce = randomBytes(16).toString("base64");
  next();
});

app.use((_req, res, next) => {
  const nonce = res.locals.nonce;
  const csp = [
    "default-src 'self'",
    `script-src 'self' 'nonce-${nonce}' 'strict-dynamic'`,
    "style-src 'self' 'unsafe-inline' https://fonts.googleapis.com",
    "img-src 'self' data:",
    "font-src 'self' https://fonts.gstatic.com",
    `connect-src 'self' ${AUTH_DOMAIN}`,
    "frame-ancestors 'self'",
    `form-action 'self' ${AUTH_DOMAIN}`,
    "base-uri 'self'",
  ].join("; ");
  res.setHeader("Content-Security-Policy", csp);
  next();
});

if (DEVELOPMENT) {
  console.log("Starting development server");
  const viteDevServer = await import("vite").then((vite) =>
    vite.createServer({
      server: { middlewareMode: true },
    })
  );
  app.use(viteDevServer.middlewares);
  app.use(async (req, res, next) => {
    try {
      const source = await viteDevServer.ssrLoadModule("./server/app.ts");
      return await source.app(req, res, next);
    } catch (error) {
      if (typeof error === "object" && error instanceof Error) {
        viteDevServer.ssrFixStacktrace(error);
      }
      next(error);
    }
  });
} else {
  console.log("Starting production server");
  app.use(
    "/assets",
    express.static("build/client/assets", { immutable: true, maxAge: "1y" })
  );
  app.use(morgan("tiny"));
  app.use(express.static("build/client", { maxAge: "1h" }));
  app.use(await import(BUILD_PATH).then((mod) => mod.app));
}

app.listen(PORT, () => {
  console.log(`Server is running on http://localhost:${PORT}`);
});
