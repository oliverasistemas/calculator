import { useState, useCallback } from "react";
import { api, ApiError } from "@/api/client";
import type {
  Operation,
  CalculateRequest,
  CalculateResult,
  HistoryEntry,
} from "@/types/calculator";

interface CalculatorState {
  display: string;
  expression: string;
  history: HistoryEntry[];
  error: string | null;
  loading: boolean;
}

const INITIAL_STATE: CalculatorState = {
  display: "0",
  expression: "",
  history: [],
  error: null,
  loading: false,
};

const UNARY_OPS = new Set<Operation>(["sqrt"]);

export function useCalculator() {
  const [state, setState] = useState<CalculatorState>(INITIAL_STATE);
  // Operands are kept and sent as strings — parseFloat would corrupt
  // results beyond 2^53 when chaining (see CalculateRequest).
  const [firstOperand, setFirstOperand] = useState<string | null>(null);
  const [pendingOp, setPendingOp] = useState<Operation | null>(null);
  const [resetOnNext, setResetOnNext] = useState(false);

  const appendDigit = useCallback(
    (digit: string) => {
      setState((prev) => {
        const current = resetOnNext || prev.display === "0" ? "" : prev.display;
        if (digit === "." && current.includes(".")) return prev;
        const next = current + digit || "0";
        return { ...prev, display: next, error: null };
      });
      if (resetOnNext) setResetOnNext(false);
    },
    [resetOnNext]
  );

  const setOperation = useCallback(
    (op: Operation) => {
      if (UNARY_OPS.has(op)) {
        performCalculation(op, state.display);
        return;
      }

      const current = state.display;

      if (firstOperand !== null && pendingOp && !resetOnNext) {
        performCalculation(pendingOp, firstOperand, current, op);
      } else {
        setFirstOperand(current);
        setPendingOp(op);
        setResetOnNext(true);
      }
    },
    [state.display, firstOperand, pendingOp, resetOnNext]
  );

  const calculate = useCallback(() => {
    if (firstOperand === null || !pendingOp) return;
    performCalculation(pendingOp, firstOperand, state.display);
  }, [firstOperand, pendingOp, state.display]);

  const performCalculation = async (
    op: Operation,
    a: string,
    b?: string,
    nextOp?: Operation
  ) => {
    setState((prev) => ({ ...prev, loading: true, error: null }));

    const body: CalculateRequest = { operation: op, a };
    if (b !== undefined) body.b = b;

    try {
      const result = await api.post<CalculateResult>("/calculate", body);
      setState((prev) => ({
        ...prev,
        display: result.resultText,
        expression: result.expression,
        history: [
          { expression: result.expression, result: result.resultText },
          ...prev.history,
        ].slice(0, 50),
        loading: false,
      }));

      if (nextOp) {
        setFirstOperand(result.resultText);
        setPendingOp(nextOp);
        setResetOnNext(true);
      } else {
        setFirstOperand(null);
        setPendingOp(null);
        setResetOnNext(true);
      }
    } catch (err) {
      const message =
        err instanceof ApiError ? err.message : "Calculation failed";
      setState((prev) => ({ ...prev, loading: false, error: message }));
    }
  };

  const clear = useCallback(() => {
    setState((prev) => ({ ...INITIAL_STATE, history: prev.history }));
    setFirstOperand(null);
    setPendingOp(null);
    setResetOnNext(false);
  }, []);

  const clearAll = useCallback(() => {
    setState(INITIAL_STATE);
    setFirstOperand(null);
    setPendingOp(null);
    setResetOnNext(false);
  }, []);

  const toggleSign = useCallback(() => {
    setState((prev) => {
      if (prev.display === "0") return prev;
      const toggled = prev.display.startsWith("-")
        ? prev.display.slice(1)
        : "-" + prev.display;
      return { ...prev, display: toggled };
    });
  }, []);

  const backspace = useCallback(() => {
    setState((prev) => {
      if (prev.display.length <= 1 || (prev.display.length === 2 && prev.display.startsWith("-"))) {
        return { ...prev, display: "0" };
      }
      return { ...prev, display: prev.display.slice(0, -1) };
    });
  }, []);

  return {
    display: state.display,
    expression: state.expression,
    history: state.history,
    error: state.error,
    loading: state.loading,
    pendingOp,
    appendDigit,
    setOperation,
    calculate,
    clear,
    clearAll,
    toggleSign,
    backspace,
  };
}
