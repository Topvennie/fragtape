import { Center, CenterProps } from "@mantine/core"
import { FragtapeIcon } from "../icons/FragtapeIcon"
import { cn } from "@/lib/utils"

type Props = {
  isFetchingNextPage?: boolean;
  hasNextPage?: boolean;
  size?: string;
} & CenterProps

export const Loading = ({ isFetchingNextPage, hasNextPage, size, className, ...props }: Props) => {
  const isFiltered = isFetchingNextPage !== undefined && hasNextPage !== undefined

  let filtereClassNames = ""
  if (isFiltered) filtereClassNames = isFetchingNextPage ? "flex" : hasNextPage ? "invisible" : "hidden"

  if (!size) size = "size-36"

  return (
    <Center className={cn(`mt-48 ${filtereClassNames}`, className)} {...props}>
      <FragtapeIcon animated className={`text-(--mantine-color-primary-6) ${size}`} />
    </Center>
  )
}
