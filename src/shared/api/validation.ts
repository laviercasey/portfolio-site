import { NextRequest, NextResponse } from 'next/server';

const SAFE_ID_REGEX = /^[a-zA-Z0-9_-]+$/;

export function isValidId(id: string): boolean {
  return SAFE_ID_REGEX.test(id) && id.length <= 128;
}

export const ALLOWED_CONTENT_SECTIONS = ['homepage', 'contact', 'career'] as const;
export type ContentSection = (typeof ALLOWED_CONTENT_SECTIONS)[number];

export function isValidContentSection(section: string): section is ContentSection {
  return (ALLOWED_CONTENT_SECTIONS as readonly string[]).includes(section);
}

export const ALLOWED_CATEGORIES = ['web', 'mobile', 'data', 'research', 'other'] as const;

export function isValidCategory(category: string): boolean {
  return (ALLOWED_CATEGORIES as readonly string[]).includes(category);
}

export function sanitizeBackendError(status: number): NextResponse {
  const message =
    status === 401 ? 'Unauthorized'
    : status === 403 ? 'Forbidden'
    : status === 404 ? 'Not found'
    : status === 400 ? 'Bad request'
    : status === 409 ? 'Conflict'
    : 'An error occurred';
  return NextResponse.json({ error: message }, { status });
}

export function checkOrigin(request: NextRequest): boolean {
  const siteUrl = process.env.NEXT_PUBLIC_SITE_URL;
  if (!siteUrl) return false;

  const origin = request.headers.get('origin');
  if (origin) return origin === siteUrl;

  const referer = request.headers.get('referer');
  if (!referer) return false;
  try {
    return new URL(referer).origin === new URL(siteUrl).origin;
  } catch {
    return false;
  }
}

const rateLimitStore = new Map<string, { count: number; resetAt: number }>();

const CLEANUP_INTERVAL = 60_000;
let lastCleanup = Date.now();

function cleanupExpired() {
  const now = Date.now();
  if (now - lastCleanup < CLEANUP_INTERVAL) return;
  lastCleanup = now;
  for (const [key, entry] of rateLimitStore) {
    if (entry.resetAt < now) rateLimitStore.delete(key);
  }
}

export function rateLimit(
  ip: string,
  key: string,
  maxRequests: number,
  windowMs: number,
): boolean {
  cleanupExpired();
  const storeKey = `${key}:${ip}`;
  const now = Date.now();
  const entry = rateLimitStore.get(storeKey);

  if (!entry || entry.resetAt < now) {
    rateLimitStore.set(storeKey, { count: 1, resetAt: now + windowMs });
    return true;
  }

  if (entry.count >= maxRequests) return false;
  entry.count++;
  return true;
}

export function getClientIp(request: NextRequest): string {
  return request.headers.get('x-real-ip') ?? 'unknown';
}
