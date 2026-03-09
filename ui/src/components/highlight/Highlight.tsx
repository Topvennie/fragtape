import { Highlight as HighlightType } from "@/lib/types/highlight";
import { formatDurationS } from "@/lib/utils";
import { LuClock } from "react-icons/lu";
import { LoadableVideo } from "../atoms/LoadableVideo";
import { FragtapeIcon } from "../icons/FragtapeIcon";

type Props = {
  highlight: HighlightType;
}

export const Highlight = ({ highlight }: Props) => {
  if (!highlight.generated) return <Rendering highlight={highlight} />

  return <Rendered highlight={highlight} />
}

const Rendering = ({ highlight }: { highlight: HighlightType }) => {
  return (
    <div className="overflow-hidden rounded-lg">
      <div className="relative aspect-video">
        <div className="absolute inset-0 grid place-items-center">
          <FragtapeIcon animated className="size-12 text-(--mantine-color-primary-6)" />
        </div>
      </div>
      <div className="flex flex-col bg-(--mantine-color-background-light) px-4 py-2 lg:gap-2 lg:py-4">
        <p className="text-md font-bold text-white lg:text-lg">{highlight.title}</p>
        <div className="flex items-center gap-2">
          <LuClock className="size-4 text-secondary" />
          <p className="text-sm text-secondary lg:text-md">
            {`Round ${highlight.round} · ${formatDurationS(highlight.durationS)}`}
          </p>
        </div>
      </div>
    </div>
  );
}

const Rendered = ({ highlight }: { highlight: HighlightType }) => {
  return (
    <div className="rounded-lg overflow-hidden">
      <div className="aspect-video">
        <LoadableVideo src={`/api/highlight/video/${highlight.id}`} />
      </div>
      <div className="flex flex-col lg:gap-2 px-4 py-2 lg:py-4 bg-(--mantine-color-background-light)">
        <p className="text-white text-md lg:text-lg font-bold">{highlight.title}</p>
        <div className="flex items-center gap-2">
          <LuClock className="text-secondary size-4" />
          <p className="text-secondary text-sm lg:text-md">{`Round ${highlight.round} · ${formatDurationS(highlight.durationS)}`}</p>
        </div>
      </div>
    </div>
  )
}
