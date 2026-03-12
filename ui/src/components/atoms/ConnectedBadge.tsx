import { Ping } from "./Ping"

export const ConnectedBadge = () => {
  return (
    <div className="rounded-full bg-(--mantine-color-green-light) px-3 py-1 flex items-center gap-1 text-green-500 text-xs">
      <Ping healthy />
      Connected
    </div>
  )
}

