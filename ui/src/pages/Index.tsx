import { AuthLayout } from "@/layout/AuthLayout"
import { NavLayout } from "@/layout/NavLayout"
import { Outlet } from "@tanstack/react-router"

export const Index = () => {
  return (
    <AuthLayout>
      <NavLayout>
        <Outlet />
      </NavLayout>
    </AuthLayout>
  )
}
