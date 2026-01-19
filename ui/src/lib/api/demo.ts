import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { apiGet, apiPost, NO_CONVERTER, NO_DATA } from "./query"
import { convertDemos } from "../types/demo"
import { STALE_TIME } from "../types/staletime"

const ENDPOINT = "demo"

export const useDemoGetAll = () => {
  return useQuery({
    queryKey: ["demo"],
    queryFn: async () => (await apiGet(ENDPOINT, convertDemos)).data,
    retry: 0,
    staleTime: STALE_TIME.MIN_5,
    throwOnError: true,
  })
}

export const useDemoUpload = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (file: File) => apiPost(`${ENDPOINT}/upload`, NO_DATA, NO_CONVERTER, [{ field: "file", file }]),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["demo"] }),
  })
}
