interface LoginResponse {
  ok: boolean;
  error?: string;
}

export async function adminLogin(password: string): Promise<LoginResponse> {
  try {
    const res = await fetch('/api/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ password }),
    });
    const data = (await res.json().catch(() => ({}))) as { error?: string };
    if (!res.ok) {
      return { ok: false, error: data.error ?? 'Login failed' };
    }
    return { ok: true };
  } catch {
    return { ok: false, error: 'Connection error' };
  }
}
