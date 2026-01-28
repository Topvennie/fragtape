import { useAuth } from "@/lib/hooks/useAuth"
import { DemoStatus, Demo as DemoType } from "@/lib/types/demo"
import { Highlight } from "@/lib/types/highlight"
import { Result, resultString } from "@/lib/types/stat"
import { formatDate } from "@/lib/utils"
import { Button, Center, Collapse } from "@mantine/core"
import { ReactNode, useMemo, useState } from "react"
import { LuChevronDown, LuClapperboard, LuTriangleAlert } from "react-icons/lu"
import { Card } from "../atoms/Card"
import { LoadableImage } from "../atoms/LoadableImage"
import { HighlightCarousel } from "../highlight/HighlightCarousel"
import { FragtapeIcon } from "../icons/FragtapeIcon"
import { DemoThumbnail } from "./DemoThumbnail"
import { useMediaQuery } from "@mantine/hooks"

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
      case DemoStatus.QueuedParse:
      case DemoStatus.Parsing:
        return <Loading />
      case DemoStatus.Failed:
        return <Failed />
      default:
        return <DemoInner demo={demo} />
    }
  }, [demo])

  return (
    <Card>
      {content}
    </Card>
  )

}

const DemoInner = ({ demo }: Props) => {
  const { user } = useAuth()
  const [clips, setClips] = useState(false)

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
      <div className="flex gap-4">
        <div className="w-16 sm:w-32 lg:w-64 aspect-square sm:aspect-video shrink-0 rounded-md overflow-hidden">
          <DemoThumbnail demo={demo} />
        </div>
        <div className="flex justify-between w-full">
          <div className="flex flex-col gap-2 justify-center">
            <div className="flex items-center gap-4">
              <p className={`text-lg sm:text-2xl font-bold uppercase ${resultColor[player.stat.result]}`}>{resultString[player.stat.result]}</p>
              <p className="sm:text-xl text-white">{score()}</p>
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
                <p className="text-secondary">{formatDate(demo.createdAt)}</p>
              </div>
              <div className="flex items-center gap-4">
                <p className="text-secondary">In group</p>
                {highlights.length > 0 && (
                  <Button variant="subtle" color="muted" onClick={() => setClips(prev => !prev)} rightSection={<LuChevronDown className={`transform duration-300 ${clips ? "rotate-180" : ""}`} />}>
                    {`${clips ? "Hide" : "Show"} clips`}
                  </Button>
                )}
              </div>
            </div>
          )}
        </div>
      </div>
      <Collapse in={clips}>
        <div className="pt-8">
          <HighlightCarousel highlights={highlights} />
        </div>
      </Collapse>
    </div>
  )
}

const Loading = () => {
  return (
    <div className="flex gap-4">
      <div className="w-64 aspect-video shrink-0 rounded-md overflow-hidden">
        <LoadableImage />
      </div>
      <div className="flex flex-col gap-2 justify-center">
        <p className="text-2xl font-bold text-secondary">Processing match...</p>
      </div>
    </div>
  )
}

const Failed = () => {
  return (
    <div className="flex gap-4">
      <div className="w-64 aspect-video shrink-0 rounded-md overflow-hidden">
        <Center className="h-full">
          <LuTriangleAlert className="size-12 text-red-400" />
        </Center>
      </div>
      <div className="flex flex-col gap-2 justify-center">
        <p className="text-2xl font-bold text-red-400">Highlight generation failed</p>
        <p className="text-secondary">{`We couldn't process this match`}</p>
      </div>
    </div>
  )
}

const ClipBadge = ({ demo, highlights }: { demo: DemoType, highlights: Highlight[] }) => {
  if (demo.status === DemoStatus.Finished && highlights.length === 0) {
    return null
  }

  let text: string;
  let icon: ReactNode;

  if (demo.status === DemoStatus.Finished) {
    text = `${highlights.length} Clip${highlights.length !== 1 ? 's' : ''} generated`
    icon = <LuClapperboard className="text-(--mantine-color-primary-6)" />
  } else {
    text = "Clips are rendering"
    icon = <FragtapeIcon animated className="size-5 text-(--mantine-color-primary-6)" />
  }

  return (
    <div className="flex items-center gap-2 bg-(--mantine-color-primary-light) rounded-lg py-2 px-4 w-fit">
      {icon}
      <p className="text-white text-sm">{text}</p>
    </div>
  )
}
