import "react-router";
import { createRequestHandler } from "@react-router/express";
import { createProxyServer } from "http-proxy-3";
import express, { type Response } from "express";

declare module "react-router" {
  interface AppLoadContext {
    nonce: string;
  }
}

const API_HOST = process.env.API_HOST;

export const app = express.Router();
const proxy = createProxyServer();

app.use(["/auth", "/api"], (req, res) => {
  req.url = req.originalUrl;
  proxy.web(req, res, { target: API_HOST });
});

app.use(
  createRequestHandler({
    build: () => import("virtual:react-router/server-build"),
    getLoadContext: (_req, res: Response) => ({
      nonce: res.locals.nonce as string,
    }),
  }),
);
