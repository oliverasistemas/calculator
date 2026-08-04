import { render } from "@testing-library/react";
import { describe, it, expect, vi } from "vitest";
import App from "./App";

vi.mock("@/api/client", async (importOriginal) => {
  const actual = await importOriginal<typeof import("@/api/client")>();
  return { ...actual, api: { post: vi.fn() } };
});

describe("App", () => {
  it("renders the calculator inside the themed shell", () => {
    const { container } = render(<App />);
    expect(container.querySelector(".app-container")).toBeInTheDocument();
    expect(container.querySelector(".display-value")).toHaveTextContent("0");
  });
});
