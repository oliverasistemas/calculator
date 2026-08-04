import { renderHook, act } from "@testing-library/react";
import { describe, it, expect, vi, beforeEach } from "vitest";
import { useCalculator } from "./useCalculator";
import { api, ApiError } from "@/api/client";

vi.mock("@/api/client", async (importOriginal) => {
  const actual = await importOriginal<typeof import("@/api/client")>();
  return { ...actual, api: { post: vi.fn() } };
});

const mockPost = vi.mocked(api.post);

function typeDigits(
  result: { current: ReturnType<typeof useCalculator> },
  digits: string
) {
  for (const d of digits) {
    act(() => result.current.appendDigit(d));
  }
}

describe("useCalculator", () => {
  beforeEach(() => {
    mockPost.mockReset();
  });

  it("starts with an empty state", () => {
    const { result } = renderHook(() => useCalculator());
    expect(result.current.display).toBe("0");
    expect(result.current.expression).toBe("");
    expect(result.current.history).toEqual([]);
    expect(result.current.error).toBeNull();
    expect(result.current.loading).toBe(false);
    expect(result.current.pendingOp).toBeNull();
  });

  describe("appendDigit", () => {
    it("replaces the leading zero and accumulates digits", () => {
      const { result } = renderHook(() => useCalculator());
      typeDigits(result, "42");
      expect(result.current.display).toBe("42");
    });

    it("ignores a second decimal point", () => {
      const { result } = renderHook(() => useCalculator());
      typeDigits(result, "1.5.2");
      expect(result.current.display).toBe("1.52");
    });

    it("starts a fresh number after a result", async () => {
      mockPost.mockResolvedValue({ result: 8, resultText: "8", expression: "5 + 3" });
      const { result } = renderHook(() => useCalculator());
      typeDigits(result, "5");
      act(() => result.current.setOperation("add"));
      typeDigits(result, "3");
      await act(async () => result.current.calculate());
      expect(result.current.display).toBe("8");

      typeDigits(result, "7");
      expect(result.current.display).toBe("7");
    });
  });

  describe("binary operations", () => {
    it("sends the operation to the API and shows the result", async () => {
      mockPost.mockResolvedValue({ result: 8, resultText: "8", expression: "5 + 3" });
      const { result } = renderHook(() => useCalculator());

      typeDigits(result, "5");
      act(() => result.current.setOperation("add"));
      expect(result.current.pendingOp).toBe("add");
      typeDigits(result, "3");
      await act(async () => result.current.calculate());

      expect(mockPost).toHaveBeenCalledWith("/calculate", {
        operation: "add",
        a: "5",
        b: "3",
      });
      expect(result.current.display).toBe("8");
      expect(result.current.expression).toBe("5 + 3");
      expect(result.current.history).toEqual([
        { expression: "5 + 3", result: "8" },
      ]);
      expect(result.current.pendingOp).toBeNull();
    });

    it("chains operations, using the intermediate result as the next operand", async () => {
      mockPost.mockResolvedValueOnce({ result: 8, resultText: "8", expression: "5 + 3" });
      const { result } = renderHook(() => useCalculator());

      typeDigits(result, "5");
      act(() => result.current.setOperation("add"));
      typeDigits(result, "3");
      // second operator triggers the pending calculation
      await act(async () => result.current.setOperation("multiply"));

      expect(mockPost).toHaveBeenCalledWith("/calculate", {
        operation: "add",
        a: "5",
        b: "3",
      });
      expect(result.current.display).toBe("8");
      expect(result.current.pendingOp).toBe("multiply");

      mockPost.mockResolvedValueOnce({ result: 16, resultText: "16", expression: "8 × 2" });
      typeDigits(result, "2");
      await act(async () => result.current.calculate());

      expect(mockPost).toHaveBeenLastCalledWith("/calculate", {
        operation: "multiply",
        a: "8",
        b: "2",
      });
      expect(result.current.display).toBe("16");
      expect(result.current.history).toHaveLength(2);
    });

    it("re-selecting an operator before typing replaces the pending one", async () => {
      const { result } = renderHook(() => useCalculator());
      typeDigits(result, "5");
      act(() => result.current.setOperation("add"));
      act(() => result.current.setOperation("subtract"));
      expect(mockPost).not.toHaveBeenCalled();
      expect(result.current.pendingOp).toBe("subtract");
    });
  });

  it("calculate is a no-op without a pending operation", async () => {
    const { result } = renderHook(() => useCalculator());
    typeDigits(result, "5");
    await act(async () => result.current.calculate());
    expect(mockPost).not.toHaveBeenCalled();
  });

  it("sqrt fires immediately without a second operand", async () => {
    mockPost.mockResolvedValue({ result: 4, resultText: "4", expression: "√16" });
    const { result } = renderHook(() => useCalculator());

    typeDigits(result, "16");
    await act(async () => result.current.setOperation("sqrt"));

    expect(mockPost).toHaveBeenCalledWith("/calculate", {
      operation: "sqrt",
      a: "16",
    });
    expect(result.current.display).toBe("4");
  });

  describe("errors", () => {
    it("shows the API error message and keeps the display", async () => {
      mockPost.mockRejectedValue(new ApiError("division by zero", 422));
      const { result } = renderHook(() => useCalculator());

      typeDigits(result, "5");
      act(() => result.current.setOperation("divide"));
      typeDigits(result, "0");
      await act(async () => result.current.calculate());

      expect(result.current.error).toBe("division by zero");
      expect(result.current.loading).toBe(false);
      expect(result.current.display).toBe("0");
      expect(result.current.history).toEqual([]);
    });

    it("falls back to a generic message for non-API errors", async () => {
      mockPost.mockRejectedValue(new TypeError("fetch failed"));
      const { result } = renderHook(() => useCalculator());

      typeDigits(result, "5");
      act(() => result.current.setOperation("add"));
      typeDigits(result, "3");
      await act(async () => result.current.calculate());

      expect(result.current.error).toBe("Calculation failed");
    });

    it("typing a digit clears the error", async () => {
      mockPost.mockRejectedValue(new ApiError("division by zero", 422));
      const { result } = renderHook(() => useCalculator());

      typeDigits(result, "5");
      act(() => result.current.setOperation("divide"));
      typeDigits(result, "0");
      await act(async () => result.current.calculate());
      expect(result.current.error).toBe("division by zero");

      typeDigits(result, "1");
      expect(result.current.error).toBeNull();
    });
  });

  describe("clear", () => {
    it("clear resets the display but keeps history", async () => {
      mockPost.mockResolvedValue({ result: 8, resultText: "8", expression: "5 + 3" });
      const { result } = renderHook(() => useCalculator());

      typeDigits(result, "5");
      act(() => result.current.setOperation("add"));
      typeDigits(result, "3");
      await act(async () => result.current.calculate());

      act(() => result.current.clear());
      expect(result.current.display).toBe("0");
      expect(result.current.pendingOp).toBeNull();
      expect(result.current.history).toHaveLength(1);
    });

    it("clearAll also wipes history", async () => {
      mockPost.mockResolvedValue({ result: 8, resultText: "8", expression: "5 + 3" });
      const { result } = renderHook(() => useCalculator());

      typeDigits(result, "5");
      act(() => result.current.setOperation("add"));
      typeDigits(result, "3");
      await act(async () => result.current.calculate());

      act(() => result.current.clearAll());
      expect(result.current.history).toEqual([]);
    });
  });

  describe("toggleSign", () => {
    it("negates and un-negates the display", () => {
      const { result } = renderHook(() => useCalculator());
      typeDigits(result, "5");
      act(() => result.current.toggleSign());
      expect(result.current.display).toBe("-5");
      act(() => result.current.toggleSign());
      expect(result.current.display).toBe("5");
    });

    it("leaves zero alone", () => {
      const { result } = renderHook(() => useCalculator());
      act(() => result.current.toggleSign());
      expect(result.current.display).toBe("0");
    });
  });

  describe("backspace", () => {
    it("removes the last digit, bottoming out at zero", () => {
      const { result } = renderHook(() => useCalculator());
      typeDigits(result, "12");
      act(() => result.current.backspace());
      expect(result.current.display).toBe("1");
      act(() => result.current.backspace());
      expect(result.current.display).toBe("0");
    });

    it("collapses a single negated digit to zero", () => {
      const { result } = renderHook(() => useCalculator());
      typeDigits(result, "5");
      act(() => result.current.toggleSign());
      act(() => result.current.backspace());
      expect(result.current.display).toBe("0");
    });
  });

  describe("result display", () => {
    it("shows the backend's resultText verbatim", async () => {
      mockPost.mockResolvedValue({
        result: 0.3,
        resultText: "0.3",
        expression: "0.1 + 0.2",
      });
      const { result } = renderHook(() => useCalculator());

      typeDigits(result, "0.1");
      act(() => result.current.setOperation("add"));
      typeDigits(result, "0.2");
      await act(async () => result.current.calculate());

      expect(result.current.display).toBe("0.3");
    });

    it("preserves digits beyond float64 precision from resultText", async () => {
      // 3^35: the float64 `result` is off by 3; resultText is exact.
      mockPost.mockResolvedValue({
        result: 50031545098999704,
        resultText: "50031545098999707",
        expression: "3 ^ 35",
      });
      const { result } = renderHook(() => useCalculator());

      typeDigits(result, "3");
      act(() => result.current.setOperation("power"));
      typeDigits(result, "35");
      await act(async () => result.current.calculate());

      expect(result.current.display).toBe("50031545098999707");
      expect(result.current.history).toEqual([
        { expression: "3 ^ 35", result: "50031545098999707" },
      ]);
    });
  });
});
