export type Result<Data, Error> =
  | { type: "success"; data: Data }
  | { type: "failure"; error: Error };