import { createRootRouteWithContext, createRoute, createRouter } from "@tanstack/react-router";
import { App } from "./App";
import { Index } from "./pages/Index";
import { Error404 } from "./pages/404";
import { Error } from "./pages/Error";
import { Home } from "./pages/Home";
import { Admin } from "./pages/Admin";

const root = createRootRouteWithContext()({
  component: App,
})

const index = createRoute({
  getParentRoute: () => root,
  id: "public-layout",
  component: Index,
})

const home = createRoute({
  getParentRoute: () => index,
  path: "/",
  component: Home,
})

const admin = createRoute({
  getParentRoute: () => index,
  path: "/admin",
  component: Admin,
})

const routeTree = root.addChildren([
  index.addChildren([
    home,
    admin,
  ]),
])

export const router = createRouter({
  routeTree,
  defaultPreload: "render",
  defaultPreloadStaleTime: 0, // Data is immediatly marked as stale and will refetch when the user navigates to the page
  scrollRestoration: true,
  defaultErrorComponent: Error,
  defaultNotFoundComponent: Error404,
})

declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}
