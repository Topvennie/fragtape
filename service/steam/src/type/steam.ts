import { z } from "zod";

export const nextMatchReq = z.object({
  webApiKey: z.string().min(1),
  steamId: z.number(),
  authToken: z.string().min(1),
  matchToken: z.string().min(1),
})
export type NextMatchReq = z.infer<typeof nextMatchReq>

export type NextMatchResp = {
  nextCode?: string;
  demoUrl?: string;
  code?: number;
  error?: string;
}

export type SteamUserCfg = {
  accountName: string;
  password: string;
  sharedSecret: string;
}
