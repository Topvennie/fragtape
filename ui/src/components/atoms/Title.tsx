import { cn } from "@/lib/utils"
import { TitleProps, Title as MantineTitle } from "@mantine/core"

const orderClasses = [
  "",
  "text-2xl sm:text-3xl md:text-4xl",
  "text-xl sm:text-2xl md:text-3xl",
  "text-[1rem] sm:text-xl md:text-2xl",
  "text-[1rem] sm:text-lg md:text-xl",
  "text-sm sm:text-md md:text-lg",
  "text-xs sm:text-sm md:text-md",
]

export const Title = ({ children, className, ...props }: TitleProps) => {
  return (
    <MantineTitle c="white" className={cn("font-black wrap-anywhere break-normal whitespace-pre-wrap", orderClasses[props.order ?? 1], className)} textWrap="wrap" {...props}>{children}</MantineTitle>
  )
}
