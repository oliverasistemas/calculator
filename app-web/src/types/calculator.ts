export type Operation =
  | "add"
  | "subtract"
  | "multiply"
  | "divide"
  | "power"
  | "sqrt"
  | "percentage";

export interface CalculateRequest {
  operation: Operation;
  a: number;
  b?: number;
}

export interface CalculateResult {
  result: number;
  expression: string;
}

export interface HistoryEntry {
  expression: string;
  result: number;
}
