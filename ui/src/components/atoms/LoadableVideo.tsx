import { cn } from "@/lib/utils";
import { ComponentProps, useEffect, useState } from "react";
import { FragtapeIcon } from "../icons/FragtapeIcon";

type Props = ComponentProps<"video">

export const LoadableVideo = ({ src, onLoadedData, onError, className, ...props }: Props) => {
  const [loaded, setLoaded] = useState(false)
  const [failed, setFailed] = useState(false)

  const handleLoaded = (e: React.SyntheticEvent<HTMLVideoElement, Event>) => {
    setLoaded(true)
    setFailed(false)
    onLoadedData?.(e)
  }

  const handleFailed = (e: React.SyntheticEvent<HTMLVideoElement, Event>) => {
    setLoaded(true)
    setFailed(true)
    onError?.(e)
  }

  useEffect(() => {
    setLoaded(false);
    setFailed(false);
  }, [src]);

  return (
    <div className="relative h-full w-full">
      {!loaded && !failed && (
        <div className="absolute inset-0 grid place-items-center">
          <FragtapeIcon animated className="size-12 text-(--mantine-color-primary-6)" />
        </div>
      )}

      <video
        preload="metadata"
        controls
        src={src}
        onLoadedData={handleLoaded}
        onError={handleFailed}
        className={cn("transition-opacity duration-300", loaded ? "opacity-100" : "opacity-0", className)}
        {...props}
      />
    </div>
  )
}
