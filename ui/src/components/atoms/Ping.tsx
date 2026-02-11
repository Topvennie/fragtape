type Props = {
  healthy: boolean
}

export const Ping = ({ healthy }: Props) => {
  const bg = healthy ? "bg-green-500" : "bg-red-500"

  return (
    <div className="relative">
      <div className={`w-2 h-2 rounded-full ${bg}`} />
      <div className={`w-2 h-2 rounded-full absolute inset-0 animate-slow-ping ${bg}`} />
    </div>
  )
}
