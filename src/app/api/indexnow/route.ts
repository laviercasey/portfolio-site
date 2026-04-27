import { NextRequest, NextResponse } from 'next/server';
import { timingSafeEqual } from 'node:crypto';
import { siteConfig as siteJson } from '@/shared/config';
import type { SiteConfig } from '@/shared/types';

const site = siteJson as SiteConfig;
const MAX_URLS = 10_000;

function safeEqualStrings(a: string, b: string): boolean {
  const bufA = Buffer.from(a);
  const bufB = Buffer.from(b);
  if (bufA.length !== bufB.length) return false;
  return timingSafeEqual(bufA, bufB);
}

export async function POST(request: NextRequest) {
  const key = process.env.INDEXNOW_KEY;
  if (!key) {
    return NextResponse.json({ error: 'INDEXNOW_KEY not configured' }, { status: 500 });
  }

  const apiSecret = process.env.INDEXNOW_API_SECRET;
  const authHeader = request.headers.get('authorization') ?? '';
  const expected = `Bearer ${apiSecret ?? ''}`;
  if (!apiSecret || !safeEqualStrings(authHeader, expected)) {
    return NextResponse.json({ error: 'Unauthorized' }, { status: 401 });
  }

  let body: unknown;
  try {
    body = await request.json();
  } catch {
    return NextResponse.json({ error: 'Invalid JSON' }, { status: 400 });
  }

  const { urls } = body as { urls?: unknown };
  if (!Array.isArray(urls) || urls.length === 0) {
    return NextResponse.json({ error: 'urls array is required' }, { status: 400 });
  }

  if (urls.length > MAX_URLS) {
    return NextResponse.json({ error: `Too many URLs (max ${MAX_URLS})` }, { status: 400 });
  }

  const ownHost = new URL(site.url).host;
  const allValid = urls.every((u: unknown) => {
    if (typeof u !== 'string') return false;
    try { return new URL(u).host === ownHost; } catch { return false; }
  });
  if (!allValid) {
    return NextResponse.json({ error: 'All URLs must belong to the site domain' }, { status: 400 });
  }

  const response = await fetch('https://api.indexnow.org/indexnow', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json; charset=utf-8' },
    body: JSON.stringify({ host: ownHost, key, urlList: urls }),
  });

  if (!response.ok) {
    return NextResponse.json(
      { error: 'IndexNow API error', status: response.status },
      { status: 502 },
    );
  }

  return NextResponse.json({ submitted: urls.length });
}
