const SLUG_REGEX = /^[a-zA-Z0-9_-]+$/;

let cached: string | null = null;

export function getAdminSlug(): string {
  if (cached) return cached;

  const raw = process.env.ADMIN_SLUG?.trim() || 'adm';

  if (!SLUG_REGEX.test(raw) || raw.length > 64) {
    throw new Error(
      `Invalid ADMIN_SLUG: must match ${SLUG_REGEX} and be ≤64 chars. Got: ${JSON.stringify(raw)}`,
    );
  }

  cached = raw;
  return raw;
}
