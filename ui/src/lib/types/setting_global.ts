import { API } from "./api";
import { JSONBody } from "./general";

import z from "zod";

export interface SettingGlobal {
  demoUpload: boolean;
  customCriteria: boolean;
  chatCommand: boolean;
  chatTrigger: string;
}

export const convertSettingGlobal = (s: API.SettingGlobal): SettingGlobal => {
  return {
    demoUpload: s.demo_upload,
    customCriteria: s.custom_criteria,
    chatCommand: s.chat_command,
    chatTrigger: s.chat_trigger,
  }
}

export const convertSettingGlobalSchema = (s: SettingGlobal): SettingGlobalSchema => {
  return s as SettingGlobalSchema
}

export const settingGlobalSchema = z.object({
  demoUpload: z.boolean(),
  customCriteria: z.boolean(),
  chatCommand: z.boolean(),
  chatTrigger: z.string().min(3),
})
export type SettingGlobalSchema = z.infer<typeof settingGlobalSchema> & JSONBody;
