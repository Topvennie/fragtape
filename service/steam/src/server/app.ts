import express from "express";
import type { SteamService } from "../steam/steam-service.js";
import { registerSteam } from "./routes/steam.js";

export function buildApp(steam: SteamService) {
  const app = express();
  app.use(express.json());

  registerSteam(app, steam);

  return app;
}
