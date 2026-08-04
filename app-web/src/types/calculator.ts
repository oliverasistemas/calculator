export type Operation =
  | "add"
  | "subtract"
  | "multiply"
  | "divide"
  | "power"
  | "sqrt"
  | "percentage";

// Operands and results travel as strings end to end: a round-trip through a
// JS number would crush integers beyond 2^53 (e.g. 3^35 = 50031545098999707)
// to the nearest float64.
export interface CalculateRequest {
  operation: Operation;
  a: string;
  b?: string;
}

export interface CalculateResult {
  // Nearest float64 only; resultText carries the authoritative digits.
  result: number;
  resultText: string;
  expression: string;
}

export interface HistoryEntry {
  expression: string;
  result: string;
}
