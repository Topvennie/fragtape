import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { apiDelete, apiGet, apiPost } from "./query"
import { convertSettingUser, SettingUserSteamSchema } from "../types/setting_user"
import { STALE_TIME } from "../types/staletime"

const ENDPOINT = "setting/user"

export const useSettingUserGet = () => {
  return useQuery({
    queryKey: ["setting", "user"],
    queryFn: async () => (await apiGet(ENDPOINT, convertSettingUser)).data,
    staleTime: STALE_TIME.MIN_30,
  })
}

export const useSettingUserSteamConnect = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (steam: SettingUserSteamSchema) => apiPost(`${ENDPOINT}/steam`, steam),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["setting", "user"] }),
  })
}

export const useSettingUserSteamDisconnect = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: () => apiDelete(`${ENDPOINT}/steam`),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["setting", "user"] }),
  })
}
