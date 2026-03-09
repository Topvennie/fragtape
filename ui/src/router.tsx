import { createRootRouteWithContext, createRoute, createRouter } from "@tanstack/react-router";
import { App } from "./App";
import { Index } from "./pages/Index";
import { Error404 } from "./pages/404";
import { Error } from "./pages/Error";
import { Overview } from "./pages/Overview";
import { SettingUser } from "./pages/SettingUser";
import { Admin } from "./pages/admin/Admin";
import { SettingGlobal } from "./pages/admin/SettingGlobal";

const root = createRootRouteWithContext()({
  component: App,
})

// Public (still behind auth)

const index = createRoute({
  getParentRoute: () => root,
  id: "public-layout",
  component: Index,
})

const overview = createRoute({
  getParentRoute: () => index,
  path: "/",
  component: Overview,
})

const settingUser = createRoute({
  getParentRoute: () => index,
  path: "/setting",
  component: SettingUser,
})

// Admin

const admin = createRoute({
  getParentRoute: () => index,
  path: "/admin",
  component: Admin,
})

const settingGlobal = createRoute({
  getParentRoute: () => admin,
  path: "/setting",
  component: SettingGlobal,
})


const routeTree = root.addChildren([
  index.addChildren([
    overview,
    settingUser,
    admin.addChildren([
      settingGlobal,
    ]),
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
