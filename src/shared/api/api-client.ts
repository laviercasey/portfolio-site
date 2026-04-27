export class ApiError extends Error {
  status: number;

  constructor(status: number, message: string) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
  }
}

export class ApiClient {
  private baseUrl: string;
  private token?: string;

  constructor(baseUrl: string, token?: string) {
    this.baseUrl = baseUrl;
    this.token = token;
  }

  withToken(token: string): ApiClient {
    return new ApiClient(this.baseUrl, token);
  }

  private async request<T>(path: string, options: RequestInit = {}): Promise<T> {
    const url = `${this.baseUrl}${path}`;
    const headers: Record<string, string> = {
      ...((options.headers as Record<string, string>) || {}),
    };

    if (!(options.body instanceof FormData)) {
      headers['Content-Type'] = 'application/json';
    }

    if (this.token) {
      headers['Authorization'] = `Bearer ${this.token}`;
    }

    const res = await fetch(url, { ...options, headers });

    if (!res.ok) {
      const err = await res.json().catch(() => ({ error: res.statusText }));
      throw new ApiError(res.status, err.error || res.statusText);
    }

    if (res.status === 204) return undefined as T;
    return res.json();
  }

  get<T>(path: string, init?: RequestInit): Promise<T> {
    return this.request<T>(path, { ...init, method: 'GET' });
  }

  post<T>(path: string, body?: unknown, init?: RequestInit): Promise<T> {
    return this.request<T>(path, {
      ...init,
      method: 'POST',
      body: body instanceof FormData ? body : body ? JSON.stringify(body) : undefined,
    });
  }

  put<T>(path: string, body?: unknown, init?: RequestInit): Promise<T> {
    return this.request<T>(path, {
      ...init,
      method: 'PUT',
      body: body ? JSON.stringify(body) : undefined,
    });
  }

  patch<T>(path: string, body?: unknown, init?: RequestInit): Promise<T> {
    return this.request<T>(path, {
      ...init,
      method: 'PATCH',
      body: body ? JSON.stringify(body) : undefined,
    });
  }

  delete(path: string, init?: RequestInit): Promise<void> {
    return this.request<void>(path, { ...init, method: 'DELETE' });
  }
}

export const serverApi = new ApiClient(
  process.env.API_INTERNAL_URL || 'http://localhost:8080',
);

export const clientApi = new ApiClient('');
