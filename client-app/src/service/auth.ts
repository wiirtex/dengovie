import api from "../api";
import type {Result} from "../types/result.ts";
import type {AxiosError} from "axios";

interface Login {
  alias: string,
  password: string,
}

type LoginError = 'invalid_password' | 'server_error' | 'unknown';

export async function login(p: Login) : Promise<Result<null, LoginError>> {
  try {
    await api.post("/auth/login", {
      alias: p.alias,
      otp: p.password,
    });

    return { type: "success", data: null };
  } catch (error) {
    const e = error as AxiosError;

    console.log({e});

    if (e.status === 500) {
      return { type: "failure", error: 'server_error' };
    } else if (e.status === 403) {
      const errorData = e.response?.data as { ErrorReason: 'invalid_otp' } | undefined;

      if (errorData?.ErrorReason === 'invalid_otp') {
        return { type: "failure", error: 'invalid_password' };
      }
    }

    return { type: "failure", error: 'unknown' };
  }
}

export async function logout() : Promise<Result<null, 'server_error'>> {
  try {
    await api.post('/auth/logout');

    return { type: "success", data: null };
  } catch (error) {
    return { type: "failure", error: 'server_error' };
  }
}
