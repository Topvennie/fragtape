import GlobalOffensive from "globaloffensive";
import { NextMatchReq, NextMatchResp } from "../type/steam.js";
import { getNextMatchSharingCode } from "./match-history.js";
import { SteamSession } from "./steam-session.js";

export class SteamService {
  private queue: Promise<any> = Promise.resolve();
  private session: SteamSession;

  constructor(
    cfg: {
      accountName: string,
      password: string,
      sharedSecret: string
    },
  ) {
    this.session = new SteamSession(cfg)
  }

  start(): void {
    this.session.start();
  }

  async nextDemoUrl(req: NextMatchReq): Promise<NextMatchResp> {
    this.queue = this.queue.then(async () => {
      const resp = await getNextMatchSharingCode(req);

      if (resp.error) return resp;

      // No new match
      if (!resp.nextCode || resp.code == 202) return resp;

      // Get new match info
      let matchData: GlobalOffensive.MatchesData
      try {
        matchData = await this.session.requestMatchList(resp.nextCode);
      } catch (e) {
        resp.error = `${e}`
        return resp
      }

      resp.demoUrl = this.extractDemoUrl(matchData);

      return resp;
    });

    return await this.queue;
  }

  private extractDemoUrl(matchData: GlobalOffensive.MatchesData) {
    for (let match of matchData.matches) {
      for (let round of match.roundstatsall) {
        if (round.map) return round.map
      }
    }

    return ""
  }
}
