import { API } from "./api";

export interface User {
  id: number;
  uid: number;
  name: string;
  displayName: string;
  avatarUrl: string;
  admin: boolean;
}

export const convertUser = (user: API.User): User => {
  return {
    id: user.id,
    uid: user.uid,
    name: user.name,
    displayName: user.display_name,
    avatarUrl: user.avatar_url,
    admin: user.admin,
  }
}

export const convertUsers = (users: API.User[]): User[] => {
  return users.map(convertUser)
}

export interface UserFilterResult {
  users: User[];
  total: number;
}

export const convertUserFilterResult = (r: API.UserFilterResult): UserFilterResult => {
  return {
    users: convertUsers(r.users),
    total: r.total,
  }
}

export interface UserFilter {
  name?: string;
  admin?: boolean;
  real?: boolean;
}
