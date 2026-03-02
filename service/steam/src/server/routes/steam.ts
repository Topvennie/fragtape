import type { Express } from "express";
import type { SteamService } from "../../steam/steam-service.js";
import { nextMatchReq } from "../../type/steam.js";

export function registerSteam(app: Express, steam: SteamService) {
  app.post("/steam/next-demo", async (req, res) => {
    const parsed = nextMatchReq.safeParse(req.body)
    if (!parsed.success) return res.status(400).json({ error: "invalid_body" })

    try {
      const out = await steam.nextDemoUrl(parsed.data);

      return res.status(200).json(out);
    } catch (e) {
      console.error(e)
      return res.status(500).json({ error: "internal_error", message: `${e}` });
    }
  });
}
