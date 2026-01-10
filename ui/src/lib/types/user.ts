import { API } from "./api";

export interface User {
  id: number;
  uid: string;
  name: string;
  displayName: string;
  avatarUrl: string;
}

export const convertUser = (user: API.User): User => {
  return {
    id: user.id,
    uid: user.uid,
    name: user.name,
    displayName: user.display_name,
    avatarUrl: user.avatar_url
  }
}
