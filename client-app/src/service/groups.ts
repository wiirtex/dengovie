import type {Result} from "../types/result.ts";
import api from "../api";
import type {AxiosError} from "axios";
import {transliterate} from "../utils.ts";

export interface GroupUser {
  ID: string,
  Name: string,
  Alias: string
}

function getTelegram(name: string) {
  const n =  name.split(" ")[0];

  return '@' + transliterate(n.toLowerCase());
}

export async function users(): Promise<Result<GroupUser[], 'not_found' | 'unknown'>> {
  try {
    const {data} = await api.get<GroupUser[]>("/groups/0/users");
    
    return { type: "success", data: data.map(x => ({...x, Alias: getTelegram(x.Name)} as GroupUser)) };
  } catch (error) {
    const e = error as AxiosError;

    if (e.status === 404) {
      return { type: "failure", error: 'not_found' };
    }

    return { type: "failure", error: 'unknown' };
  }
}