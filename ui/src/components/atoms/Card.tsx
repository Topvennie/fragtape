import { CardProps, Card as MCard } from "@mantine/core"

type Props = CardProps

export const Card = (props: Props) => {
  return <MCard shadow="lg" padding="xl" bg="background.9" {...props} />
}
