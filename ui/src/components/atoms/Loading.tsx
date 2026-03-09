import { Center, CenterProps } from "@mantine/core"
import { FragtapeIcon } from "../icons/FragtapeIcon"
import { cn } from "@/lib/utils"

type Props = {
  isFetchingNextPage?: boolean;
  hasNextPage?: boolean;
} & CenterProps

export const Loading = ({ isFetchingNextPage, hasNextPage, className, ...props }: Props) => {
  const isFiltered = isFetchingNextPage !== undefined && hasNextPage !== undefined

  let filtereClassNames = ""
  if (isFiltered) filtereClassNames = isFetchingNextPage ? "flex" : hasNextPage ? "invisible" : "hidden"

  return (
    <Center className={cn(`mt-48 ${filtereClassNames}`, className)} {...props}>
      <FragtapeIcon animated className="size-36 text-(--mantine-color-primary-6)" />
    </Center>
  )
}
