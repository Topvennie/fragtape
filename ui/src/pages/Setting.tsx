import { Title } from "@/components/atoms/Title"
import { FragtapeIcon } from "@/components/icons/FragtapeIcon"
import { SettingSteam } from "@/components/setting/SettingSteam"
import { useSettingUserGet } from "@/lib/api/setting_user"
import { Center, Stack } from "@mantine/core"

export const Setting = () => {
  const { isLoading } = useSettingUserGet()

  if (isLoading) {
    return (
      <Center>
        <FragtapeIcon animated className="text-(--mantine-color-primary-6) size-12" />
      </Center>
    )
  }

  return (
    <Stack gap="xl">
      <Title order={2} className="font-bold">User Settings</Title>
      <p className="text-secondary mb-4">Manage your external integrarions and account connections.</p>

      <SettingSteam />
    </Stack>
  )
}
