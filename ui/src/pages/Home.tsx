import { Title } from "@/components/atoms/Title"
import { Demo } from "@/components/demo/Demo"
import { FragtapeIcon } from "@/components/icons/FragtapeIcon"
import { useDemoGetAll, useDemoUpload } from "@/lib/api/demo"
import { getErrorMessage } from "@/lib/utils"
import { Button, Center, FileButton, Group, Stack } from "@mantine/core"
import { notifications } from "@mantine/notifications"
import { useMemo, useState } from "react"
import { LuCircleCheckBig, LuClock } from "react-icons/lu"

export const Home = () => {
  const { data: demos, isLoading } = useDemoGetAll()
  const demoUpload = useDemoUpload()

  const [uploading, setUploading] = useState(false)

  const handleUpload = (file: File | null) => {
    if (!file) return

    setUploading(true)

    demoUpload.mutateAsync(file, {
      onSuccess: () => notifications.show({ title: "Demo uploaded", message: "Come back later to see your highlights" }),
      onError: async error => {
        const msg = await getErrorMessage(error)
        notifications.show({ color: "red", message: msg })
      },
      onSettled: () => setUploading(false)
    })
  }

  const content = useMemo(() => {
    if (isLoading) return (
      <Center className="mt-48">
        <FragtapeIcon animated className="size-36 text-(--mantine-color-primary-6)" />
      </Center>
    )
    if (demos?.length === 0) return <NoDemos />
    return (
      <>
        {demos?.map(d => <Demo key={d.id} demo={d} />)}
      </>
    )
  }, [demos, isLoading])

  return (
    <Stack>
      <Group justify="space-between">
        <Title order={2} className="font-bold">Recent Matches</Title>
        <FileButton onChange={handleUpload}>
          {props => <Button loading={uploading} {...props}>Upload</Button>}
        </FileButton>
      </Group>
      {content}
    </Stack>
  )
}

const NoDemos = () => {
  return (
    <Stack align="center" className="text-center">
      <div className="p-4 rounded-full bg-(--mantine-color-background-8)">
        <LuClock className="text-primary size-6" />
      </div>
      <p className="text-primary font-bold">No matches tracked yet</p>
      <p className="text-secondary whitespace-pre-wrap text-balance max-w-xl">{`Once you finish a match with this Steam account it will automatically be visible here.\nOr you can manually upload an old demo`}</p>
      <Stack gap={0} align="start" pl="xl">
        <Group>
          <LuCircleCheckBig className="text-(--mantine-color-primary-6)" />
          <p className="text-secondary">Make sure your next game code is correct</p>
        </Group>
        <Group>
          <LuCircleCheckBig className="text-(--mantine-color-primary-6)" />
          <p className="text-secondary">{`Play a match, we'll handle the rest automatically`}</p>
        </Group>
      </Stack>
    </Stack>
  )
}
