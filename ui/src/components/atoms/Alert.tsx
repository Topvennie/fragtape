import { cn } from "@/lib/utils";
import { ComponentProps, ReactNode } from "react";

type Props = {
  icon?: ReactNode;
  title?: string;
  border?: boolean;
  color: Color;
} & ComponentProps<"div">

type Color = "orange" | "blue" | "green" | "red"

const backgroundColor: Record<Color, string> = {
  "orange": "bg-orange-400/10",
  "blue": "bg-blue-400/10",
  "green": "bg-green-400/10",
  "red": "bg-red-400/10",
}

const borderColor: Record<Color, string> = {
  "orange": "border-orange-400/40",
  "blue": "border-blue-400/40",
  "green": "border-green-400/40",
  "red": "bg-border-400/40",
}

export const Alert = ({ color, icon = null, title = "", border = false, children, className, ...props }: Props) => {
  return (
    <div className={cn(`flex items-start gap-4 p-4 rounded-md ${backgroundColor[color]} ${border ? "border " + borderColor[color] : ""}`, className)} {...props}>
      {icon}
      <div className="flex flex-col gap-2">
        {title && <div className="font-bold text-white">{title}</div>}
        <div className="text-secondary whitespace-pre-wrap">{children}</div>
      </div>
    </div>
  )
}
