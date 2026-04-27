import { siteConfig as siteJson } from '@/shared/config';
import type { SiteConfig } from '@/shared/types';
import { getAdminSlug } from '@/shared/lib';

const site = siteJson as SiteConfig;

export async function GET() {
  const adminSlug = getAdminSlug();
  const adminPath = `/${adminSlug}/`;

  const body = `User-agent: *
Allow: /
Disallow: ${adminPath}
Disallow: /api/

User-agent: Yandex
Allow: /
Disallow: ${adminPath}
Disallow: /api/
Clean-param: utm_source&utm_medium&utm_campaign&utm_term&utm_content&ref&gclid&yclid&fbclid&_openstat /
Crawl-delay: 2

Sitemap: ${site.url}/sitemap.xml
Host: ${new URL(site.url).host}`;

  return new Response(body, {
    headers: { 'Content-Type': 'text/plain' },
  });
}
