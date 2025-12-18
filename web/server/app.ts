import "react-router";
import { createRequestHandler } from "@react-router/express";
import { createProxyServer } from "http-proxy-3";
import express from "express";

const API_HOST = process.env.API_HOST;

export const app = express();
const proxy = createProxyServer();

app.use(["/auth", "/api"], (req, res) => {
  proxy.web(req, res, { target: API_HOST });
});

app.use(
  createRequestHandler({
    build: () => import("virtual:react-router/server-build"),
  }),
);
