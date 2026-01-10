import { useState, useEffect, useCallback, useMemo, PropsWithChildren } from "react";
import { notifications } from "@mantine/notifications";
import { isResponseNot200Error } from "../api/query";
import { useUser, useUserLogin, useUserLogout } from "../api/user";
import { AuthContext } from "../contexts/authContext";
import { User } from "../types/user";

export const AuthProvider = ({ children }: PropsWithChildren) => {
  const [user, setUser] = useState<User | null>(null);
  const [forbidden, setForbidden] = useState(false);
  const [error, setError] = useState<Error | null>(null)

  const { data, isLoading, error: userError } = useUser();
  const { mutate: logoutMutation } = useUserLogout();

  useEffect(() => {
    if (data) {
      setUser(data);
      setForbidden(false);
    }
  }, [data]);

  useEffect(() => {
    if (userError) {
      if (!isResponseNot200Error(userError)) {
        setError(userError)
        return
      }

      if (userError.response.status === 403) setForbidden(true)
      else if (userError.response.status !== 401) setError(userError)

      return
    }

    setForbidden(false);
    setError(null)
  }, [userError]);

  const logout = useCallback(() => {
    logoutMutation(undefined, {
      onSuccess: () => notifications.show({ message: "Logged out" }),
      onError: (err) => { console.log(`Logout failed ${err}`) },
      onSettled: () => setUser(null),
    });
  }, [logoutMutation]);

  const value = useMemo(() => ({ user, isLoading, forbidden, error, login: useUserLogin, logout }), [user, isLoading, forbidden, error, logout]);

  return <AuthContext value={value}>{children}</AuthContext>;
}
