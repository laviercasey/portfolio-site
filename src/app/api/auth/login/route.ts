import { NextRequest, NextResponse } from 'next/server';
import { cookies } from 'next/headers';
import { rateLimit, getClientIp } from '@/shared/api';

const BACKEND = process.env.API_INTERNAL_URL || 'http://localhost:8080';

export async function POST(request: NextRequest) {
  const ip = getClientIp(request);
  if (!rateLimit(ip, 'login', 5, 15 * 60 * 1000)) {
    return NextResponse.json({ error: 'Too many attempts' }, { status: 429 });
  }

  const body = await request.json().catch(() => null);
  if (!body?.password || typeof body.password !== 'string') {
    return NextResponse.json({ error: 'Password required' }, { status: 400 });
  }

  const res = await fetch(`${BACKEND}/api/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ password: body.password }),
  });

  if (!res.ok) {
    return NextResponse.json({ error: 'Invalid credentials' }, { status: res.status });
  }

  const { token, expiresAt } = await res.json();

  const cookieStore = await cookies();
  cookieStore.set('auth-token', token, {
    httpOnly: true,
    secure: process.env.COOKIE_SECURE !== 'false',
    sameSite: 'strict',
    path: '/',
    expires: new Date(expiresAt),
  });

  return NextResponse.json({ ok: true });
}
