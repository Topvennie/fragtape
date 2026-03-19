import z from "zod";
import { API } from "./api";
import { JSONBody } from "./general";

export interface SettingUser {
  connectedSteam: boolean;
  connectedFaceit: boolean;
  firstTimeWizard: boolean;
}

export const convertSettingUser = (s: API.SettingUser): SettingUser => {
  return {
    connectedSteam: s.connected_steam,
    connectedFaceit: s.connected_faceit,
    firstTimeWizard: s.first_time_wizard,
  }
}

const matchTokenRegex = /^CSGO-[A-Za-z0-9]{5}(?:-[A-Za-z0-9]{5}){4}$/
const authTokenRegex = /^[A-Za-z0-9]{4}-[A-Za-z0-9]{5}-[A-Za-z0-9]{4}$/

export const settingUserSteamSchema = z.object({
  matchToken: z
    .string()
    .trim()
    .min(1, "Match token is required")
    .regex(
      matchTokenRegex,
      "Invalid match token format. Expected: CSGO-xxxxx-xxxxx-xxxxx-xxxxx-xxxxx (letters/numbers)"
    ),
  authenticationToken: z
    .string()
    .trim()
    .min(1, "Authentication token is required")
    .regex(
      authTokenRegex,
      "Invalid authentication token format. Expected: xxxx-xxxxx-xxxx (letters/numbers)"
    ),
  importOld: z.boolean(),
})
export type SettingUserSteamSchema = z.infer<typeof settingUserSteamSchema> & JSONBody;
