import type { Metadata } from 'next';
import { setRequestLocale } from 'next-intl/server';
import { contentService } from '@/entities/content';
import { JsonLd, buildBreadcrumbList } from '@/shared/lib';
import { siteConfig as siteJson } from '@/shared/config';
import type { SiteConfig } from '@/shared/types';
import { ContactPage } from '@/widgets/contact-page';

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
    ? 'Контакты — фриланс и сотрудничество'
    : 'Contact — Freelance & Collaboration';
  const description = isRu
    ? 'Связаться с Casey Laviere: фриланс на PHP/Python/TypeScript, консультации, трудоустройство. Модель оплаты 50/50, ответ в течение дня.'
    : 'Reach Casey Laviere for PHP/Python/TypeScript freelance work, consulting, or full-time offers. 50/50 payment model. Replies within one business day.';
  const localizedName = isRu ? (site.nameRu ?? site.name) : site.name;
  return {
    title,
    description,
    keywords: isRu
      ? ['нанять разработчика', 'заказать разработку', 'фриланс PHP', 'фриланс Python', 'контакты разработчика', 'Telegram разработчик', 'удалённая разработка']
      : ['hire developer', 'freelance PHP', 'freelance Python', 'contact developer', 'remote developer', 'available for hire'],
    alternates: {
      canonical: `/${locale}/contact`,
      languages: { ru: '/ru/contact', en: '/en/contact', 'x-default': '/ru/contact' },
    },
    openGraph: {
      title,
      description,
      type: 'website',
      url: `${site.url}/${locale}/contact`,
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

export default async function ContactRoute({ params }: { params: Promise<{ locale: string }> }) {
  const { locale } = await params;
  setRequestLocale(locale);

  const { homepage, contact: contactConfig } = await contentService.getHomepageAndContact();

  const isRu = locale === 'ru';
  const breadcrumb = buildBreadcrumbList(locale, [
    { name: isRu ? 'Главная' : 'Home', path: '' },
    { name: isRu ? 'Контакты' : 'Contact', path: '/contact' },
  ]);

  const contactPageJsonLd = {
    '@context': 'https://schema.org',
    '@type': 'ContactPage',
    name: isRu ? 'Контакты' : 'Contact',
    url: `${site.url}/${locale}/contact`,
    inLanguage: isRu ? 'ru-RU' : 'en-US',
    mainEntity: {
      '@type': 'Person',
      name: site.name,
      email: `mailto:${site.email}`,
      url: site.url,
      sameAs: homepage.socialLinks
        ?.filter((l) => l.platform !== 'Email')
        .map((l) => l.url) || [],
    },
  };

  return (
    <>
      <JsonLd data={breadcrumb} />
      <JsonLd data={contactPageJsonLd} />
      <ContactPage
        config={contactConfig}
        socialLinks={homepage.socialLinks}
      />
    </>
  );
}
