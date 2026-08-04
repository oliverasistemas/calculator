import { describe, it, expect } from "vitest";
import { formatResult } from "./format";

describe("formatResult", () => {
  it("rounds float artifacts to a clean value", () => {
    expect(formatResult(0.1 + 0.2)).toBe("0.3");
  });

  it("keeps integers unchanged", () => {
    expect(formatResult(8)).toBe("8");
    expect(formatResult(-42)).toBe("-42");
  });

  it("keeps ordinary decimals unchanged", () => {
    expect(formatResult(2.5)).toBe("2.5");
    expect(formatResult(-0.125)).toBe("-0.125");
  });

  it("formats large non-safe integers", () => {
    expect(formatResult(1e30)).toBe("1e+30");
  });

  it("trims to 12 significant digits", () => {
    expect(formatResult(Math.sqrt(2))).toBe("1.41421356237");
  });
});
