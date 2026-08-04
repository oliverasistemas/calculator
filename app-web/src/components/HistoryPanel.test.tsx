import { render, screen, fireEvent } from "@testing-library/react";
import { describe, it, expect, vi } from "vitest";
import HistoryPanel from "./HistoryPanel";

describe("HistoryPanel", () => {
  it("renders nothing when history is empty", () => {
    const { container } = render(
      <HistoryPanel history={[]} onClear={vi.fn()} />
    );
    expect(container.querySelector(".history")).toBeNull();
  });

  it("renders history entries", () => {
    const history = [
      { expression: "5 + 3", result: 8 },
      { expression: "10 / 2", result: 5 },
    ];
    render(<HistoryPanel history={history} onClear={vi.fn()} />);
    expect(screen.getByText("5 + 3 =")).toBeInTheDocument();
    expect(screen.getByText("8")).toBeInTheDocument();
    expect(screen.getByText("10 / 2 =")).toBeInTheDocument();
    expect(screen.getByText("5")).toBeInTheDocument();
  });

  it("formats results with float artifacts", () => {
    render(
      <HistoryPanel
        history={[{ expression: "0.1 + 0.2", result: 0.1 + 0.2 }]}
        onClear={vi.fn()}
      />
    );
    expect(screen.getByText("0.1 + 0.2 =")).toBeInTheDocument();
    expect(screen.getByText("0.3")).toBeInTheDocument();
  });

  it("calls onClear when Clear is clicked", () => {
    const onClear = vi.fn();
    render(
      <HistoryPanel
        history={[{ expression: "1 + 1", result: 2 }]}
        onClear={onClear}
      />
    );
    fireEvent.click(screen.getByText("Clear"));
    expect(onClear).toHaveBeenCalled();
  });

  it("shows History header", () => {
    render(
      <HistoryPanel
        history={[{ expression: "1 + 1", result: 2 }]}
        onClear={vi.fn()}
      />
    );
    expect(screen.getByText("History")).toBeInTheDocument();
  });
});
