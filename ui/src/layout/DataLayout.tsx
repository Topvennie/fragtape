import { useSettingGlobalGet } from "@/lib/api/setting_global"
import { useSettingUserGet } from "@/lib/api/setting_user"
import { PropsWithChildren } from "react"

export const DataLayout = ({ children }: PropsWithChildren) => {
  const { isLoading: isLoadingSetting } = useSettingGlobalGet()
  const { isLoading: isLoadingUserSetting } = useSettingUserGet()

  if (isLoadingSetting || isLoadingUserSetting) return null

  return children
}
