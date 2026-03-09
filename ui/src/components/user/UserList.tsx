import { User as UserType } from "@/lib/types/user"
import { ReactNode } from "react";
import { UserIcon } from "./UserIcon";
import { FragtapeIcon } from "../icons/FragtapeIcon";

type Props = {
  users: UserType[];
  rightSection?: (user: UserType) => ReactNode;
  isLoading?: boolean;
}

export const UserList = ({ users, rightSection = () => null, isLoading = false }: Props) => {
  if (!users.length) {
    return <p className="text-secondary text-center">No users</p>
  }

  return (
    <div className="divide-y divide-gray-800">
      {isLoading
        ? <FragtapeIcon animated className="size-12" />
        : users.map(u => <User key={u.id} user={u} right={rightSection(u)} />)
      }
    </div>
  )
}

const User = ({ user, right }: { user: UserType, right: ReactNode }) => {
  return (
    <div className="flex items-center justify-between py-2">
      <div className="flex items-center gap-2">
        <UserIcon user={user} />
        <div>
          <p className="text-white">{user.name}</p>
          <p className="text-secondary text-sm/3">{user.displayName}</p>
        </div>
      </div>
      {right}
    </div>
  )
}
