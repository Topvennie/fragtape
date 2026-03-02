import GlobalOffensive from "globaloffensive";
import SteamTotp from "steam-totp";
import SteamUser from "steam-user";
import { SteamUserCfg } from "../type/steam";

export class SteamSession {
  private user = new SteamUser();
  private csgo = new GlobalOffensive(this.user)

  private cfg: SteamUserCfg;

  private loggedOnResolve!: () => void;
  private connectedGcResolve!: () => void;

  private loggedOn = new Promise<void>((r) => (this.loggedOnResolve = r));
  private connectedGc = new Promise<void>((r) => this.connectedGcResolve = r);

  constructor(cfg: SteamUserCfg) {
    this.cfg = cfg;
  }

  start() {
    this.user.logOn({
      accountName: this.cfg.accountName,
      password: this.cfg.password,
      twoFactorCode: SteamTotp.getAuthCode(this.cfg.sharedSecret),
    })

    this.user.on("steamGuard", (_domain: string | null, callback: (code: string) => void) => {
      const code = SteamTotp.getAuthCode(this.cfg.sharedSecret);
      callback(code);
    });

    this.user.on("loggedOn", async () => {
      console.debug(`Logged into steam as ${this.user.steamID?.getSteamID64()}`)
      this.user.setPersona(SteamUser.EPersonaState.Online);
      this.user.gamesPlayed([730])

      this.loggedOnResolve()
    })

    this.csgo.on("connectedToGC", () => {
      console.debug("Connected to gc")
      this.connectedGcResolve()
    })
  }

  private async waitReady(): Promise<void> {
    await this.loggedOn
    await this.connectedGc
  }

  async requestMatchList(shareCode: string, timeoutMs = 15_000): Promise<GlobalOffensive.MatchesData> {
    await this.waitReady();

    return await new Promise((resolve, reject) => {
      const t = setTimeout(() => {
        cleanup();
        reject(new Error("gc timeout reached"));
      }, timeoutMs);

      const onMatchList = (_: GlobalOffensive.Match[], data: GlobalOffensive.MatchesData) => {
        cleanup();
        resolve(data);
      };

      const cleanup = () => {
        clearTimeout(t);
        this.csgo.removeListener("matchList", onMatchList);
      };

      this.csgo.on("matchList", onMatchList);
      this.csgo.requestGame(shareCode);
    });
  }
}

