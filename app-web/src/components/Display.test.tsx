import { render, screen } from "@testing-library/react";
import { describe, it, expect } from "vitest";
import Display from "./Display";

describe("Display", () => {
  it("renders the value", () => {
    render(<Display value="42" expression="" error={null} />);
    expect(screen.getByText("42")).toBeInTheDocument();
  });

  it("renders the expression", () => {
    render(<Display value="8" expression="5 + 3" error={null} />);
    expect(screen.getByText("5 + 3")).toBeInTheDocument();
  });

  it("renders an error message", () => {
    render(<Display value="0" expression="" error="division by zero" />);
    expect(screen.getByText("division by zero")).toBeInTheDocument();
  });

  it("renders non-breaking space when expression is empty", () => {
    const { container } = render(
      <Display value="0" expression="" error={null} />
    );
    const expr = container.querySelector(".display-expression");
    expect(expr?.textContent).toBe("\u00A0");
  });
});
