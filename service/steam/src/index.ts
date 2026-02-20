import express from "express";
import { initConfig } from "./config.js";

async function main() {
  const cfg = initConfig();

  const port = cfg.getNumber("service.steam.port", 8001);

  const app = express();

  app.get("/health", (_req, res) => {
    res.status(200).send("ok");
  });

  app.listen(port, "0.0.0.0", () => {
    console.log(`steamgc listening on :${port} (env=${cfg.env})`);
  });
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
