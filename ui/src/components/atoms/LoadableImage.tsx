import { cn } from "@/lib/utils";
import { Skeleton } from "@mantine/core";
import { ComponentProps, useState } from "react";

type Props = ComponentProps<"img">

export const LoadableImage = ({ className, ...props }: Props) => {
  const [loaded, setLoaded] = useState(false)

  return (
    <div className="relative">
      {!loaded && (
        <div className="absolute inset-0 flex items-center justify-center">
          <Skeleton animate radius="md" className="w-full h-full" />
        </div>
      )}

      <img
        onLoad={() => setLoaded(true)}
        className={cn("w-full h-full object-cover transition-opacity duration-300", loaded ? "opacity-100" : "opacity-0", className)}
        {...props}
      />
    </div>
  )
}
