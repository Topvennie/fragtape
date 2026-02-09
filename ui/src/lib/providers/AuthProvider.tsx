import { notifications } from "@mantine/notifications";
import { PropsWithChildren, useCallback, useEffect, useMemo, useState } from "react";
import { isResponseNot200Error } from "../api/query";
import { useUser, useUserLogin, useUserLogout } from "../api/user";
import { AuthContext } from "../contexts/authContext";

export const AuthProvider = ({ children }: PropsWithChildren) => {
  const [forbidden, setForbidden] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const { data, isLoading, isFetching, error: userError } = useUser();
  const { mutate: logoutMutation } = useUserLogout();

  const user = data ?? null;

  useEffect(() => {
    if (userError) {
      if (!isResponseNot200Error(userError)) {
        setError(userError);
        return;
      }

      if (userError.response.status === 403) setForbidden(true);
      else if (userError.response.status !== 401) setError(userError);
      return;
    }

    setForbidden(false);
    setError(null);
  }, [userError]);

  const logout = useCallback(() => {
    logoutMutation(undefined, {
      onSuccess: () => notifications.show({ message: "Logged out" }),
      onError: (err) => console.log(`Logout failed ${err}`),
    });
  }, [logoutMutation]);

  const value = useMemo(
    () => ({
      user,
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
