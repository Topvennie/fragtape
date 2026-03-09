import { AdminConfiguration } from "@/components/admin/AdminConfiguration"
import { AdminTeam } from "@/components/admin/AdminTeam"
import { Loading } from "@/components/atoms/Loading"
import { Page, PageTitle } from "@/components/atoms/Page"
import { useUserAdmin } from "@/lib/api/user"

export const SettingGlobal = () => {
  // Preload some data
  const { isLoading: isLoading } = useUserAdmin()

  if (isLoading) return <Loading />

  return (
    <Page>
      <PageTitle
        title="Admin Settings"
        description="Manage global platform configuration, permissions and admin access."
      />

      <AdminTeam />

      <AdminConfiguration />
    </Page>
  )
}
