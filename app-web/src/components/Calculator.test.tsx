import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import { describe, it, expect, vi, beforeEach } from "vitest";
import Calculator from "./Calculator";
import { api, ApiError } from "@/api/client";

vi.mock("@/api/client", async (importOriginal) => {
  const actual = await importOriginal<typeof import("@/api/client")>();
  return { ...actual, api: { post: vi.fn() } };
});

const mockPost = vi.mocked(api.post);

const displayValue = (container: HTMLElement) =>
  container.querySelector(".display-value");
const displayExpression = (container: HTMLElement) =>
  container.querySelector(".display-expression");

describe("Calculator", () => {
  beforeEach(() => {
    mockPost.mockReset();
  });

  it("performs a calculation end to end through the keypad", async () => {
    mockPost.mockResolvedValue({ result: 8, expression: "5 + 3" });
    const { container } = render(<Calculator />);

    fireEvent.click(screen.getByText("5"));
    fireEvent.click(screen.getByRole("button", { name: /plus/i }));
    fireEvent.click(screen.getByText("3"));
    fireEvent.click(screen.getByText("="));

    await waitFor(() => {
      expect(displayValue(container)).toHaveTextContent("8");
    });
    expect(mockPost).toHaveBeenCalledWith("/calculate", {
      operation: "add",
      a: 5,
      b: 3,
    });
    expect(displayExpression(container)).toHaveTextContent("5 + 3");
  });

  it("shows an error from the API in the display", async () => {
    mockPost.mockRejectedValue(new ApiError("division by zero", 422));
    const { container } = render(<Calculator />);

    fireEvent.click(screen.getByText("5"));
    fireEvent.click(screen.getByText("÷"));
    fireEvent.click(screen.getByText("0"));
    fireEvent.click(screen.getByText("="));

    await waitFor(() => {
      expect(screen.getByText("division by zero")).toBeInTheDocument();
    });
    expect(displayValue(container)).toHaveTextContent("0");
  });

  it("adds results to the history panel and AC clears them", async () => {
    mockPost.mockResolvedValue({ result: 8, expression: "5 + 3" });
    render(<Calculator />);

    fireEvent.click(screen.getByText("5"));
    fireEvent.click(screen.getByRole("button", { name: /plus/i }));
    fireEvent.click(screen.getByText("3"));
    fireEvent.click(screen.getByText("="));

    await waitFor(() => {
      expect(screen.getByText("History")).toBeInTheDocument();
    });
    expect(screen.getByText("5 + 3 =")).toBeInTheDocument();

    fireEvent.click(screen.getByText("AC"));
    expect(screen.queryByText("History")).not.toBeInTheDocument();
  });
});
