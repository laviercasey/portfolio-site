import type { Metadata } from 'next';
import { setRequestLocale } from 'next-intl/server';
import { Suspense } from 'react';
import { projectsService } from '@/entities/project';
import { ProjectCard } from '@/entities/project';
import { ProjectFilter } from '@/features/project-filter';
import { JsonLd, buildBreadcrumbList } from '@/shared/lib';
import { siteConfig as siteJson } from '@/shared/config';
import type { SiteConfig } from '@/shared/types';

const site = siteJson as SiteConfig;

export const revalidate = 600;

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string }>;
}): Promise<Metadata> {
  const { locale } = await params;
  const isRu = locale === 'ru';
  const title = isRu
    ? 'Проекты — PHP, Python, ML, Telegram-боты'
    : 'Projects — PHP, Python, ML & Telegram bots';
  const description = isRu
    ? 'Проекты Casey Laviere: веб-приложения на Laravel и Next.js, ML-модели, Telegram-боты и инструменты анализа данных. Реальные кейсы с исходниками и метриками.'
    : 'Casey Laviere projects: Laravel and Next.js web apps, ML models, Telegram bots, and data analysis tools. Real case studies with sources and measurable outcomes.';
  const localizedName = isRu ? (site.nameRu ?? site.name) : site.name;
  return {
    title,
    description,
    keywords: isRu
      ? ['проекты разработчика', 'портфолио PHP', 'портфолио Python', 'Laravel проекты', 'Next.js проекты', 'кейсы разработки', 'ML проекты', 'Telegram боты']
      : ['developer projects', 'PHP portfolio', 'Python portfolio', 'Laravel projects', 'Next.js projects', 'case studies', 'ML projects', 'Telegram bots'],
    alternates: {
      canonical: `/${locale}/projects`,
      languages: { ru: '/ru/projects', en: '/en/projects', 'x-default': '/ru/projects' },
    },
    openGraph: {
      title,
      description,
      type: 'website',
      url: `${site.url}/${locale}/projects`,
      locale: isRu ? 'ru_RU' : 'en_US',
      siteName: localizedName,
    },
    twitter: {
      card: 'summary_large_image',
      title,
      description,
    },
  };
}

interface ProjectsPageProps {
  params: Promise<{ locale: string }>;
  searchParams: Promise<{ category?: string; status?: string; tag?: string }>;
}

export default async function ProjectsPage({ params, searchParams }: ProjectsPageProps) {
  const { locale } = await params;
  setRequestLocale(locale);
  const filters = await searchParams;
  const allProjects = await projectsService.list();

  const filtered = allProjects
    .filter((p) => {
      if (filters.category && filters.category !== 'all') {
        if (p.category !== filters.category) return false;
      }
      if (filters.status) {
        if (p.status !== filters.status) return false;
      }
      if (filters.tag) {
        if (!p.tags.includes(filters.tag)) return false;
      }
      return true;
    })
    .sort((a, b) => a.order - b.order);

  const isRu = locale === 'ru';
  const breadcrumb = buildBreadcrumbList(locale, [
    { name: isRu ? 'Главная' : 'Home', path: '' },
    { name: isRu ? 'Проекты' : 'Projects', path: '/projects' },
  ]);

  return (
    <div className="container py-24">
      <JsonLd data={breadcrumb} />
      <div className="mb-14">
        <h1 className="font-heading text-5xl font-bold mb-3 gradient-text">
          {locale === 'ru' ? 'Проекты' : 'Projects'}
        </h1>
        <p className="text-lg text-muted-foreground">
          {locale === 'ru'
            ? 'Веб-приложения, ML-модели, боты и инструменты для данных'
            : 'Web apps, ML models, bots, and data tools I\'ve built'}
        </p>
      </div>

      <Suspense fallback={null}>
        <ProjectFilter />
      </Suspense>

      {filtered.length === 0 ? (
        <div className="text-center py-24 text-muted-foreground glass-card rounded-2xl">
          {locale === 'ru' ? 'Проекты не найдены' : 'No projects found'}
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
          {filtered.map((project) => (
            <ProjectCard key={project.id} project={project} />
          ))}
        </div>
      )}
    </div>
  );
}
