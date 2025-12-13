import "react-router";
import { createRequestHandler } from "@react-router/express";
import { createProxyMiddleware } from "http-proxy-middleware";
import express from "express";

const API_HOST = process.env.API_HOST || "http://localhost:8080";

export const app = express();

app.use(
  ["/auth", "/api"],
  createProxyMiddleware({
    target: API_HOST,
    changeOrigin: true,
  }),
);

app.use(
  createRequestHandler({
    build: () => import("virtual:react-router/server-build"),
  }),
);
