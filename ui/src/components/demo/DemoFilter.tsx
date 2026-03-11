import { DemoFilter as DemoFilterType, DemoSource } from "@/lib/types/demo";
import { Dispatch, SetStateAction, useLayoutEffect, useMemo, useRef, useState } from "react";
import { HighlightFilter } from "./Demo";
import { useSettingGlobalGet } from "@/lib/api/setting_global";
import { Result, resultString } from "@/lib/types/stat";
import { Group, Stack } from "@mantine/core";
import { useLocalStorage } from "@mantine/hooks";
import { LuEyeOff } from "react-icons/lu";
import { Section } from "../atoms/Page";
import { Segment } from "../molecules/Segment";
import { Card } from "../atoms/Card";
import { LoadingOverlay } from "../atoms/LoadingOverlay";
import { DatePickerInput } from "../molecules/DatePickerInput";
import { DatesRangeValue } from "@mantine/dates";

type Props = {
  highlight: HighlightFilter;
  setHighlight: Dispatch<SetStateAction<HighlightFilter>>;
  demo: DemoFilterType;
  setDemo: Dispatch<SetStateAction<DemoFilterType>>;
  loading: boolean;
  total: number;
}


type DemosFilterType = "highlight" | "demo" | "none"

type Keys = keyof DemoFilterType
type Values = DemoFilterType[Keys]

export const DemoFilter = ({ highlight, setHighlight, demo, setDemo, loading, total }: Props) => {
  const [filter, setFilter] = useLocalStorage({ key: "fragtape-ui-demo-filter", defaultValue: "none" })

  const innerRef = useRef<HTMLDivElement>(null);
  const [height, setHeight] = useState(0);

  useLayoutEffect(() => {
    if (filter === "none") {
      setHeight(0);
      return;
    }

    if (innerRef.current) {
      setHeight(innerRef.current.scrollHeight);
    }
  }, [filter]);

  const content = useMemo(() => {
    switch (filter) {
      case "highlight":
        return <Highlight highlight={highlight} setHighlight={setHighlight} />
      case "demo":
        return <Demo demo={demo} setDemo={setDemo} />
      default:
        return null
    }
  }, [filter, demo, highlight]) // eslint-disable-line react-hooks/exhaustive-deps

  return (
    <Section
      title="Filters"
      card={false}
      rightSection={
        <Group gap={0} justify="end">
          <Segment
            segmentProps={{
              data: [
                { value: "highlight", label: "Clips" },
                { value: "demo", label: "Matches" },
                { value: "none", label: <LuEyeOff className="size-5" /> },
              ],
              value: filter,
              onChange: (e) => setFilter(e as DemosFilterType),
            }}
          />
          <div className="min-w-[11ch] text-right text-white">
            {!loading && <p>{`${total} Match${total !== 1 ? "es" : ""}`}</p>}
          </div>
        </Group>
      }
    >
      <div
        className="overflow-hidden transition-[height] duration-300"
        style={{ height: `${height}px` }}
      >
        <div ref={innerRef}>
          <Card className="relative">
            <LoadingOverlay loading={loading} />

            {content}
          </Card>
        </div>
      </div>
    </Section>
  );
}

const Highlight = ({ highlight, setHighlight }: Pick<Props, "highlight" | "setHighlight">) => {
  return (
    <Stack align="flex-start">
      <Segment
        segmentProps={{
          data: [
            { value: "me", label: "My clips" },
            { value: "group", label: "Group" },
          ],
          value: highlight,
          onChange: (e) => setHighlight(e as HighlightFilter),
        }}
        label="Visible Clips"
      />
    </Stack>
  )
}

const Demo = ({ demo, setDemo }: Pick<Props, "demo" | "setDemo">) => {
  const { data: settings } = useSettingGlobalGet()
  const sources = [
    { value: "all", label: "All sources" },
    { value: DemoSource.Steam, label: "Steam" },
  ]
  if (settings?.demoUpload) sources.push({ value: DemoSource.Manual, label: "Manual" })

  const handleDemoChange = (k: Keys, v: Values) => {
    setDemo(prev => ({ ...prev, [k]: v }))
  }

  const handleDateChange = (r: DatesRangeValue) => {
    if (r[0] && r[1] || (!r[0] && !r[1])) {
      handleDemoChange("playedAtStart", r[0] ?? undefined)
      handleDemoChange("playedAtEnd", r[1] ?? undefined)
    }
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
      <Segment
        segmentProps={{
          data: sources,
          value: !demo.source ? "all" : demo.source,
          onChange: (e) => handleDemoChange("source", e !== "all" ? (e as DemoSource) : undefined)
        }}
        label="Source"
      />
      <Segment
        segmentProps={{
          data: [
            { value: "all", label: "Any result" },
            ...Object.values(Result).map((r) => ({
              value: r,
              label: resultString[r],
            })),
          ],
          value: !demo.result ? "all" : demo.result,
          onChange: (e) => handleDemoChange("result", e !== "all" ? (e as Result) : undefined)
        }}
        label="Result"
        labelProps={{
          className: "md:text-right"
        }}
        className="md:justify-self-end"
      />
      <Segment
        segmentProps={{
          data: [
            { value: "all", label: "Any" },
            { value: "true", label: "Has highlights" },
            { value: "false", label: "No highlights" },
          ],
          value: demo.hasHighlight === undefined ? "all" : demo.hasHighlight.toString(),
          onChange: (e) => handleDemoChange("hasHighlight", e !== "all" ? e === "true" : undefined)
        }}
        label="Highlights in match"
      />
      <DatePickerInput
        dateProps={{
          type: "range",
          value: [demo.playedAtStart ?? null, demo.playedAtEnd ?? null],
          onChange: handleDateChange,
          clearable: true,
          numberOfColumns: 2,
        }}
        label="Time range"
        labelProps={{
          className: "md:text-right"
        }}
        className="md:justify-self-end min-w-72"
      />
    </div>
  )
}

