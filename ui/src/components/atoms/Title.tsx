import { cn } from "@/lib/utils"
import { TitleProps, Title as MantineTitle } from "@mantine/core"

const orderClasses = [
  "",
  "font-black text-xl sm:text-2xl md:text-3xl",
  "font-extrabold text-lg sm:text-xl md:text-2xl",
  "font-extrabold text-[1rem] sm:text-lg md:text-xl",
  "font-bold text-[1rem] sm:text-md md:text-lg",
  "font-bold text-sm sm:text-md md:text-md",
  "font-semibold text-xs sm:text-sm md:text-sm",
]

export const Title = ({ children, className, ...props }: TitleProps) => {
  return (
    <MantineTitle c="white" className={cn("wrap-anywhere break-normal whitespace-pre-wrap", orderClasses[props.order ?? 1], className)} textWrap="wrap" {...props}>{children}</MantineTitle>
  )
}
