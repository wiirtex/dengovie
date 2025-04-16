import type {Result} from "../types/result.ts";
import api from "../api";
import type {AxiosError} from "axios";

export interface Me {
  user_id: string,
  name: string,
  alias: string
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

export async function updName(newName: string): Promise<Result<undefined, 'not_authorized' | 'unknown'>> {
  try {
    await api.post("/user/update_name", {new_name: newName});

    return { type: "success", data: undefined };
  } catch (error) {
    const e = error as AxiosError;

    if (e.status === 403) {
      return { type: "failure", error: 'not_authorized' };
    }

    return { type: "failure", error: 'unknown' };
  }
}

export async function deleteAccount(): Promise<Result<undefined, 'not_authorized' | 'unknown'>> {
  try {
    await api.delete("/user/delete");

    return { type: "success", data: undefined };
  } catch (error) {
    const e = error as AxiosError;

    if (e.status === 403) {
      return { type: "failure", error: 'not_authorized' };
    }

    return { type: "failure", error: 'unknown' };
  }
}