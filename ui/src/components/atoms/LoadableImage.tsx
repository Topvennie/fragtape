import { cn } from "@/lib/utils";
import { ComponentProps, useState } from "react";
import { FragtapeIcon } from "../icons/FragtapeIcon";

type Props = ComponentProps<"img">

export const LoadableImage = ({ className, ...props }: Props) => {
  const [loaded, setLoaded] = useState(false)

  return (
    <div className="relative h-full w-full">
      {!loaded && (
        <div className="absolute inset-0 grid place-items-center">
          <FragtapeIcon animated className="size-12 text-(--mantine-color-primary-6)" />
        </div>
      )}

      <img
        onLoad={() => setLoaded(true)}
        loading="lazy"
        className={cn(`h-full w-full object-cover transition-opacity duration-300 ${loaded ? "opacity-100" : "opacity-0"}`, className)}
        {...props}
      />
    </div>
  )
}

