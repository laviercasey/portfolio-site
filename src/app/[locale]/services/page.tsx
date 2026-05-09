import type { Metadata } from 'next';
import { setRequestLocale } from 'next-intl/server';
import { JsonLd, buildBreadcrumbList } from '@/shared/lib';
import { siteConfig as siteJson } from '@/shared/config';
import type { SiteConfig } from '@/shared/types';
import { servicesService } from '@/entities/service';
import {
  ServicesHero,
  ServiceBlock,
  ServicesProcess,
  ServicesFaq,
  ServicesCTA,
} from '@/widgets/services-page';

export const revalidate = 600;

const site = siteJson as SiteConfig;

interface ServicesPageProps {
  params: Promise<{ locale: string }>;
}

export async function generateMetadata({ params }: ServicesPageProps): Promise<Metadata> {
  const { locale } = await params;
  const isRu = locale === 'ru';

  const title = isRu
    ? 'Услуги — разработка сайтов, Telegram-ботов, ML-моделей на заказ'
    : 'Services — Custom websites, Telegram bots, and ML models';
  const description = isRu
    ? 'Разработка сайтов на Laravel и Next.js, Telegram-ботов на Python, ML-моделей и парсеров. Цены от 20 000 ₽, работа по договору, прозрачные сроки. Свободна для удалённых проектов.'
    : 'Custom development services: Laravel and Next.js websites, Python Telegram bots, ML models, parsers. Transparent pricing, fixed timelines, contract-based remote work.';

  return {
    title,
    description,
    keywords: isRu
      ? [
          'разработка сайта на заказ',
          'сайт под ключ',
          'разработка Telegram бота',
          'телеграм бот на python',
          'сайт на Laravel',
          'сайт на Next.js',
          'разработка веб-приложения',
          'ML модель на заказ',
          'разработка парсера',
          'фриланс разработчик',
        ]
      : [
          'custom web development',
          'Laravel development services',
          'Next.js development services',
          'Telegram bot development',
          'Python developer for hire',
          'ML model development',
          'data parser development',
          'freelance full-stack developer',
        ],
    alternates: {
      canonical: `/${locale}/services`,
      languages: {
        ru: '/ru/services',
        en: '/en/services',
        'x-default': '/en/services',
      },
    },
    openGraph: {
      title,
      description,
      type: 'website',
      url: `${site.url}/${locale}/services`,
      locale: isRu ? 'ru_RU' : 'en_US',
      siteName: isRu ? (site.nameRu ?? site.name) : site.name,
    },
    twitter: {
      card: 'summary_large_image',
      title,
      description,
    },
  };
}

export default async function ServicesPage({ params }: ServicesPageProps) {
  const { locale } = await params;
  setRequestLocale(locale);

  const isRu = locale === 'ru';
  const { services, faqs, processSteps } = await servicesService.getPageData();

  const breadcrumb = buildBreadcrumbList(locale, [
    { name: isRu ? 'Главная' : 'Home', path: '' },
    { name: isRu ? 'Услуги' : 'Services', path: '/services' },
  ]);

  const faqSchema = faqs.length > 0 ? {
    '@context': 'https://schema.org',
    '@type': 'FAQPage',
    inLanguage: isRu ? 'ru-RU' : 'en-US',
    mainEntity: faqs.map((f) => ({
      '@type': 'Question',
      name: f.question[isRu ? 'ru' : 'en'],
      acceptedAnswer: { '@type': 'Answer', text: f.answer[isRu ? 'ru' : 'en'] },
    })),
  } : null;

  const serviceSchema = services.length > 0 ? {
    '@context': 'https://schema.org',
    '@type': 'ProfessionalService',
    name: isRu
      ? `${site.nameRu ?? site.name} — разработка на заказ`
      : `${site.name} — custom development`,
    url: `${site.url}/${locale}/services`,
    inLanguage: isRu ? 'ru-RU' : 'en-US',
    description: isRu
      ? 'Разработка сайтов на Laravel и Next.js, Telegram-ботов на Python, ML-моделей и парсеров на заказ.'
      : 'Custom Laravel and Next.js websites, Python Telegram bots, ML models, and parsers.',
    areaServed: isRu ? 'Удалённо, по всему миру' : 'Worldwide, remote',
    provider: {
      '@type': 'Person',
      name: isRu ? (site.nameRu ?? site.name) : site.name,
      url: site.url,
      email: `mailto:${site.email}`,
    },
    hasOfferCatalog: {
      '@type': 'OfferCatalog',
      name: isRu ? 'Услуги разработки' : 'Development services',
      itemListElement: services.map((s) => ({
        '@type': 'Offer',
        itemOffered: { '@type': 'Service', name: s.title[isRu ? 'ru' : 'en'], description: s.lead[isRu ? 'ru' : 'en'] },
        priceSpecification: {
          '@type': 'PriceSpecification',
          price: isRu ? s.priceRu : s.priceEn,
          priceCurrency: isRu ? 'RUB' : 'USD',
        },
      })),
    },
  } : null;

  return (
    <>
      <JsonLd data={breadcrumb} />
      {faqSchema && <JsonLd data={faqSchema} />}
      {serviceSchema && <JsonLd data={serviceSchema} />}

      <ServicesHero services={services} />

      <div className="container">
        {services.map((service, i) => (
          <ServiceBlock key={service.id} service={service} index={i} />
        ))}
      </div>

      <ServicesProcess steps={processSteps} />
      <ServicesFaq faqs={faqs} />
      <ServicesCTA />
    </>
  );
}
