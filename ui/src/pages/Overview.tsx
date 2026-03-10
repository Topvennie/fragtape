import { BottomOfPage } from "@/components/atoms/ButtomOfPage"
import { Card } from "@/components/atoms/Card"
import { Loading } from "@/components/atoms/Loading"
import { LoadingOverlay } from "@/components/atoms/LoadingOverlay"
import { Page, PageTitle, Section } from "@/components/atoms/Page"
import { Demo, HighlightFilter } from "@/components/demo/Demo"
import { Segment } from "@/components/molecules/Segment"
import { SettingNoConnections } from "@/components/setting/SettingNoConnections"
import { useDemoGetFiltered, useDemoUpload } from "@/lib/api/demo"
import { useSettingGlobalGet } from "@/lib/api/setting_global"
import { useSettingUserGet } from "@/lib/api/setting_user"
import { DemoFilter, DemoSource } from "@/lib/types/demo"
import { Result, resultString } from "@/lib/types/stat"
import { getErrorMessage } from "@/lib/utils"
import { Button, Center, Collapse, FileButton, Group, Stack } from "@mantine/core"
import { notifications } from "@mantine/notifications"
import { Dispatch, SetStateAction, useState } from "react"
import { LuCircleCheckBig, LuClock, LuEyeOff, LuFilter } from "react-icons/lu"
import useInfiniteScroll from "react-infinite-scroll-hook"

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
  const [demoFilter, setDemoFilter] = useState<DemoFilter>({})
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
      <DemosFilter highlight={highlightFilter} setHighlight={setHighlightFilter} demo={demoFilter} setDemo={setDemoFilter} loading={isLoading} total={result.total} />

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

type DemosFilterProps = {
  highlight: HighlightFilter;
  setHighlight: Dispatch<SetStateAction<HighlightFilter>>;
  demo: DemoFilter;
  setDemo: Dispatch<SetStateAction<DemoFilter>>;
  loading: boolean;
  total: number;
}

type DemosFilterType = "highlight" | "demo" | "none"

type DemoFilterKeys = keyof DemoFilter
type DemoFilterValues = DemoFilter[DemoFilterKeys]

const DemosFilter = ({ highlight, setHighlight, demo, setDemo, loading, total }: DemosFilterProps) => {
  const [filter, setFilter] = useState<DemosFilterType>("highlight")

  const { data: settings } = useSettingGlobalGet()

  const sources = [
    { value: "all", label: "All" },
    { value: DemoSource.Steam, label: "Steam" },
  ]
  if (settings?.demoUpload) sources.push({ value: DemoSource.Manual, label: "Manual" })

  const handleDemoChange = (k: DemoFilterKeys, v: DemoFilterValues) => {
    setDemo(prev => ({ ...prev, [k]: v }))
  }

  return (
    <Section
      title={`Filters`}
      card={false}
      rightSection={(
        <Group gap={0} justify="end">
          <Segment
            data={[
              { value: "highlight", label: "Clips" },
              { value: "demo", label: "Matches" },
              { value: "none", label: <LuEyeOff className="size-5" /> }
            ]}
            value={filter}
            onChange={e => setFilter(e as DemosFilterType)}
          />
          <div className="text-white min-w-[11ch] text-right">
            {!loading && <p>{`${total} Matche${total !== 1 ? 's' : ''}`}</p>}
          </div>
        </Group>
      )}
    >
      <Collapse in={filter !== "none"}>
        <Card className="relative">
          <LoadingOverlay loading={loading} />
          {filter === "highlight" && (
            <Stack align="flex-start">
              <Segment
                data={[
                  { value: "me", label: "My clips" },
                  { value: "group", label: "Group" },
                ]}
                value={highlight}
                onChange={e => setHighlight(e as HighlightFilter)}
                label="Visible Clips"
              />
            </Stack>
          )}

          {filter === "demo" && (
            <Stack align="flex-start">
              <Segment
                data={sources}
                value={!demo.source ? "all" : demo.source}
                onChange={e => handleDemoChange("source", e !== "all" ? e as DemoSource : undefined)}
                label="Source"
              />
              <Segment
                data={[
                  { value: "all", label: "All" },
                  ...Object.values(Result).map(r => ({ value: r, label: resultString[r] })),
                ]}
                value={!demo.result ? "all" : demo.result}
                onChange={e => handleDemoChange("result", e !== "all" ? e as Result : undefined)}
                label="Result"
              />
            </Stack>
          )}
        </Card>
      </Collapse>
    </Section>
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
