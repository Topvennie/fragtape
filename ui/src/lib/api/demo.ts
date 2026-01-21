import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { convertDemos, DemoStatus } from "../types/demo"
import { apiGet, apiPost, NO_CONVERTER, NO_DATA } from "./query"
import { STALE_TIME } from "../types/staletime"

const ENDPOINT = "demo"

export const useDemoGetAll = () => {
  return useQuery({
    queryKey: ["demo"],
    queryFn: async () => (await apiGet(ENDPOINT, convertDemos)).data,
    staleTime: STALE_TIME.MIN_5,
    refetchInterval: (query) => {
      const demos = query.state.data
      if (!demos) return false

      const poll = demos.some(d => ![DemoStatus.Finished, DemoStatus.Failed].includes(d.status))
      return poll ? STALE_TIME.SEC_5 : false
    },
    retry: 0,
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
