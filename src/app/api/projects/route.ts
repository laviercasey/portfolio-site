import { NextRequest, NextResponse } from 'next/server';
import { getAuthToken } from '@/shared/lib/server';
import { isValidCategory, sanitizeBackendError } from '@/shared/api';
import { pingIndexNow } from '@/shared/lib';

const BACKEND = process.env.API_INTERNAL_URL || 'http://localhost:8080';

export async function GET(request: NextRequest) {
  const { searchParams } = new URL(request.url);
  const rawCategory = searchParams.get('category');
  const category = rawCategory && isValidCategory(rawCategory) ? rawCategory : null;
  const q = category ? `?category=${encodeURIComponent(category)}` : '';

  const res = await fetch(`${BACKEND}/api/projects${q}`, {
    next: { revalidate: 60 },
  });
  if (!res.ok) return sanitizeBackendError(res.status);
  const data = await res.json();
  return NextResponse.json(data);
}

export async function POST(request: NextRequest) {
  const token = await getAuthToken();
  if (!token) return NextResponse.json({ error: 'Unauthorized' }, { status: 401 });

  const body = await request.json().catch(() => null);
  if (!body) return NextResponse.json({ error: 'Invalid JSON' }, { status: 400 });

  const res = await fetch(`${BACKEND}/api/projects`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(body),
  });

  if (!res.ok) return sanitizeBackendError(res.status);
  const data = await res.json();

  const slug = typeof data?.slug === 'string' ? data.slug : undefined;
  await pingIndexNow(['/projects', ...(slug ? [`/projects/${slug}`] : [])]);

  return NextResponse.json(data);
}
