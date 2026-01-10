import { AuthLayout } from "@/layout/AuthLayout"
import { Outlet } from "@tanstack/react-router"

export const Index = () => {
  return (
    <div className="w-screen h-screen bg-(--mantine-color-background-8)">
      <AuthLayout>
        <Outlet />
      </AuthLayout>
    </div>
  )
}
