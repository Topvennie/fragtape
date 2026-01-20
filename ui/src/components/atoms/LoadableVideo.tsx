import { cn } from "@/lib/utils";
import { Skeleton } from "@mantine/core";
import { ComponentProps, useState } from "react";

type Props = ComponentProps<"video">

export const LoadableVideo = ({ className, ...props }: Props) => {
  const [loaded, setLoaded] = useState(false)

  return (
    <div className="relative h-full">
      {!loaded && (
        <div className="absolute inset-0 flex items-center justify-center">
          <Skeleton animate className="w-full h-full" />
        </div>
      )}

      <video
        controls
        onLoadedData={() => setLoaded(true)}
        className={cn("transition-opacity duration-300", loaded ? "opacity-100" : "opacity-0", className)}
        {...props}
      />
    </div>
  )
}
