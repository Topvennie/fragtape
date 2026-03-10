import { useInfiniteQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import { convertDemoFilterResult, DemoFilter, DemoFilterResult, DemoStatus } from "../types/demo"
import { STALE_TIME } from "../types/staletime"
import { apiGet, apiPost, NO_CONVERTER, NO_DATA } from "./query"

const ENDPOINT = "demo"

const PAGE_LIMIT = 10

export const useDemoGetFiltered = (filter?: DemoFilter) => {
  const { data, isLoading, fetchNextPage, isFetchingNextPage, hasNextPage, error, refetch, isFetching } = useInfiniteQuery({
    queryKey: ["demo", "filtered", JSON.stringify(filter)],
    queryFn: async ({ pageParam = 1 }) => {
      const queryParams = new URLSearchParams({
        page: pageParam.toString(),
        limit: PAGE_LIMIT.toString(),
      })

      if (filter?.source !== undefined) queryParams.append("source", filter?.source)
      if (filter?.result !== undefined) queryParams.append("result", filter?.result)
      if (filter?.PlayedAtStart !== undefined) queryParams.append("played_at_start", filter.PlayedAtStart.toISOString())
      if (filter?.PlayedAtEnd !== undefined) queryParams.append("played_at_end", filter.PlayedAtEnd.toISOString())
      if (filter?.HasHighlight !== undefined) queryParams.append("has_highlight", filter.HasHighlight.toString())

      const url = `${ENDPOINT}/filtered?${queryParams.toString()}`
      return (await apiGet(url, convertDemoFilterResult)).data
    },
    initialPageParam: 1,
    getNextPageParam: (lastPage, allPages) => {
      return lastPage.demos.length < PAGE_LIMIT ? undefined : allPages.length + 1
    },
    staleTime: STALE_TIME.MIN_5,
    refetchInterval: (query) => {
      const data = query.state.data
      if ((data?.pages.length ?? 0) === 0) return false

      const poll = data?.pages.flatMap(d => d.demos).some(d => ![DemoStatus.Finished, DemoStatus.Failed].includes(d.status))
      return poll ? STALE_TIME.SEC_5 : false
    },
  })

  const result: DemoFilterResult = { demos: data?.pages.flatMap(d => d.demos) ?? [], total: data?.pages?.[0].total ?? 0 }

  return {
    result,
    isLoading,
    fetchNextPage,
    isFetchingNextPage,
    hasNextPage,
    error,
    refetch,
    isFetching,
  }
}

export const useDemoUpload = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (file: File) => apiPost(`${ENDPOINT}/upload`, NO_DATA, NO_CONVERTER, [{ field: "file", file }]),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["demo"] }),
  })
}
