import { NextMatchReq, NextMatchResp } from "../type/steam";

type SteamResp = {
  result?: {
    nextcode?: string;
  };
}

export async function getNextMatchSharingCode(req: NextMatchReq): Promise<NextMatchResp> {
  const u = new URL("https://api.steampowered.com/ICSGOPlayers_730/GetNextMatchSharingCode/v1")
  u.searchParams.set("key", req.webApiKey)
  u.searchParams.set("steamid", req.steamId.toString())
  u.searchParams.set("steamidkey", req.authToken)
  u.searchParams.set("knowncode", req.matchToken)

  const resp: NextMatchResp = {}
  let res: Response

  try {
    res = await fetch(u, { method: "GET" })
  } catch (e) {
    resp.error = `${e}`
    return resp
  }

  resp.code = res.status

  if (!res.ok) {
    resp.error = res.statusText
    return resp
  }

  const json = (await res.json()) as SteamResp
  resp.nextCode = json.result?.nextcode

  return resp
}
