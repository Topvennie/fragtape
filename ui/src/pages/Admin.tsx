import { useAuth } from "@/lib/hooks/useAuth"
import { Forbidden } from "./Forbidden"
import { Center, Stack } from "@mantine/core"
import { Title } from "@/components/atoms/Title"
import { AdminTeam } from "@/components/admin/AdminTeam"
import { AdminConfiguration } from "@/components/admin/AdminConfiguration"
import { FragtapeIcon } from "@/components/icons/FragtapeIcon"
import { useSettingGlobalGet } from "@/lib/api/setting_global"
import { useUserAdmin } from "@/lib/api/user"

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
    <Stack gap="xl">
      <Title order={2} className="font-bold">Admin Settings</Title>
      <p className="text-secondary mb-4">Manage global platform configuration, permissions and admin access.</p>

      <AdminTeam />

      <AdminConfiguration />
    </Stack>
  )
}
