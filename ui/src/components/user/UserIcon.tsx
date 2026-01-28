import { User } from "@/lib/types/user"
import { Avatar } from "@mantine/core";

type Props = {
  user: User;
}

export const UserIcon = ({ user }: Props) => {
  return <Avatar src={user.avatarUrl} name={user.name} alt={user.name} className="border-2 border-(--mantine-color-primary-6)" />
}
