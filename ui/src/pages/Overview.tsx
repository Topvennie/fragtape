import { Alert } from "@/components/atoms/Alert"
import { BottomOfPage } from "@/components/atoms/ButtomOfPage"
import { LinkButton } from "@/components/atoms/LinkButton"
import { Loading } from "@/components/atoms/Loading"
import { Page, PageTitle, Section } from "@/components/atoms/Page"
import { Demo, HighlightFilter } from "@/components/demo/Demo"
import { Segment } from "@/components/molecules/Segment"
import { useDemoGetFiltered, useDemoUpload } from "@/lib/api/demo"
import { useSettingGlobalGet } from "@/lib/api/setting_global"
import { useSettingUserGet } from "@/lib/api/setting_user"
import { getErrorMessage } from "@/lib/utils"
import { Button, Center, FileButton, Group, Stack } from "@mantine/core"
import { notifications } from "@mantine/notifications"
import { useState } from "react"
import { LuArrowRight, LuCircleCheckBig, LuClock, LuFilter, LuTriangleAlert } from "react-icons/lu"
import useInfiniteScroll from "react-infinite-scroll-hook"

export const Overview = () => {
  const { data: settingGlobal } = useSettingGlobalGet()
  const { data: settingUser } = useSettingUserGet()

  const [uploading, setUploading] = useState(false)

  // Preload first demo data
  const { result, isLoading } = useDemoGetFiltered()
  const demoUpload = useDemoUpload()

  if (isLoading) return <Loading />

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

  return (
    <Page>
      {!settingUser?.connectedSteam && <NoConnections />}

      <PageTitle
        title="Recent Matches"
        rightSection={settingGlobal?.demoUpload && (
          <FileButton onChange={handleUpload}>
            {props => <Button loading={uploading} {...props}>Upload</Button>}
          </FileButton>
        )}
      />

      {result.demos.length > 0
        ? <Demos />
        : <NoDemos />
      }

    </Page>
  )
}

const NoConnections = () => {
  return (
    <Alert
      title="Missing Account Connections"
      icon={<LuTriangleAlert className="size-6 text-(--mantine-color-primary-6)" />}
      color="orange"
      border
    >
      {`Your profile has no connections.\nWe cannot fetch your matches or generate highlights until you add at least one connection.`}
      <div className="mt-4">
        <LinkButton to="/setting" rightSection={<LuArrowRight />}>
          Go to Settings
        </LinkButton>
      </div>
    </Alert>
  )
}

const Demos = () => {
  const { result, isFetchingNextPage, hasNextPage, fetchNextPage } = useDemoGetFiltered()
  const demos = result.demos

  const [sentryRef] = useInfiniteScroll({
    loading: isFetchingNextPage,
    hasNextPage: Boolean(hasNextPage),
    onLoadMore: fetchNextPage,
    rootMargin: "0px",
  });

  const [highlightFilter, setHighlightFilter] = useState<HighlightFilter>("me")

  const handleFilterChange = (value: string) => {
    setHighlightFilter(value as HighlightFilter)
  }

  if (demos.length === 0) return <NoDemos />

  return (
    <Section
      rightSection={(
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
      )}
    >
      {demos.length > 0
        ? demos.map(d => <Demo key={d.id} demo={d} highlightFilter={highlightFilter} />)
        : <NoDemosFiltered />
      }

      <Loading isFetchingNextPage={isFetchingNextPage} hasNextPage={hasNextPage} />

      <BottomOfPage ref={sentryRef} showLoading={isFetchingNextPage} hasNextPage={hasNextPage} />
    </Section>
  )
}

const NoDemosFiltered = () => {
  return (
    <Stack align="center" className="text-center text-secondary">
      <LuFilter className="text-primary size-6" />
      <p className="text-primary font-bold">No demos match your current filters</p>
      <p className="whitespace-pre-wrap text-balance max-w-xl">{`We couldn't find any matches with the filters you selected\nTry loosening your filters to include more demos`}</p>
      <Button className="mt-2">
        Clear Filters
      </Button>
    </Stack>
  )
}

const NoDemos = () => {
  return (
    <Section align="center">
      <Center>
        <Stack align="center" className="text-center text-secondary">
          <LuClock className="text-primary size-6" />
          <p className="text-primary font-bold">No matches tracked yet</p>
          <p className="whitespace-pre-wrap text-balance max-w-xl">{`Once you finish a match with this Steam account it will automatically be visible here.`}</p>
          <Stack gap={0} align="start" pl="xl">
            <Group>
              <LuCircleCheckBig className="text-(--mantine-color-primary-6)" />
              <p>Make sure you have an account linked</p>
            </Group>
            <Group>
              <LuCircleCheckBig className="text-(--mantine-color-primary-6)" />
              <p>{`Play a match, we handle the rest automatically`}</p>
            </Group>
          </Stack>
        </Stack>
      </Center>
    </Section>
  )
}
