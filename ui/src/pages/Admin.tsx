import { AdminConfiguration } from "@/components/admin/AdminConfiguration"
import { AdminTeam } from "@/components/admin/AdminTeam"
import { Page, PageTitle } from "@/components/atoms/Page"
import { FragtapeIcon } from "@/components/icons/FragtapeIcon"
import { useSettingGlobalGet } from "@/lib/api/setting_global"
import { useUserAdmin } from "@/lib/api/user"
import { useAuth } from "@/lib/hooks/useAuth"
import { Center } from "@mantine/core"
import { Forbidden } from "./Forbidden"

export const Admin = () => {
  const { user } = useAuth()

  const { isLoading: isLoadingAdmins } = useUserAdmin()
  const { isLoading: isLoadingSetting } = useSettingGlobalGet()

  if (!user?.admin) return <Forbidden />

  if (isLoadingAdmins || isLoadingSetting) {
    return (
      <Center>
        <FragtapeIcon animated className="text-(--mantine-color-primary-6) size-12" />
      </Center>
    )
  }

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
