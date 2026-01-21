import { useAuth } from "@/lib/hooks/useAuth"
import { DemoStatus, Demo as DemoType } from "@/lib/types/demo"
import { Highlight } from "@/lib/types/highlight"
import { Result, resultString } from "@/lib/types/stat"
import { formatDate } from "@/lib/utils"
import { Button, Collapse } from "@mantine/core"
import { ReactNode, useMemo, useState } from "react"
import { LuChevronDown, LuClapperboard } from "react-icons/lu"
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
  const { user } = useAuth()

  const [clips, setClips] = useState(false)

  const highlights = useMemo(() => demo.players.flatMap(p => p.highlights), [demo])
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
    <Card>
      <div className="flex flex-col">
        <div className="flex gap-4">
          <div className="w-62 aspect-video rounded-md overflow-hidden">
            <DemoThumbnail demo={demo} />
          </div>
          <div className="flex justify-between w-full">
            <div className="flex flex-col gap-2 justify-center">
              <div className="flex items-center gap-4">
                <p className={`text-2xl font-bold uppercase ${resultColor[player.stat.result]}`}>{resultString[player.stat.result]}</p>
                <p className="text-xl text-white">{score()}</p>
              </div>
              <div className="space-x-4 text-secondary">
                K <span className="text-white">{player.stat.kills}</span>
                D <span className="text-white">{player.stat.deaths}</span>
                A <span className="text-white">{player.stat.assists}</span>
              </div>
              <ClipBadge demo={demo} highlights={highlights} />
            </div>
            <div className="flex flex-col items-end justify-between">
              <div>
                <p className="text-secondary">{formatDate(demo.createdAt)}</p>
              </div>
              <div className="flex items-center gap-4">
                <p className="text-secondary">In group</p>
                <Button variant="subtle" color="muted" onClick={() => setClips(prev => !prev)} rightSection={<LuChevronDown className={`transform duration-300 ${clips ? "rotate-180" : ""}`} />}>
                  {`${clips ? "Hide" : "Show"} clips`}
                </Button>
              </div>
            </div>
          </div>
        </div>
        <Collapse in={clips}>
          <div className="pt-8">
            <HighlightCarousel highlights={highlights} />
          </div>
        </Collapse>
      </div>
    </Card>
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
