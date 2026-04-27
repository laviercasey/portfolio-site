import { NextRequest, NextResponse } from 'next/server';
import { getAuthToken } from '@/shared/lib/server';
import { isValidContentSection, sanitizeBackendError } from '@/shared/api';
import { isSafeUrl } from '@/shared/lib';

const BACKEND = process.env.API_INTERNAL_URL || 'http://localhost:8080';

function hasUnsafeSocialLink(body: unknown): boolean {
  if (!body || typeof body !== 'object') return false;
  const links = (body as { socialLinks?: unknown }).socialLinks;
  if (!Array.isArray(links)) return false;
  return links.some((l) => l && typeof l === 'object' && !isSafeUrl((l as { url?: unknown }).url));
}

export async function GET() {
  const res = await fetch(`${BACKEND}/api/content`, {
    next: { revalidate: 60 },
  });
  if (!res.ok) return sanitizeBackendError(res.status);
  const data = await res.json();
  return NextResponse.json(data);
}

export async function PUT(request: NextRequest) {
  const token = await getAuthToken();
  if (!token) return NextResponse.json({ error: 'Unauthorized' }, { status: 401 });

  const { searchParams } = new URL(request.url);
  const section = searchParams.get('section') ?? 'homepage';

  if (!isValidContentSection(section)) {
    return NextResponse.json({ error: 'Invalid section' }, { status: 400 });
  }

  const body = await request.json().catch(() => null);
  if (!body) return NextResponse.json({ error: 'Invalid JSON' }, { status: 400 });

  if (hasUnsafeSocialLink(body)) {
    return NextResponse.json(
      { error: 'Social link URLs must use http(s), mailto, or tel scheme' },
      { status: 400 },
    );
  }

  const res = await fetch(`${BACKEND}/api/content/${section}`, {
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
