import { Title } from "@/components/atoms/Title";
import { useAuth } from "@/lib/hooks/useAuth";
import { Button, Center, Stack } from "@mantine/core";
import { useNavigate } from "@tanstack/react-router";
import { LuArrowLeft, LuLock } from "react-icons/lu";

export const Forbidden = () => {
  const { logout } = useAuth()
  const navigate = useNavigate()

  const handleReturn = () => {
    logout()
    navigate({ to: "/" })
  }

  return (
    <Center h="100vh">
      <Stack align="center" gap="xl">
        <div className="p-4 rounded-full bg-(--mantine-color-primary-2) border-2 border-(--mantine-color-primary-4)">
          <LuLock className="size-16 text-(--mantine-color-primary-6)" />
        </div>
        <Title order={2} className="text-center">Restricted Area</Title>
        <p className="text-secondary whitespace-pre-wrap text-center">{`You don't have access to this area.`}</p>
        <Button onClick={handleReturn} leftSection={<LuArrowLeft />}>Go back home</Button>
      </Stack>
    </Center>
  )
}
