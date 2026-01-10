import { Demo as DemoType } from "@/lib/types/demo"

type Props = {
  demo: DemoType
}

export const Demo = ({ demo }: Props) => {
  return <p>{demo.id}</p>
}
