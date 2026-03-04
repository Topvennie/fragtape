import { Alert } from "@/components/atoms/Alert"
import { LinkButton } from "@/components/atoms/LinkButton"
import { Title } from "@/components/atoms/Title"
import { Demo, HighlightFilter } from "@/components/demo/Demo"
import { FragtapeIcon } from "@/components/icons/FragtapeIcon"
import { Segment } from "@/components/molecules/Segment"
import { useDemoGetAll, useDemoUpload } from "@/lib/api/demo"
import { useSettingGlobalGet } from "@/lib/api/setting_global"
import { useSettingUserGet } from "@/lib/api/setting_user"
import { Demo as DemoType } from "@/lib/types/demo"
import { getErrorMessage } from "@/lib/utils"
import { Button, Center, FileButton, Group, Stack } from "@mantine/core"
import { notifications } from "@mantine/notifications"
import { useMemo, useState } from "react"
import { LuArrowRight, LuCircleCheckBig, LuClock, LuTriangleAlert } from "react-icons/lu"

export const Home = () => {
  const { data: demos, isLoading: isLoadingDemos } = useDemoGetAll()
  const demoUpload = useDemoUpload()

  const [uploading, setUploading] = useState(false)

  const { data: settingGlobal } = useSettingGlobalGet()
  const { data: settingUser } = useSettingUserGet()

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
    if (isLoadingDemos) return (
      <Center className="mt-48">
        <FragtapeIcon animated className="size-36 text-(--mantine-color-primary-6)" />
      </Center>
    )
    if (demos?.length === 0) return <NoDemos />
    return <Demos demos={demos ?? []} />
  }, [demos, isLoadingDemos])

  return (
    <Stack>
      {!settingUser?.connectedSteam && <NoConnections />}
      <Group justify="space-between">
        <Title order={2} className="font-bold">Recent Matches</Title>
        {settingGlobal?.demoUpload && (
          <FileButton onChange={handleUpload}>
            {props => <Button loading={uploading} {...props}>Upload</Button>}
          </FileButton>
        )}
      </Group>
      {content}
    </Stack>
  )
}

const NoConnections = () => {
  return (
    <Alert
      title="Missing Account Connections"
      icon={<LuTriangleAlert className="size-6 text-(--mantine-color-primary-6)" />}
      color="orange"
      border
    >Round 3 · 00:33
      No accounts are linked to your profile. We cannot fetch your matches or generate highlights until you connect at least one account.
      <div className="mt-2">
        <LinkButton to="/setting" rightSection={<LuArrowRight />}>
          Go to Settings
        </LinkButton>
      </div>
    </Alert>
  )
}

const NoDemos = () => {
  return (
    <Stack align="center" className="text-center">
      <div className="p-4 rounded-full bg-(--mantine-color-background-8)">
        <LuClock className="text-primary size-6" />
      </div>
      <p className="text-primary font-bold">No matches tracked yet</p>
      <p className="text-secondary whitespace-pre-wrap text-balance max-w-xl">{`Once you finish a match with this Steam account it will automatically be visible here.`}</p>
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

const Demos = ({ demos }: { demos: DemoType[] }) => {
  const [highlightFilter, setHighlightFilter] = useState<HighlightFilter>("me")

  const handleFilterChange = (value: string) => {
    setHighlightFilter(value as HighlightFilter)
  }

  return (
    <>
      <Segment
        data={[
          { value: "me", label: "Only my clips" },
          { value: "group", label: "Me + group" },
          { value: "match", label: "Everyone" },
        ]}
        value={highlightFilter}
        onChange={handleFilterChange}
        className="ml-auto"
      />
      {demos.map(d => <Demo key={d.id} demo={d} highlightFilter={highlightFilter} />)}
    </>
  )
}

