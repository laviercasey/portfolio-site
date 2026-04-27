const SAFE_SCHEMES = /^(https?:|mailto:|tel:)/i;

export function isSafeUrl(value: unknown): value is string {
  return typeof value === 'string' && SAFE_SCHEMES.test(value.trim());
}

export function safeHref(value: string): string {
  return isSafeUrl(value) ? value : '#';
}
