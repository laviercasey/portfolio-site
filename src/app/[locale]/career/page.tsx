import type { Metadata } from 'next';
import { setRequestLocale } from 'next-intl/server';
import { careerService } from '@/entities/career';
import { projectsService } from '@/entities/project';
import { JsonLd, buildBreadcrumbList } from '@/shared/lib';
import { siteConfig as siteJson } from '@/shared/config';
import type { SiteConfig } from '@/shared/types';
import { CareerPage } from '@/widgets/career-page';

const site = siteJson as SiteConfig;

export const revalidate = 3600;

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string }>;
}): Promise<Metadata> {
  const { locale } = await params;
  const isRu = locale === 'ru';
  const title = isRu
    ? 'Опыт и образование — Full-Stack Developer'
    : 'Experience & Education — Full-Stack Developer';
  const description = isRu
    ? 'Путь Casey Laviere: образование, системное администрирование, full-stack разработка, ML-исследования, фриланс-проекты для клиентов из России и зарубежья.'
    : 'Casey Laviere background: formal education, sysadmin years, full-stack engineering, ML research, and freelance work for clients in Russia and worldwide.';
  const localizedName = isRu ? (site.nameRu ?? site.name) : site.name;
  return {
    title,
    description,
    keywords: isRu
      ? ['резюме разработчика', 'опыт Full-Stack', 'образование программиста', 'карьера в IT', 'Casey Laviere резюме', 'системный администратор', 'фриланс проекты']
      : ['developer resume', 'full-stack experience', 'software engineer background', 'IT career', 'Casey Laviere CV', 'sysadmin experience', 'freelance projects'],
    alternates: {
      canonical: `/${locale}/career`,
      languages: { ru: '/ru/career', en: '/en/career', 'x-default': '/ru/career' },
    },
    openGraph: {
      title,
      description,
      type: 'website',
      url: `${site.url}/${locale}/career`,
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

export default async function CareerRoute({ params }: { params: Promise<{ locale: string }> }) {
  const { locale } = await params;
  setRequestLocale(locale);
  const [career, projects] = await Promise.all([
    careerService.getAll(),
    projectsService.list(),
  ]);

  const isRu = locale === 'ru';
  const breadcrumb = buildBreadcrumbList(locale, [
    { name: isRu ? 'Главная' : 'Home', path: '' },
    { name: isRu ? 'Опыт и образование' : 'Experience & Education', path: '/career' },
  ]);

  return (
    <>
      <JsonLd data={breadcrumb} />
      <CareerPage career={career} locale={locale} projects={projects} />
    </>
  );
}
