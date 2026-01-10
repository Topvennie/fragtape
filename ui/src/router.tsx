import { createRootRouteWithContext, createRoute, createRouter } from "@tanstack/react-router";
import { App } from "./App";
import { Index } from "./pages/Index";
import { Error404 } from "./pages/404";
import { Error } from "./pages/Error";

const root = createRootRouteWithContext()({
  component: App,
})

const index = createRoute({
  getParentRoute: () => root,
  id: "public-layout",
  component: Index,
})

const routeTree = root.addChildren([
  index.addChildren([
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
