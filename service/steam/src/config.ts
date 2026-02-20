import fs from "node:fs";
import path from "node:path";
import yaml from "js-yaml";
import dotenv from "dotenv";

type AnyObj = Record<string, any>;

function envKeyFor(dotKey: string): string {
  return dotKey.replace(/\./g, "_").toUpperCase();
}

function getFromObject(obj: AnyObj, dotKey: string): unknown {
  const parts = dotKey.split(".");

  let cur: any = obj;
  for (const p of parts) {
    if (cur == null || typeof cur !== "object") return undefined;
    cur = cur[p];
  }

  return cur;
}

function readYamlConfig(envName: string): AnyObj {
  const file = path.resolve(process.cwd(), "config", `${envName}.yml`);
  const content = fs.readFileSync(file, "utf8");
  const parsed = yaml.load(content);
  if (!parsed || typeof parsed !== "object") return {};
  return parsed as AnyObj;
}

export function initConfig() {
  const res = dotenv.config();
  if (res.error && (res.error as any).code !== "ENOENT") {
    console.error("Failed to load .env file", res.error);
  }

  const env = process.env[envKeyFor("app.env")] ?? "development";

  const cfg = readYamlConfig(env);

  function getString(key: string, def: string): string {
    const ek = envKeyFor(key);
    const v = process.env[ek];

    if (v != null && v !== "") return v;

    const fromYaml = getFromObject(cfg, key);
    if (typeof fromYaml === "string" && fromYaml !== "") return fromYaml;

    return def;
  }

  function getNumber(key: string, def: number): number {
    const ek = envKeyFor(key);
    const v = process.env[ek];

    if (v != null && v !== "") {
      const n = Number(v);
      if (!Number.isNaN(n)) return n;
    }

    const fromYaml = getFromObject(cfg, key);
    if (typeof fromYaml === "number") return fromYaml;
    if (typeof fromYaml === "string") {
      const n = Number(fromYaml);
      if (!Number.isNaN(n)) return n;
    }

    return def;
  }

  return {
    env,
    getString,
    getNumber
  };
}
