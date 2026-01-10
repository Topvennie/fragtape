import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { convertUser } from "../types/user";
import { apiGet, apiPost } from "./query";
import { STALE_TIME } from "../types/staletime";

const ENDPOINT_AUTH = "auth"
const ENDPOINT_USER = "user"

export const useUser = () => {
  return useQuery({
    queryKey: ["user"],
    queryFn: async () => (await apiGet(`${ENDPOINT_USER}/me`, convertUser)).data,
    retry: 0,
    staleTime: STALE_TIME.MIN_30,
  })
}

export const useUserLogin = () => {
  window.location.href = `/api/${ENDPOINT_AUTH}/login/steam`
}

export const useUserLogout = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async () => (await apiPost(`${ENDPOINT_AUTH}/logout`)).data,
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["user"] })
  })
}
