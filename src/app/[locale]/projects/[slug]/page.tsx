import type { Metadata } from 'next';
import { setRequestLocale } from 'next-intl/server';
import { notFound } from 'next/navigation';
import { projectsService } from '@/entities/project';
import { siteConfig as siteJson } from '@/shared/config';
import type { SiteConfig } from '@/shared/types';
import { ProjectDetail } from '@/widgets/project-detail';

export const revalidate = 600;

const site = siteJson as SiteConfig;

interface ProjectPageProps {
  params: Promise<{ locale: string; slug: string }>;
}

export async function generateMetadata({ params }: ProjectPageProps): Promise<Metadata> {
  const { locale, slug } = await params;
  const isRu = locale === 'ru';

  try {
    const project = await projectsService.getBySlug(slug);
    if (!project) return {};

    const title = isRu ? project.title.ru : project.title.en;
    const description = isRu ? project.shortDescription.ru : project.shortDescription.en;
    const published = project.createdAt;
    const modified = project.createdAt;

    return {
      title,
      description,
      keywords: [...project.tags, ...project.techStack],
      alternates: {
        canonical: `/${locale}/projects/${slug}`,
        languages: {
          ru: `/ru/projects/${slug}`,
          en: `/en/projects/${slug}`,
          'x-default': `/ru/projects/${slug}`,
        },
      },
      openGraph: {
        title: `${title} | ${site.name}`,
        description,
        type: 'article',
        publishedTime: published,
        modifiedTime: modified,
        authors: [site.url],
        tags: project.tags,
        url: `${site.url}/${locale}/projects/${slug}`,
        siteName: site.name,
        locale: isRu ? 'ru_RU' : 'en_US',
      },
      twitter: {
        card: 'summary_large_image',
        title: `${title} — ${site.name}`,
        description,
      },
    };
  } catch {
    return {};
  }
}

export function generateStaticParams() {
  return [];
}

export default async function ProjectPage({ params }: ProjectPageProps) {
  const { locale, slug } = await params;
  setRequestLocale(locale);

  const project = await projectsService.getBySlug(slug);
  if (!project) notFound();

  const isRu = locale === 'ru';
  const title = isRu ? project.title.ru : project.title.en;
  const description = isRu ? project.shortDescription.ru : project.shortDescription.en;
  const articleBody = isRu ? project.description.ru : project.description.en;
  const published = project.createdAt;
  const modified = project.createdAt;
  const pageUrl = `${site.url}/${locale}/projects/${slug}`;
  const ogUrl = `${pageUrl}/opengraph-image`;
  const image = project.thumbnailUrl
    ? (project.thumbnailUrl.startsWith('http') ? project.thumbnailUrl : `${site.url}${project.thumbnailUrl}`)
    : ogUrl;

  const articleJsonLd = {
    '@context': 'https://schema.org',
    '@type': 'Article',
    headline: title,
    description,
    image,
    datePublished: new Date(published).toISOString(),
    dateModified: new Date(modified).toISOString(),
    inLanguage: isRu ? 'ru-RU' : 'en-US',
    keywords: [...project.tags, ...project.techStack].join(', '),
    articleBody,
    mainEntityOfPage: { '@type': 'WebPage', '@id': pageUrl },
    author: { '@type': 'Person', name: site.name, url: site.url },
    publisher: {
      '@type': 'Person',
      name: site.name,
      url: site.url,
      logo: { '@type': 'ImageObject', url: `${site.url}${site.ogImage}` },
    },
  };

  const breadcrumbJsonLd = {
    '@context': 'https://schema.org',
    '@type': 'BreadcrumbList',
    itemListElement: [
      { '@type': 'ListItem', position: 1, name: isRu ? 'Главная' : 'Home', item: `${site.url}/${locale}` },
      { '@type': 'ListItem', position: 2, name: isRu ? 'Проекты' : 'Projects', item: `${site.url}/${locale}/projects` },
      { '@type': 'ListItem', position: 3, name: title, item: pageUrl },
    ],
  };

  return (
    <>
      <script
        type="application/ld+json"
        dangerouslySetInnerHTML={{ __html: JSON.stringify(articleJsonLd).replace(/</g, '\\u003c') }}
      />
      <script
        type="application/ld+json"
        dangerouslySetInnerHTML={{ __html: JSON.stringify(breadcrumbJsonLd).replace(/</g, '\\u003c') }}
      />
      <ProjectDetail project={project} />
    </>
  );
}
