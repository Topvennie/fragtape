import { useAuth } from "@/lib/hooks/useAuth"
import { DemoStatus, Demo as DemoType } from "@/lib/types/demo"
import { Highlight } from "@/lib/types/highlight"
import { Result, resultString } from "@/lib/types/stat"
import { cn, formatDate } from "@/lib/utils"
import { Button, Collapse, Group, Stack } from "@mantine/core"
import { useMediaQuery } from "@mantine/hooks"
import { ReactNode, useEffect, useMemo, useState } from "react"
import { LuChevronDown, LuClapperboard, LuTriangleAlert } from "react-icons/lu"
import { Card } from "../atoms/Card"
import { FragtapeIcon } from "../icons/FragtapeIcon"
import { DemoThumbnail } from "./DemoThumbnail"
import { HighlightCarousel } from "../highlight/HighlightCarousel"

type Props = {
  demo: DemoType;
  highlightFilter?: HighlightFilter;
}

export type HighlightFilter = "me" | "group"

const resultColor: Record<Result, string> = {
  [Result.Win]: "text-green-400",
  [Result.Loss]: "text-red-400",
  [Result.Tie]: "text-white",
}

export const Demo = ({ demo, highlightFilter }: Props) => {
  const content = useMemo(() => {
    switch (demo.status) {
      case DemoStatus.QueuedDownload:
        return <DownloadQueued />
      case DemoStatus.Downloading:
        return <Download />
      case DemoStatus.QueuedParse:
        return <ParseQueued />
      case DemoStatus.Parsing:
        return <Parse />
      case DemoStatus.QueuedRender:
      case DemoStatus.Rendering:
      case DemoStatus.QueuedFinalize:
      case DemoStatus.Finalize:
      case DemoStatus.Finished:
        return <Finished demo={demo} highlightFilter={highlightFilter} />
      case DemoStatus.Failed:
        return <Failed />
    }
  }, [demo, highlightFilter])

  return (
    <Card>
      {content}
    </Card>
  )
}

const Finished = ({ demo, highlightFilter = "me" }: Props) => {
  const { user } = useAuth()
  const [showClips, setShowClips] = useState(false)

  useEffect(() => setShowClips(false), [highlightFilter])

  const mdPoint = useMediaQuery('(width >= 48em)')

  const player = demo.players.find(p => p.user.id === user?.id)
  if (!player) return null // Shouldn't really be possible

  const highlights = demo.players.flatMap(p => p.highlights)
  const filteredHighlights = highlightFilter === "me" ? player?.highlights ?? [] : highlights

  const score = () => {
    const winnerRounds = Math.max(demo.stats.roundsCt, demo.stats.roundsT)
    const loserRounds = Math.min(demo.stats.roundsCt, demo.stats.roundsT)

    if (player.stat.result === Result.Win)
      return `${winnerRounds} - ${loserRounds}`

    return `${loserRounds} - ${winnerRounds}`
  }

  return (
    <Stack>
      <Group align="stretch" wrap="nowrap">
        <div className="w-16 xs:w-32 sm:w-42 lg:w-64 aspect-square sm:aspect-video shrink-0 rounded-md overflow-hidden h-fit">
          <DemoThumbnail demo={demo} />
        </div>
        <Stack gap={0} w="100%" justify="space-between">
          <Group justify="space-between" align="start">
            <Stack gap={0}>
              <Group>
                <p className={`text-lg sm:text-xl lg:text-2xl font-bold uppercase ${resultColor[player.stat.result]}`}>{resultString[player.stat.result]}</p>
                <p className="sm:text-lg lg:text-xl text-white">{score()}</p>
              </Group>
              <Group className="text-secondary">
                <p>K <span className="text-white">{player.stat.kills}</span></p>
                <p>D <span className="text-white">{player.stat.deaths}</span></p>
                <p>A <span className="text-white">{player.stat.assists}</span></p>
              </Group>
            </Stack>
            {mdPoint && <p className="text-secondary">{formatDate(demo.playedAt)}</p>}
          </Group>
          {mdPoint && (
            <Group justify="space-between">
              <ClipBadge demo={demo} highlights={highlights} filtered={filteredHighlights} className="mt-auto" />
              {filteredHighlights.length > 0 && (
                <Button variant="subtle" color="muted" onClick={() => setShowClips(prev => !prev)} rightSection={<LuChevronDown className={`transform duration-300 ${showClips ? "rotate-180" : ""}`} />}>
                  {`${showClips ? "Hide" : "Show"} clips`}
                </Button>
              )}
            </Group>
          )}
        </Stack>
      </Group>
      <Collapse in={showClips}>
        <div className="pt-8">
          {/* Prevent the carousel to preload the video's unless asked */}
          {showClips && <HighlightCarousel highlights={filteredHighlights} />}
        </div>
      </Collapse>
    </Stack>
  )
}

const DownloadQueued = () => {
  return <Loading text="Download queued" />
}

const Download = () => {
  return <Loading text="Downloading match..." />
}

const ParseQueued = () => {
  return <Loading text="Parsing queued" />
}

const Parse = () => {
  return <Loading text="Parsing match..." />
}

const Loading = ({ text }: { text: string }) => {
  return (
    <div className="flex items-center gap-4">
      <FragtapeIcon animated={true} className="size-8 text-(--mantine-color-primary-6)" />
      <div className="flex flex-col gap-2 justify-center">
        <p className="text-xl font-bold text-secondary">{text}</p>
      </div>
    </div>
  )
}

const Failed = () => {
  return (
    <div className="flex items-center gap-4">
      <LuTriangleAlert className="size-8 text-red-400" />
      <div className="flex flex-col justify-center">
        <p className="text-xl font-bold text-red-400">Highlight generation failed</p>
        <p className="text-secondary text-sm">{`We couldn't process this match`}</p>
      </div>
    </div>
  )
}

const ClipBadge = ({ demo, highlights, filtered, className }: { demo: DemoType, highlights: Highlight[], filtered: Highlight[], className: string }) => {
  if (demo.status === DemoStatus.Finished && highlights.length === 0) {
    return null
  }

  const rendering = filtered.some(h => !h.generated)

  let icon: ReactNode

  if (rendering) {
    icon = <FragtapeIcon animated className="size-4 text-(--mantine-color-primary-6)" />
  } else {
    icon = <LuClapperboard className="text-(--mantine-color-primary-6)" />
  }

  return (
    <Group className={cn("bg-(--mantine-color-primary-light) rounded-lg py-1 lg:py-2 px-2 lg:px-4 w-fit", className)}>
      {icon}
      <p className="text-white text-sm">{filtered.length} <span className="text-[0.65rem] align-top text-secondary">{`/ ${highlights.length}`}</span>{` Clip${highlights.length !== 1 ? 's' : ''}`}</p>
    </Group>
  )
}
