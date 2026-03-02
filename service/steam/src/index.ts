import { initConfig } from "./config.js";
import { buildApp } from "./server/app.js";
import { SteamService } from "./steam/steam-service.js";

async function main() {
  const cfg = initConfig();

  const port = cfg.getNumber("service.steam.port", 8001);

  const accountName = cfg.getString("worker.steam.username", "");
  const password = cfg.getString("worker.steam.password", "");
  const sharedSecret = cfg.getString("worker.steam.shared_secret", "");

  const steam = new SteamService({ accountName, password, sharedSecret });

  steam.start();

  const app = buildApp(steam);
  app.listen(port, "0.0.0.0", () => {
    // eslint-disable-next-line no-console
    console.log(`steam service listening on :${port}`);
  });
}

main().catch((e) => {
  // eslint-disable-next-line no-console
  console.error(e);
  process.exit(1);
});
