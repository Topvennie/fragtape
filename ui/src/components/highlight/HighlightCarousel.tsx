import { Highlight as HighlightType } from "@/lib/types/highlight"
import { Carousel } from "@mantine/carousel";
import { Highlight } from "./Highlight";

type Props = {
  highlights: HighlightType[];
}

export const HighlightCarousel = ({ highlights }: Props) => {
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
