import { AuthLayout } from "@/layout/AuthLayout"
import { Outlet } from "@tanstack/react-router"

export const Index = () => {
  return (
    <AuthLayout>
      <Outlet />
    </AuthLayout>
  )
}
