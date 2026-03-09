import { AdminLayout } from "@/layout/AdminLayout"
import { Outlet } from "@tanstack/react-router"

export const Admin = () => {
  return (
    <AdminLayout>
      <Outlet />
    </AdminLayout>
  )
}
