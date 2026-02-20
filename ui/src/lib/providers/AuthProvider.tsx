import { notifications } from "@mantine/notifications";
import { PropsWithChildren, useCallback, useEffect, useMemo, useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { isResponseNot200Error } from "../api/query";
import { useUser, useUserLogin, useUserLogout } from "../api/user";
import { AuthContext } from "../contexts/authContext";

export const AuthProvider = ({ children }: PropsWithChildren) => {
  const queryClient = useQueryClient();

  const [forbidden, setForbidden] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const { data: user, isLoading, isFetching, error: userError } = useUser();
  const { mutate: logoutMutation } = useUserLogout();

  useEffect(() => {
    if (!userError) {
      setForbidden(false);
      setError(null);
      return;
    }

    if (!isResponseNot200Error(userError)) {
      setError(userError);
      return;
    }

    if (userError.response.status === 403) setForbidden(true);
    else if (userError.response.status !== 401) setError(userError);
  }, [userError]);

  const logout = useCallback(() => {
    logoutMutation(undefined, {
      onSuccess: async () => {
        queryClient.setQueryData(["user"], null)
        await queryClient.cancelQueries({ queryKey: ["user"] });

        notifications.show({ message: "Logged out" });
      },
      onError: (err) => console.log(`Logout failed ${err}`),
    });
  }, [logoutMutation, queryClient]);

  const value = useMemo(
    () => ({
      user: user ?? null,
      // important: treat "loading/fetching" as "unknown auth", not "logged out"
      isLoading: isLoading || isFetching,
      forbidden,
      error,
      login: useUserLogin,
      logout,
    }),
    [user, isLoading, isFetching, forbidden, error, logout]
  );

  return <AuthContext value={value}>{children}</AuthContext>;
};

