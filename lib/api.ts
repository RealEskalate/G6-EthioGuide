// Centralized API client for EthioGuide backend
// Base URL: proxied via Next.js rewrites to avoid CORS

export const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "/api/v1";

export type ApiError = {
  error: {
    code: string;
    message: string;
    details?: Record<string, unknown>;
  };
};

type HttpMethod = "GET" | "POST" | "PATCH" | "DELETE";

type RequestOptions = {
  method?: HttpMethod;
  headers?: Record<string, string>;
  body?: unknown;
  query?: Record<string, string | number | boolean | undefined | null>;
  authToken?: string | null;
  idempotencyKey?: string;
};

function buildQuery(params?: RequestOptions["query"]): string {
  if (!params) return "";
  const search = new URLSearchParams();
  Object.entries(params).forEach(([key, value]) => {
    if (value === undefined || value === null) return;
    search.set(key, String(value));
  });
  const qs = search.toString();
  return qs ? `?${qs}` : "";
}

export async function apiRequest<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const { method = "GET", headers = {}, body, query, authToken, idempotencyKey } = options;
  const language = localStorage.getItem("i18nextLng") || "en";
  console.log("lang",language)
  const url = `${API_BASE_URL}${path}${buildQuery(query)}`;
  const finalHeaders: HeadersInit = {
    "Content-Type": "application/json",
    ...headers,
    "lang": "am",
  };

  if (authToken) {
    (finalHeaders as Record<string, string>)["Authorization"] = `Bearer ${authToken}`;
  }
  if (idempotencyKey) {
    (finalHeaders as Record<string, string>)["Idempotency-Key"] = idempotencyKey;
  }

  const res = await fetch(url, {
    method,
    headers: finalHeaders,
    body: body !== undefined ? JSON.stringify(body) : undefined,
    cache: "no-store",
    // credentials omitted: backend is stateless JWT via Authorization header
  });

  const text = await res.text();
  let json: unknown = null;
  try {
    json = text ? JSON.parse(text) : null;
  } catch {
    // noop, will treat as text
  }

  if (!res.ok) {
    const errObj = (typeof json === 'object' && json && 'error' in (json as Record<string, unknown>))
      ? (json as { error: ApiError['error'] }).error
      : undefined
    const apiErr: ApiError = errObj
      ? { error: errObj }
      : { error: { code: `HTTP_${res.status}`, message: res.statusText || "Request failed" } };
    throw apiErr;
  }

  const data = (typeof json === 'object' && json !== null && 'data' in (json as Record<string, unknown>))
    ? (json as { data: unknown }).data
    : json
  return data as T;
}

// Convenience endpoints
export const ProceduresAPI = {
  list: (params: { page?: number; limit?: number; q?: string; category?: string; orgId?: string; sort?: string; updatedSince?: string; lang?: string }) =>
    apiRequest<{ data: unknown[]; page: number; limit: number; total: number; hasNext: boolean }>(
      "/procedures",
      { query: params }
    ),
  get: (id: string) => apiRequest<unknown>(`/procedures/${id}`),
  // authToken optional now; backend mock expected to accept unauthenticated creation
  createChecklist: (procedureId: string, authToken?: string) =>
    apiRequest<unknown>("/checklists", { method: "POST", authToken, body: { procedure_id: procedureId } }),
  // list feedback entries for a procedure (pageable)
  listFeedback: (procedureId: string, params?: { page?: number; limit?: number }, authToken?: string) =>
    apiRequest<{ data: unknown[]; total: number; page: number; limit: number; hasNext?: boolean }>(
      `/procedures/${procedureId}/feedback`,
      { query: params, authToken }
    ),
  submitFeedback: (authToken: string, procedureId: string, payload: { type: string; body: string }) =>
    apiRequest<unknown>(`/procedures/${procedureId}/feedback`, { method: "POST", authToken, body: payload }),
};

export const ChecklistsAPI = {
  getMine: (authToken?: string) => apiRequest<unknown>("/checklists/myProcedures", { authToken }),
  getOne: (userProcedureId: string, authToken?: string) => apiRequest<unknown>(`/checklists/${userProcedureId}`, { authToken }),
  patchItem: (checklistID: string, body: Record<string, unknown>, authToken?: string) =>
    apiRequest<unknown>(`/checklists/${checklistID}`, { method: "PATCH", authToken, body }),
};

export const NoticesAPI = {
  list: (params: { orgId?: string; procedureId?: string; pinned?: boolean; lang?: string }) =>
  apiRequest<unknown[]>("/notices", { query: params }),
};

export const SearchAPI = {
  search: (params: { q: string; lang?: string; sort?: string }) => apiRequest<unknown>("/search", { query: params }),
};

// Categories endpoint helper
export const CategoriesAPI = {
  list: (params?: { page?: number; limit?: number; sortOrder?: 'asc' | 'desc'; parentID?: string; organizationID?: string; title?: string }) =>
  apiRequest<{ data: unknown[]; page: number; limit: number; total: number }>("/categories", { query: params })
};


