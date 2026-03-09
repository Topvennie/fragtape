import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { convertSettingGlobal, SettingGlobalSchema } from "../types/setting_global"
import { STALE_TIME } from "../types/staletime"
import { apiGet, apiPost } from "./query"

const ENDPOINT = "setting/global"

export const useSettingGlobalGet = () => {
  return useQuery({
    queryKey: ["setting", "global"],
    queryFn: async () => (await apiGet(ENDPOINT, convertSettingGlobal)).data,
    staleTime: STALE_TIME.MIN_30,
  })
}

export const useSettingGlobalUpdate = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (setting: SettingGlobalSchema) => apiPost(ENDPOINT, setting),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["setting", "global"] })
  })
}
