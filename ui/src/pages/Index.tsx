import { UserFirstTime } from "@/components/user/UserFirstTime"
import { AuthLayout } from "@/layout/AuthLayout"
import { DataLayout } from "@/layout/DataLayout"
import { NavLayout } from "@/layout/NavLayout"
import { Outlet } from "@tanstack/react-router"

export const Index = () => {
  return (
    <AuthLayout>
      <DataLayout>
        <NavLayout>
          <UserFirstTime />

          <Outlet />
        </NavLayout>
      </DataLayout>
    </AuthLayout>
  )
}
