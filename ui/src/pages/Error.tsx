import { isResponseNot200Error } from "@/lib/api/query";
import { useAuth } from "@/lib/hooks/useAuth";
import { Button, Center, Container, Stack, Title } from "@mantine/core";
import { ErrorComponentProps, useNavigate } from "@tanstack/react-router";
import { Error404 } from "./404";
import Smoke from "@/assets/smoke.webm"
import { LuArrowLeft } from "react-icons/lu";
import { FragtapeIcon } from "@/components/icons/FragtapeIcon";

export const Error = ({ error, reset }: ErrorComponentProps) => {
  const { logout } = useAuth()
  const navigate = useNavigate()

  if (isResponseNot200Error(error)) {
    switch (error.response.status) {
      case 404:
        return (
          <Container fluid className="pt-[10%]">
            <Error404 />
          </Container>
        )
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
    <div className="relative w-screen h-screen overflow-hidden bg-[url(/src/assets/smoke.webm)]">
      <video
        autoPlay
        loop
        muted
        playsInline
        preload="auto"
        className="absolute inset-0 h-full w-full object-cover"
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
