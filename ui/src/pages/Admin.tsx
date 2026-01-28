import { useAuth } from "@/lib/hooks/useAuth"
import { Forbidden } from "./Forbidden"
import { Stack } from "@mantine/core"
import { Title } from "@/components/atoms/Title"
import { AdminTeam } from "@/components/admin/AdminTeam"

export const Admin = () => {
  const { user } = useAuth()
  if (!user?.admin) return <Forbidden />

  return (
    <Stack>
      <Title order={2} className="font-bold">Admin Settings</Title>
      <p className="text-secondary mb-4">Manage global platform configuration, permissions and admin access.</p>

      <AdminTeam />
    </Stack>
  )
}
