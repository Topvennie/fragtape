import { useInfiniteQuery, useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { STALE_TIME } from "../types/staletime";
import { convertUser, convertUserFilterResult, convertUsers, User, UserFilter, UserFilterResult } from "../types/user";
import { apiDelete, apiGet, apiPost } from "./query";

const ENDPOINT_AUTH = "auth"
const ENDPOINT_USER = "user"

const PAGE_LIMIT = 5

export const useUser = () => {
  return useQuery({
    queryKey: ["user"],
    queryFn: async () => (await apiGet(`${ENDPOINT_USER}/me`, convertUser)).data,
    retry: 0,
    staleTime: STALE_TIME.MIN_30,
  })
}

export const useUserAdmin = () => {
  return useQuery({
    queryKey: ["user", "admin"],
    queryFn: async () => (await apiGet(`${ENDPOINT_USER}/admin`, convertUsers)).data,
    staleTime: STALE_TIME.MIN_30,
    retry: 0,
    throwOnError: true,
  })
}

export const useUserFiltered = (filter?: UserFilter) => {
  const { data, isLoading, fetchNextPage, isFetchingNextPage, hasNextPage, error, refetch, isFetching } = useInfiniteQuery({
    queryKey: ["user", "filtered", JSON.stringify(filter)],
    queryFn: async ({ pageParam = 1 }) => {
      const queryParams = new URLSearchParams({
        page: pageParam.toString(),
        limit: PAGE_LIMIT.toString(),
      })

      if (filter?.name !== undefined && filter.name !== "") {
        queryParams.append("name", filter.name)
      }
      if (filter?.admin !== undefined) {
        queryParams.append("admin", String(filter.admin))
      }
      if (filter?.real !== undefined) {
        queryParams.append("real", String(filter.real))
      }

      const url = `${ENDPOINT_USER}/filtered?${queryParams.toString()}`
      return (await apiGet(url, convertUserFilterResult)).data
    },
    initialPageParam: 1,
    getNextPageParam: (lastPage, allPages) => {
      return lastPage.users.length < PAGE_LIMIT ? undefined : allPages.length + 1
    },
    staleTime: STALE_TIME.MIN_5,
    throwOnError: true,
  })

  const result: UserFilterResult = { users: data?.pages.flatMap(p => p.users) ?? [], total: data?.pages?.[0].total ?? 0 }

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

export const useUserCreateAdmin = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (user: Pick<User, "id">) => apiPost(`${ENDPOINT_USER}/admin/${user.id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["user", "admin"] })
      queryClient.invalidateQueries({ queryKey: ["user", "filtered"] })
    },
  })
}

export const useUserDeleteAdmin = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (user: Pick<User, "id">) => apiDelete(`${ENDPOINT_USER}/admin/${user.id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["user", "admin"] })
      queryClient.invalidateQueries({ queryKey: ["user", "filtered"] })
    },
  })
}
