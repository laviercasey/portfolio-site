import { NextRequest, NextResponse } from 'next/server';
import { getAuthToken } from '@/shared/lib/server';
import { isValidId, sanitizeBackendError } from '@/shared/api';
import { pingIndexNow } from '@/shared/lib';

const BACKEND = process.env.API_INTERNAL_URL || 'http://localhost:8080';

export async function PUT(request: NextRequest, { params }: { params: Promise<{ id: string }> }) {
  const token = await getAuthToken();
  if (!token) return NextResponse.json({ error: 'Unauthorized' }, { status: 401 });

  const { id } = await params;
  if (!isValidId(id)) return NextResponse.json({ error: 'Invalid ID' }, { status: 400 });

  const body = await request.json().catch(() => null);
  if (!body) return NextResponse.json({ error: 'Invalid JSON' }, { status: 400 });

  const res = await fetch(`${BACKEND}/api/services/faqs/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
    body: JSON.stringify(body),
  });

  if (!res.ok) return sanitizeBackendError(res.status);
  const data = await res.json();
  await pingIndexNow(['/services']);
  return NextResponse.json(data);
}

export async function DELETE(_req: NextRequest, { params }: { params: Promise<{ id: string }> }) {
  const token = await getAuthToken();
  if (!token) return NextResponse.json({ error: 'Unauthorized' }, { status: 401 });

  const { id } = await params;
  if (!isValidId(id)) return NextResponse.json({ error: 'Invalid ID' }, { status: 400 });

  const res = await fetch(`${BACKEND}/api/services/faqs/${id}`, {
    method: 'DELETE',
    headers: { Authorization: `Bearer ${token}` },
  });

  if (res.status === 204) {
    await pingIndexNow(['/services']);
    return new NextResponse(null, { status: 204 });
  }
  if (!res.ok) return sanitizeBackendError(res.status);
  await pingIndexNow(['/services']);
  return NextResponse.json({ ok: true });
}
