import { siteConfig as siteJson } from '@/shared/config';
import type { SiteConfig } from '@/shared/types';

const site = siteJson as SiteConfig;
const LOCALES = ['ru', 'en'] as const;

export async function pingIndexNow(paths: string[]): Promise<void> {
  const key = process.env.INDEXNOW_KEY;
  if (!key) return;

  const host = new URL(site.url).host;
  const urlList = paths
    .flatMap((p) => LOCALES.map((l) => `${site.url}/${l}${p.startsWith('/') ? p : `/${p}`}`))
    .filter((u, i, arr) => arr.indexOf(u) === i)
    .slice(0, 10_000);

  if (urlList.length === 0) return;

  try {
    await fetch('https://api.indexnow.org/indexnow', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json; charset=utf-8' },
      body: JSON.stringify({ host, key, urlList }),
      signal: AbortSignal.timeout(5_000),
    });
  } catch {
  }
}
