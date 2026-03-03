import { useAuth } from "@/lib/hooks/useAuth"
import { DemoStatus, Demo as DemoType } from "@/lib/types/demo"
import { Highlight } from "@/lib/types/highlight"
import { Result, resultString } from "@/lib/types/stat"
import { formatDate } from "@/lib/utils"
import { Button, Collapse } from "@mantine/core"
import { useMediaQuery } from "@mantine/hooks"
import { ReactNode, useMemo, useState } from "react"
import { LuChevronDown, LuClapperboard, LuTriangleAlert } from "react-icons/lu"
import { Card } from "../atoms/Card"
import { HighlightCarousel } from "../highlight/HighlightCarousel"
import { FragtapeIcon } from "../icons/FragtapeIcon"
import { DemoThumbnail } from "./DemoThumbnail"

type Props = {
  demo: DemoType
}

const resultColor: Record<Result, string> = {
  [Result.Win]: "text-green-400",
  [Result.Loss]: "text-red-400",
  [Result.Tie]: "text-white",
}

export const Demo = ({ demo }: Props) => {
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
        return <Finished demo={demo} />
      case DemoStatus.Failed:
        return <Failed />
    }
  }, [demo])

  return (
    <Card>
      {content}
    </Card>
  )
}

const Finished = ({ demo }: Props) => {
  const { user } = useAuth()
  const [showClips, setShowClips] = useState(false)

  const smPoint = useMediaQuery('(width >= 40em)')

  const highlights = useMemo(() => demo.players.flatMap(p => p.highlights.filter(h => h.generated)), [demo])
  const player = demo.players.find(p => p.user.id === user?.id)
  if (!player) return null // Shouldn't really be possible

  const score = () => {
    const winnerRounds = Math.max(demo.stats.roundsCt, demo.stats.roundsT)
    const loserRounds = Math.min(demo.stats.roundsCt, demo.stats.roundsT)

    if (player.stat.result === Result.Win)
      return `${winnerRounds} - ${loserRounds}`

    return `${loserRounds} - ${winnerRounds}`
  }

  return (
    <div className="flex flex-col">
      <div className="flex items-center gap-4">
        <div className="w-16 xs:w-32 sm:w-42 lg:w-64 aspect-square sm:aspect-video shrink-0 rounded-md overflow-hidden h-fit">
          <DemoThumbnail demo={demo} />
        </div>
        <div className="flex justify-between w-full">
          <div className="flex flex-col gap-2 justify-center">
            <div className="flex items-center gap-4">
              <p className={`text-lg sm:text-xl lg:text-2xl font-bold uppercase ${resultColor[player.stat.result]}`}>{resultString[player.stat.result]}</p>
              <p className="sm:text-lg lg:text-xl text-white">{score()}</p>
            </div>
            <div className="space-x-4 text-secondary">
              K <span className="text-white">{player.stat.kills}</span>
              D <span className="text-white">{player.stat.deaths}</span>
              A <span className="text-white">{player.stat.assists}</span>
            </div>
            <ClipBadge demo={demo} highlights={highlights} />
          </div>
          {smPoint && (
            <div className="flex flex-col items-end justify-between">
              <div>
                <p className="text-secondary">{formatDate(demo.playedAt)}</p>
              </div>
              {highlights.length > 0 && (
                <Button variant="subtle" color="muted" onClick={() => setShowClips(prev => !prev)} rightSection={<LuChevronDown className={`transform duration-300 ${showClips ? "rotate-180" : ""}`} />}>
                  {`${showClips ? "Hide" : "Show"} clips`}
                </Button>
              )}
            </div>
          )}
        </div>
      </div>
      <Collapse in={showClips}>
        <div className="pt-8">
          <HighlightCarousel highlights={highlights} />
        </div>
      </Collapse>
    </div>
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

const ClipBadge = ({ demo, highlights }: { demo: DemoType, highlights: Highlight[] }) => {
  if (DemoStatus.Finished && highlights.length === 0) {
    return null
  }

  const rendering = [DemoStatus.QueuedRender, DemoStatus.Rendering].includes(demo.status)

  let text: string;
  let icon: ReactNode;

  if (rendering) {
    text = "Clips are rendering"
    icon = <FragtapeIcon animated className="size-4 text-(--mantine-color-primary-6)" />
  } else {
    text = `${highlights.length} Clip${highlights.length !== 1 ? 's' : ''} generated`
    icon = <LuClapperboard className="text-(--mantine-color-primary-6)" />
  }

  return (
    <div className="flex items-center gap-2 bg-(--mantine-color-primary-light) rounded-lg py-1 lg:py-2 px-2 lg:px-4 w-fit">
      {icon}
      <p className="text-white text-sm">{text}</p>
    </div>
  )
}
