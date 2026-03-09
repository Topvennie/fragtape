import Smoke from "@/assets/smoke.webm";
import { FragtapeIcon } from "@/components/icons/FragtapeIcon";
import { isResponseNot200Error } from "@/lib/api/query";
import { useAuth } from "@/lib/hooks/useAuth";
import { Button, Center, Stack, Title } from "@mantine/core";
import { ErrorComponentProps, useNavigate } from "@tanstack/react-router";
import { LuArrowLeft } from "react-icons/lu";
import { Error404 } from "./404";

export const Error = ({ error, reset }: ErrorComponentProps) => {
  const { logout } = useAuth()
  const navigate = useNavigate()

  if (isResponseNot200Error(error)) {
    switch (error.response.status) {
      case 404: return <Error404 />
      case 401:
        logout()
        navigate({ to: "/" })
        break
    }
  }

  const handleReturn = () => {
    reset()
    navigate({ to: "/" })
  }

  return (
    <div className="fixed top-0 left-0 z-100 h-screen w-screen bg-[url(/src/assets/smoke.webm)]">
      <video
        autoPlay
        loop
        muted
        playsInline
        preload="auto"
        className="absolute object-cover h-full"
      >
        <source src={Smoke} type="video/webm" />
      </video>
      <Center h="100vh" className="relative z-10 text-primary">
        <Stack align="center">
          <FragtapeIcon className="my-8 size-14 text-(--mantine-color-primary-6) animate-pulse-extreme" />
          <Title order={2} className="text-center">Server got smoked</Title>
          <p className="text-secondary whitespace-pre-wrap text-center">{`We can't see a thing right now.\nThe connection is smoked.`}</p>
          <Button onClick={handleReturn} leftSection={<LuArrowLeft />}>Go back</Button>
        </Stack>
      </Center>
    </div>
  )
}
