import { Button, Center, Stack, Text, Title } from "@mantine/core";
import { useNavigate } from "@tanstack/react-router";

export const Forbidden = () => {
  const navigate = useNavigate()

  const handleReturn = () => {
    navigate({ to: "." })
  }

  return (
    <Center h="100%">
      <Stack align="center" gap={0}>
        <Text fw={600}>403</Text>
        <Title fw={600} className="mt-2">Forbidden</Title>
        <Text c="gray" className="mt-6">Not enough permissions</Text>
        <Button onClick={handleReturn} className="mt-6">
          Return to the start
        </Button>
      </Stack>
    </Center>
  );
}
