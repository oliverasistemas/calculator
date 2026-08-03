import { API_BASE_URL } from "@/constants";

const REQUEST_TIMEOUT_MS = 30000;

export class ApiError extends Error {
  status: number;
  data?: unknown;

  constructor(message: string, status: number, data?: unknown) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.data = data;
  }
}

async function request<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const url = `${API_BASE_URL}${endpoint}`;
  const controller = new AbortController();
  const timeout = setTimeout(() => controller.abort(), REQUEST_TIMEOUT_MS);

  try {
    const response = await fetch(url, {
      ...options,
      signal: controller.signal,
      headers: {
        "Content-Type": "application/json",
        ...options.headers,
      },
    });

    const data = response.status === 204 ? null : await response.json();

    if (!response.ok) {
      throw new ApiError(
        data.error || data.message || "Request failed",
        response.status,
        data
      );
    }

    return data as T;
  } finally {
    clearTimeout(timeout);
  }
}

export const api = {
  post: <T>(endpoint: string, body?: unknown) =>
    request<T>(endpoint, {
      method: "POST",
      body: body ? JSON.stringify(body) : undefined,
    }),
};
