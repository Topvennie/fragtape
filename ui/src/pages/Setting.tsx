import { Page, PageTitle } from "@/components/atoms/Page"
import { FragtapeIcon } from "@/components/icons/FragtapeIcon"
import { SettingSteam } from "@/components/setting/SettingSteam"
import { useSettingUserGet } from "@/lib/api/setting_user"
import { Center } from "@mantine/core"

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
    <Page>
      <PageTitle
        title="User Settings"
        description="Manage your account connections."
      />

      <SettingSteam />
    </Page>
  )
}
