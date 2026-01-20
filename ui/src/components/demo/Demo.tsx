import { useAuth } from "@/lib/hooks/useAuth"
import { Demo as DemoType } from "@/lib/types/demo"
import { Result, resultString } from "@/lib/types/stat"
import { formatDate } from "@/lib/utils"
import { Carousel } from '@mantine/carousel'
import { Button, Collapse } from "@mantine/core"
import { useMemo, useState } from "react"
import { LuChevronDown, LuClapperboard } from "react-icons/lu"
import { Card } from "../atoms/Card"
import { Highlight } from "../highlight/Highlight"
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
              {player.highlights.length > 0 && (
                <div className="flex items-center gap-2 bg-(--mantine-color-primary-light) rounded-lg py-2 px-4 w-fit">
                  <LuClapperboard className="text-(--mantine-color-primary-6)" />
                  <p className="text-white text-sm">{`${player.highlights.length} Clip${player.highlights.length !== 1 ? 's' : ''} generated`}</p>
                </div>
              )}
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
            <DemoHighlights demo={demo} />
          </div>
        </Collapse>
      </div>
    </Card>
  )
}

const DemoHighlights = ({ demo }: { demo: DemoType }) => {
  const highlights = useMemo(() => demo.players.flatMap(p => p.highlights), [demo])

  if (highlights.length === 0) return null

  return (
    <Carousel
      slideSize={{ base: "100%", sm: "50%", lg: "33%" }}
      slideGap={{ base: "xl", sm: "lg" }}
      emblaOptions={{ align: "start", slidesToScroll: 1 }}
      styles={{
        controls: {
          left: 0,
          right: 0,
          padding: 0,
        }
      }}
      className="py-6 px-12"
    >
      {highlights.map(h => (
        <Carousel.Slide key={h.id}>
          <Highlight highlight={h} />
        </Carousel.Slide>
      ))}
    </Carousel>
  )
}
