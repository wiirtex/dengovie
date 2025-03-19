import type {Result} from "../types/result.ts";
import api from "../api";
import type {AxiosError} from "axios";

export interface Me {
  user_id: string,
  user_name: string,
  user_alias: string
}

export async function me(): Promise<Result<Me, 'not_authorized' | 'unknown'>> {
  try {
    const {data} = await api.get<Me>("/user");

    return { type: "success", data };
  } catch (error) {
    const e = error as AxiosError;

    if (e.status === 403) {
      return { type: "failure", error: 'not_authorized' };
    }

    return { type: "failure", error: 'unknown' };
  }
}