import { Demo } from "@/lib/types/demo"
import { LoadableImage } from "../atoms/LoadableImage"
import { ComponentProps } from "react"

type Props = {
  demo: Demo
} & ComponentProps<"img">

const getSrc = (map: string) => {
  return `https://raw.githubusercontent.com/ghostcap-gaming/cs2-map-images/refs/heads/main/cs2/${map}.png`
}

export const DemoThumbnail = ({ demo, ...props }: Props) => {
  return <LoadableImage src={getSrc(demo.stats.map)} {...props} />
}
