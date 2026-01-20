import { Highlight as HighlightType } from "@/lib/types/highlight"
import { LoadableVideo } from "../atoms/LoadableVideo";
import { LuClock } from "react-icons/lu";
import { formatDurationS } from "@/lib/utils";

type Props = {
  highlight: HighlightType;
}

export const Highlight = ({ highlight }: Props) => {
  return (
    <div className="rounded-lg overflow-hidden">
      <div className="aspect-video">
        <LoadableVideo src={`/api/highlight/video/${highlight.id}`} />
      </div>
      <div className="flex flex-col gap-2 p-4 bg-(--mantine-color-background-light)">
        <p className="text-white text-xl font-bold">{highlight.title}</p>
        <div className="flex items-center gap-2">
          <LuClock className="text-secondary size-4" />
          <p className="text-secondary text-md">{`Round ${highlight.round} · ${formatDurationS(highlight.durationS)}`}</p>
        </div>
      </div>
    </div>
  )
}
