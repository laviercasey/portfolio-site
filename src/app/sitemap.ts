import { MetadataRoute } from 'next';
import { siteConfig as siteJson } from '@/shared/config';
import type { SiteConfig } from '@/shared/types';
import { projectsService } from '@/entities/project';

export const revalidate = 3600;

const LOCALES = ['ru', 'en'] as const;

function buildAlternates(path: string, baseUrl: string) {
  const languages: Record<string, string> = {};
  for (const locale of LOCALES) {
    languages[locale] = `${baseUrl}/${locale}${path}`;
  }
  languages['x-default'] = `${baseUrl}/en${path}`;
  return { languages };
}

export default async function sitemap(): Promise<MetadataRoute.Sitemap> {
  const site = siteJson as SiteConfig;

  let projects: Awaited<ReturnType<typeof projectsService.list>> = [];
  try {
    projects = await projectsService.list();
  } catch {
  }

  const staticRoutes = ['', '/projects', '/career', '/contact'] as const;

  const staticEntries = LOCALES.flatMap((locale) =>
    staticRoutes.map((route) => ({
      url: `${site.url}/${locale}${route}`,
      lastModified: new Date(),
      changeFrequency: 'weekly' as const,
      priority: route === '' ? 1 : 0.8,
      alternates: buildAlternates(route, site.url),
    })),
  );

  const projectEntries = LOCALES.flatMap((locale) =>
    projects.map((project) => ({
      url: `${site.url}/${locale}/projects/${project.slug}`,
      lastModified: new Date(project.createdAt),
      changeFrequency: 'monthly' as const,
      priority: 0.6,
      alternates: buildAlternates(`/projects/${project.slug}`, site.url),
    })),
  );

  return [...staticEntries, ...projectEntries];
}
