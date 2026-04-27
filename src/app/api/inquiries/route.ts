import { NextRequest, NextResponse } from 'next/server';
import { z } from 'zod';
import { getAuthToken } from '@/shared/lib/server';
import {
  sanitizeBackendError,
  checkOrigin,
  rateLimit,
  getClientIp,
} from '@/shared/api';

const BACKEND = process.env.API_INTERNAL_URL || 'http://localhost:8080';

const TELEGRAM_NICK = /^@?[A-Za-z0-9_]{5,32}$/;

const inquirySchema = z.object({
  name: z.string().min(1).max(200),
  email: z.string().email().max(320),
  company: z.string().max(200).optional(),
  telegram: z.union([z.literal(''), z.string().regex(TELEGRAM_NICK)]).optional(),
  type: z.enum(['freelance', 'fulltime', 'collaboration', 'other']),
  budget: z.string().max(100).optional(),
  message: z.string().min(10).max(5000),
}).strict();

export async function GET() {
  const token = await getAuthToken();
  if (!token) return NextResponse.json({ error: 'Unauthorized' }, { status: 401 });

  const res = await fetch(`${BACKEND}/api/inquiries`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  if (!res.ok) return sanitizeBackendError(res.status);
  const data = await res.json();
  return NextResponse.json(data);
}

export async function POST(request: NextRequest) {
  if (!checkOrigin(request)) {
    return NextResponse.json({ error: 'Forbidden' }, { status: 403 });
  }

  const ip = getClientIp(request);
  if (!rateLimit(ip, 'inquiries', 5, 15 * 60 * 1000)) {
    return NextResponse.json({ error: 'Too many requests' }, { status: 429 });
  }

  const body = await request.json().catch(() => null);
  if (!body) return NextResponse.json({ error: 'Invalid JSON' }, { status: 400 });

  const parsed = inquirySchema.safeParse(body);
  if (!parsed.success) {
    return NextResponse.json({ error: 'Validation failed' }, { status: 400 });
  }

  const res = await fetch(`${BACKEND}/api/inquiries`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(parsed.data),
  });

  if (!res.ok) return sanitizeBackendError(res.status);
  const data = await res.json();
  return NextResponse.json(data, { status: 201 });
}
