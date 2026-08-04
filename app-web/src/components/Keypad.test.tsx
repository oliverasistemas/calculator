import { render, screen, fireEvent } from "@testing-library/react";
import { describe, it, expect, vi } from "vitest";
import Keypad from "./Keypad";

function renderKeypad(overrides = {}) {
  const props = {
    onDigit: vi.fn(),
    onOperation: vi.fn(),
    onEquals: vi.fn(),
    onClear: vi.fn(),
    onClearAll: vi.fn(),
    onToggleSign: vi.fn(),
    onBackspace: vi.fn(),
    pendingOp: null,
    loading: false,
    ...overrides,
  };
  render(<Keypad {...props} />);
  return props;
}

describe("Keypad", () => {
  it("calls onDigit when digit button clicked", () => {
    const props = renderKeypad();
    fireEvent.click(screen.getByText("7"));
    expect(props.onDigit).toHaveBeenCalledWith("7");
  });

  it("calls onEquals when = clicked", () => {
    const props = renderKeypad();
    fireEvent.click(screen.getByText("="));
    expect(props.onEquals).toHaveBeenCalled();
  });

  it("calls onClearAll when AC clicked", () => {
    const props = renderKeypad();
    fireEvent.click(screen.getByText("AC"));
    expect(props.onClearAll).toHaveBeenCalled();
  });

  it("calls onClear when C clicked", () => {
    const props = renderKeypad();
    fireEvent.click(screen.getByText("C"));
    expect(props.onClear).toHaveBeenCalled();
  });

  it("calls onToggleSign when +/- clicked", () => {
    const props = renderKeypad();
    fireEvent.click(screen.getByText("+/-"));
    expect(props.onToggleSign).toHaveBeenCalled();
  });

  it("disables buttons when loading", () => {
    renderKeypad({ loading: true });
    const equalsBtn = screen.getByText("=");
    expect(equalsBtn.closest("button")).toBeDisabled();
  });

  it("renders all digit buttons 0-9", () => {
    renderKeypad();
    for (let i = 0; i <= 9; i++) {
      expect(screen.getByText(String(i))).toBeInTheDocument();
    }
  });

  it("renders decimal point button", () => {
    renderKeypad();
    expect(screen.getByText(".")).toBeInTheDocument();
  });

  it("calls onOperation with the operation for operator buttons", () => {
    const props = renderKeypad();
    fireEvent.click(screen.getByText("÷"));
    expect(props.onOperation).toHaveBeenCalledWith("divide");
    fireEvent.click(screen.getByRole("button", { name: /plus/i }));
    expect(props.onOperation).toHaveBeenCalledWith("add");
    fireEvent.click(screen.getByText("√"));
    expect(props.onOperation).toHaveBeenCalledWith("sqrt");
  });

  it("calls onBackspace when the delete button is clicked", () => {
    const props = renderKeypad();
    fireEvent.click(screen.getByRole("button", { name: /delete/i }));
    expect(props.onBackspace).toHaveBeenCalled();
  });

  it("highlights only the pending operator", () => {
    renderKeypad({ pendingOp: "divide" });
    expect(screen.getByText("÷").closest("button")).toHaveClass("key-active");
    expect(screen.getByText("√").closest("button")).not.toHaveClass(
      "key-active"
    );
  });
});
