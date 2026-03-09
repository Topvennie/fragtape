import { useAuth } from "@/lib/hooks/useAuth"
import { Forbidden } from "@/pages/Forbidden"
import { PropsWithChildren } from "react"

export const AdminLayout = ({ children }: PropsWithChildren) => {
  const { user } = useAuth()

  if (!user?.admin) return <Forbidden />

  return children
}
