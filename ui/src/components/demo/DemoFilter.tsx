import { useSettingGlobalGet } from "@/lib/api/setting_global";
import { DemoFilter as DemoFilterType, DemoHighlightFilter, DemoSource } from "@/lib/types/demo";
import { Result, resultString } from "@/lib/types/stat";
import { Group } from "@mantine/core";
import { DatesRangeValue } from "@mantine/dates";
import { useLocalStorage } from "@mantine/hooks";
import NumberFlow from '@number-flow/react';
import { Dispatch, SetStateAction, useLayoutEffect, useMemo, useRef, useState } from "react";
import { LuEyeOff } from "react-icons/lu";
import { Card } from "../atoms/Card";
import { LoadingOverlay } from "../atoms/LoadingOverlay";
import { Section } from "../atoms/Page";
import { DatePickerInput } from "../molecules/DatePickerInput";
import { Segment } from "../molecules/Segment";

type Props = {
  highlight: DemoHighlightFilter;
  setHighlight: Dispatch<SetStateAction<DemoHighlightFilter>>;
  demo: DemoFilterType;
  setDemo: Dispatch<SetStateAction<DemoFilterType>>;
  loading: boolean;
  total: number;
}

type DemosFilterType = "highlight" | "demo" | "none"

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
          <div className={`min-w-[11ch] text-right text-white ${loading ? "invsible" : "block"}`}>
            <NumberFlow value={total} suffix={` Match${total !== 1 ? "es" : ""}`} />
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

type HighlightKeys = keyof DemoHighlightFilter
type HighlightValues = DemoHighlightFilter[HighlightKeys]

const Highlight = ({ highlight, setHighlight }: Pick<Props, "highlight" | "setHighlight">) => {
  const handleHighlightChange = (k: HighlightKeys, v: HighlightValues) => {
    setHighlight(prev => ({ ...prev, [k]: v }))
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
      <Segment
        segmentProps={{
          data: [
            { value: "me", label: "My clips" },
            { value: "group", label: "Group" },
          ],
          value: highlight.player,
          onChange: (e) => handleHighlightChange("player", e as "me" | "group"),
        }}
        label="Visible Clips"
      />
      <Segment
        segmentProps={{
          data: [
            { value: "all", label: "Any" },
            ...Array.from({ length: 3 }).map((_, idx) => ({ value: (idx + 3).toString(), label: (idx + 3).toString() }))
          ],
          value: highlight.minKillCount ? highlight.minKillCount.toString() : "all",
          onChange: (e) => handleHighlightChange("minKillCount", e !== "all" ? Number(e) : undefined),
        }}
        label="Minimum kill count"
        labelProps={{
          className: "md:text-right"
        }}
        className="md:justify-self-end"
      />
    </div>
  )
}

type DemoKeys = keyof DemoFilterType
type DemoValues = DemoFilterType[DemoKeys]

const Demo = ({ demo, setDemo }: Pick<Props, "demo" | "setDemo">) => {
  const { data: settings } = useSettingGlobalGet()
  const sources = [
    { value: "all", label: "All sources" },
    { value: DemoSource.Steam, label: "Steam" },
  ]
  if (settings?.demoUpload) sources.push({ value: DemoSource.Manual, label: "Manual" })

  const handleDemoChange = (k: DemoKeys, v: DemoValues) => {
    setDemo(prev => ({ ...prev, [k]: v }))
  }

  const handleDateChange = (r: DatesRangeValue) => {
    if (r[0] && r[1] || (!r[0] && !r[1])) {
      r[0]?.setHours(0)
      r[0]?.setMinutes(0)
      r[0]?.setSeconds(0)

      r[1]?.setHours(23)
      r[1]?.setMinutes(59)
      r[1]?.setSeconds(59)

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

