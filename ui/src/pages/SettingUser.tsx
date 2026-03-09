import { Loading } from "@/components/atoms/Loading"
import { Page, PageTitle } from "@/components/atoms/Page"
import { SettingSteam } from "@/components/setting/SettingSteam"
import { useSettingUserGet } from "@/lib/api/setting_user"

export const SettingUser = () => {
  const { isLoading } = useSettingUserGet()

  if (isLoading) return <Loading />

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
