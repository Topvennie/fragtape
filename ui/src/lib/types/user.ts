import { API } from "./api";

export interface User {
  id: number;
  uid: string;
  name: string;
  email: string;
  hasProfile: boolean;
}

export const convertUser = (user: API.User): User => {
  let name = user.name
  if (user.display_name !== "") name = user.display_name

  return {
    id: user.id,
    uid: user.uid,
    name: name,
    email: user.email,
    hasProfile: false,
  }
}
