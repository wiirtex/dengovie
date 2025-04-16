import type {Result} from "../types/result.ts";
import api from "../api";
import type {AxiosError} from "axios";

export interface DebtItem {
  "another_user_id": number,
  "another_user_name": string,
  "amount": number
}

export interface DebtsAll {
  Debts: DebtItem[]
}

export async function all(): Promise<Result<DebtsAll, 'not_found' | 'unknown'>> {
  try {
    const {data} = await api.get<DebtsAll>("/debts");

    if (data.Debts === null) {
      data.Debts = [];
    }

    return { type: "success", data: data};
  } catch (error) {
    const e = error as AxiosError;

    if (e.status === 404) {
      return { type: "failure", error: 'not_found' };
    }

    return { type: "failure", error: 'unknown' };
  }
}

export async function share(amount: number, user_ids: number[]): Promise<Result<undefined, 'unknown'>> {
  try {
    await api.post("/debts/share", { group_id: 0, amount, user_ids });

    return { type: "success", data: undefined };
  } catch (error) {
    // const e = error as AxiosError;

    return { type: "failure", error: 'unknown' };
  }
}

export async function pay(amount: number, user_id: number): Promise<Result<undefined, 'unknown'>> {
  try {
    await api.post("/debts/pay", { full: true, amount, another_user_id: user_id });

    return { type: "success", data: undefined };
  } catch (error) {
    // const e = error as AxiosError;

    return { type: "failure", error: 'unknown' };
  }
}