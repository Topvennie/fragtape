import { BottomOfPage } from "@/components/atoms/ButtomOfPage";
import { Card } from "@/components/atoms/Card";
import { Loading } from "@/components/atoms/Loading";
import { Page, PageTitle, Section } from "@/components/atoms/Page";
import { Demo, HighlightFilter } from "@/components/demo/Demo";
import { DemoFilter } from "@/components/demo/DemoFilter";
import { SettingNoConnections } from "@/components/setting/SettingNoConnections";
import { useDemoGetFiltered, useDemoUpload } from "@/lib/api/demo";
import { DemoFilter as DemoFilterType } from "@/lib/types/demo";
import { useSettingGlobalGet } from "@/lib/api/setting_global";
import { useSettingUserGet } from "@/lib/api/setting_user";
import { getErrorMessage } from "@/lib/utils";
import { Button, Center, FileButton, Group, Stack } from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { useState } from "react";
import { LuCircleCheckBig, LuClock, LuFilter } from "react-icons/lu";
import useInfiniteScroll from "react-infinite-scroll-hook";

export const Overview = () => {
  const { data: settingGlobal } = useSettingGlobalGet()
  const { data: settingUser } = useSettingUserGet()

  const [uploading, setUploading] = useState(false)

  // Preload first demo data
  const { result, isLoading } = useDemoGetFiltered({})
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
      {!settingUser?.connectedSteam && <SettingNoConnections />}

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

const Demos = () => {
  const [demoFilter, setDemoFilter] = useState<DemoFilterType>({})
  const [highlightFilter, setHighlightFilter] = useState<HighlightFilter>("me")

  const { result, isLoading, isFetchingNextPage, hasNextPage, fetchNextPage } = useDemoGetFiltered(demoFilter)
  const demos = result.demos

  const [sentryRef] = useInfiniteScroll({
    loading: isFetchingNextPage,
    hasNextPage: Boolean(hasNextPage),
    onLoadMore: fetchNextPage,
    rootMargin: "0px",
  });

  return (
    <>
      <DemoFilter highlight={highlightFilter} setHighlight={setHighlightFilter} demo={demoFilter} setDemo={setDemoFilter} loading={isLoading} total={result.total} />

      <Section
        card={false}
      >

        {demos.length > 0
          ? demos.map(d => <Card key={d.id}><Demo demo={d} highlightFilter={highlightFilter} /></Card>)
          : !isLoading && <NoDemosFiltered clearFilters={() => setDemoFilter({})} />
        }

        <Loading isFetchingNextPage={isFetchingNextPage} hasNextPage={hasNextPage} />

        <BottomOfPage ref={sentryRef} showLoading={isFetchingNextPage} hasNextPage={hasNextPage} />
      </Section>
    </>
  )
}

// The user has demos but not with the current filters
const NoDemosFiltered = ({ clearFilters }: { clearFilters: () => void }) => {
  return (
    <Card>
      <Stack align="center" className="text-center text-secondary">
        <LuFilter className="text-primary size-6" />
        <p className="text-primary font-bold">No demos match your current filters</p>
        <p className="whitespace-pre-wrap text-balance max-w-xl">{`We couldn't find any matches with the filters you selected\nTry loosening your filters to include more demos`}</p>
        <Button onClick={clearFilters} className="mt-2">
          Clear Filters
        </Button>
      </Stack>
    </Card>
  )
}

// The user doesn't have any demos yet
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
