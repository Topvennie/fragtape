import { Title } from "@/components/atoms/Title"
import { FragtapeIcon } from "@/components/icons/FragtapeIcon"
import { Button, Center, Stack } from "@mantine/core"
import { useNavigate } from "@tanstack/react-router"
import { LuArrowLeft, LuUserRound } from "react-icons/lu"

export const Error404 = () => {
  const navigate = useNavigate()

  const handleReturn = () => navigate({ to: "/" })

  return (
    <Center h="100vh">
      <Stack align="center" gap="xl">
        <div className="flex justify-center relative w-16 sm:w-32 lg:w-64 aspect-square">
          <FragtapeIcon className="text-(--mantine-color-primary-6) size-16 sm:size-32 lg:size-64  z-20 absolute left-7 -top-5 sm:left-15 sm:-top-10 lg:left-25 lg:-top-15" />
          <LuUserRound className="size-16 sm:size-32 lg:size-64 z-10 stroke-1 text-gray-500" />
        </div>
        <Title order={2} className="text-center">You missed</Title>
        <p className="text-secondary whitespace-pre-wrap text-center">{`Looks like the shot went wide.\nThe page you are looking for doesn't exist or has been removed.`}</p>
        <Button onClick={handleReturn} leftSection={<LuArrowLeft />}>Go back home</Button>
      </Stack>
    </Center>
  )
}
