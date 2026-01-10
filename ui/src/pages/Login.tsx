import { Card } from "@/components/atoms/Card"
import { Title } from "@/components/atoms/Title"
import { FragtapeIcon } from "@/components/icons/FragtapeIcon"
import { GithubIcon } from "@/components/icons/GithubIcon"
import { SteamIcon } from "@/components/icons/SteamIcon"
import { useAuth } from "@/lib/hooks/useAuth"
import { ActionIcon, Button, Center, Stack } from "@mantine/core"

export const Login = () => {
  const { login } = useAuth()

  return (
    <Center h="100vh">
      <Card className="max-w-lg">
        <Stack>
          <FragtapeIcon className="text-(--mantine-color-primary-6) h-14 aspect-square my-8" />
          <Stack align="center" gap="xl">
            <Title order={2} className="text-center">Start generating highlights</Title>
            <p className="text-secondary text-center">Connect your Steam account to automatically fetch matches and create clips from your best moments.</p>
            <Button onClick={login} fullWidth color="steam" leftSection={<SteamIcon className="h-5 w-5" />}>
              Log in with Steam
            </Button>
            <ActionIcon onClick={e => e.preventDefault()} component="a" href="https://github.com/topvennie/fragtape" rel="noopener noreferrer" target="_blank" variant="transparent">
              <GithubIcon className="text-(--mantine-color-white)" />
            </ActionIcon>
          </Stack>
        </Stack>
      </Card>
    </Center>
  )
}
