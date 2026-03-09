import { Center } from "@mantine/core"
import { FragtapeIcon } from "../icons/FragtapeIcon"

export const Loading = () => {
  return (
    <Center className="mt-48">
      <FragtapeIcon animated className="size-36 text-(--mantine-color-primary-6)" />
    </Center>
  )
}
