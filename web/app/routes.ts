import { type RouteConfig, index, route } from "@react-router/dev/routes";

export default [
  index("routes/home.tsx"),
  route("settlements", "routes/settlements/page.tsx"),
] satisfies RouteConfig;
