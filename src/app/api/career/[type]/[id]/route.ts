import { NextRequest, NextResponse } from 'next/server';
import { getAuthToken } from '@/shared/lib/server';
import { isValidId, sanitizeBackendError } from '@/shared/api';

const BACKEND = process.env.API_INTERNAL_URL || 'http://localhost:8080';

const ALLOWED_TYPES = ['education', 'work', 'certificate', 'publication'] as const;

function isValidType(t: string): boolean {
  return (ALLOWED_TYPES as readonly string[]).includes(t);
}

export async function PUT(
  request: NextRequest,
  { params }: { params: Promise<{ type: string; id: string }> },
) {
  const token = await getAuthToken();
  if (!token) return NextResponse.json({ error: 'Unauthorized' }, { status: 401 });

  const { type, id } = await params;
  if (!isValidType(type)) return NextResponse.json({ error: 'Invalid type' }, { status: 400 });
  if (!isValidId(id)) return NextResponse.json({ error: 'Invalid ID' }, { status: 400 });

  const body = await request.json().catch(() => null);
  if (!body) return NextResponse.json({ error: 'Invalid JSON' }, { status: 400 });

  const res = await fetch(`${BACKEND}/api/career/${type}/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(body),
  });

  if (!res.ok) return sanitizeBackendError(res.status);
  const data = await res.json();
  return NextResponse.json(data);
}

export async function DELETE(
  _req: NextRequest,
  { params }: { params: Promise<{ type: string; id: string }> },
) {
  const token = await getAuthToken();
  if (!token) return NextResponse.json({ error: 'Unauthorized' }, { status: 401 });

  const { type, id } = await params;
  if (!isValidType(type)) return NextResponse.json({ error: 'Invalid type' }, { status: 400 });
  if (!isValidId(id)) return NextResponse.json({ error: 'Invalid ID' }, { status: 400 });

  const res = await fetch(`${BACKEND}/api/career/${type}/${id}`, {
    method: 'DELETE',
    headers: { Authorization: `Bearer ${token}` },
  });

  if (res.status === 204) return new NextResponse(null, { status: 204 });
  if (!res.ok) return sanitizeBackendError(res.status);
  const data = await res.json().catch(() => ({ ok: true }));
  return NextResponse.json(data);
}
