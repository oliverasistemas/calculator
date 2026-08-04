import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { api, ApiError } from "./client";
import { API_BASE_URL } from "@/constants";

const fetchMock = vi.fn();

function jsonResponse(body: unknown, status = 200) {
  return {
    ok: status >= 200 && status < 300,
    status,
    json: () => Promise.resolve(body),
  };
}

describe("api client", () => {
  beforeEach(() => {
    fetchMock.mockReset();
    vi.stubGlobal("fetch", fetchMock);
  });

  afterEach(() => {
    vi.unstubAllGlobals();
    vi.useRealTimers();
  });

  it("POSTs JSON to the endpoint and returns the parsed body", async () => {
    fetchMock.mockResolvedValue(jsonResponse({ result: 8, expression: "5 + 3" }));

    const data = await api.post("/calculate", { operation: "add", a: 5, b: 3 });

    expect(data).toEqual({ result: 8, expression: "5 + 3" });
    expect(fetchMock).toHaveBeenCalledTimes(1);
    const [url, init] = fetchMock.mock.calls[0];
    expect(url).toBe(`${API_BASE_URL}/calculate`);
    expect(init.method).toBe("POST");
    expect(init.headers["Content-Type"]).toBe("application/json");
    expect(JSON.parse(init.body)).toEqual({ operation: "add", a: 5, b: 3 });
  });

  it("sends no body when none is given", async () => {
    fetchMock.mockResolvedValue(jsonResponse({ ok: true }));

    await api.post("/ping");

    const [, init] = fetchMock.mock.calls[0];
    expect(init.body).toBeUndefined();
  });

  it("returns null for a 204 response without reading the body", async () => {
    fetchMock.mockResolvedValue({
      ok: true,
      status: 204,
      json: () => Promise.reject(new Error("no body")),
    });

    await expect(api.post("/calculate")).resolves.toBeNull();
  });

  it("throws an ApiError with the server's `error` field", async () => {
    fetchMock.mockResolvedValue(
      jsonResponse({ error: "division by zero" }, 422)
    );

    const err = await api
      .post<never>("/calculate")
      .catch((e: unknown) => e as ApiError);

    expect(err).toBeInstanceOf(ApiError);
    expect(err.message).toBe("division by zero");
    expect(err.status).toBe(422);
    expect(err.data).toEqual({ error: "division by zero" });
  });

  it("falls back to the `message` field, then a generic message", async () => {
    fetchMock.mockResolvedValue(jsonResponse({ message: "bad input" }, 400));
    await expect(api.post("/calculate")).rejects.toThrow("bad input");

    fetchMock.mockResolvedValue(jsonResponse({}, 500));
    await expect(api.post("/calculate")).rejects.toThrow("Request failed");
  });

  it("aborts the request after the 30s timeout", async () => {
    vi.useFakeTimers();
    fetchMock.mockImplementation(
      (_url: string, init: RequestInit) =>
        new Promise((_resolve, reject) => {
          init.signal?.addEventListener("abort", () =>
            reject(new DOMException("The operation was aborted", "AbortError"))
          );
        })
    );

    const pending = api.post("/calculate");
    const assertion = expect(pending).rejects.toThrow(/aborted/i);
    await vi.advanceTimersByTimeAsync(30000);
    await assertion;
  });
});
