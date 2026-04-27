import type { Metadata, Viewport } from 'next';
import { Playfair_Display, DM_Sans, JetBrains_Mono } from 'next/font/google';
import { NextIntlClientProvider } from 'next-intl';
import { getMessages, setRequestLocale } from 'next-intl/server';
import { notFound } from 'next/navigation';
import { routing } from '@/shared/config';
import { Header } from '@/widgets/header';
import { Footer } from '@/widgets/footer';
import { ThemeProvider, Umami, UtmCleanup } from '@/app/providers';
import { themeScript } from '@/shared/lib';
import { contentService } from '@/entities/content';
import { CustomCursor } from '@/shared/ui';
import { siteConfig as siteJson } from '@/shared/config';
import type { SiteConfig } from '@/shared/types';
import '../globals.css';

const playfair = Playfair_Display({
  subsets: ['latin', 'cyrillic'],
  weight: ['700', '800', '900'],
  style: ['normal', 'italic'],
  variable: '--font-playfair',
  display: 'swap',
});

const dmSans = DM_Sans({
  subsets: ['latin', 'latin-ext', 'cyrillic'] as Array<'latin' | 'latin-ext'>,
  weight: ['300', '400', '500', '600'],
  variable: '--font-dm-sans',
  display: 'swap',
});

const jetbrainsMono = JetBrains_Mono({
  subsets: ['latin'],
  weight: ['400', '500'],
  variable: '--font-jetbrains-mono',
  display: 'swap',
});

const site = siteJson as SiteConfig;

export const viewport: Viewport = {
  themeColor: '#1a120c',
  width: 'device-width',
  initialScale: 1,
};

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string }>;
}): Promise<Metadata> {
  const { locale } = await params;
  const isRu = locale === 'ru';

  const localizedName = isRu ? (site.nameRu ?? site.name) : site.name;
  const title = isRu
    ? `${localizedName} — ${site.taglineRu}`
    : `${localizedName} — ${site.taglineEn}`;
  const description = isRu ? site.descriptionRu : site.descriptionEn;

  return {
    title: {
      default: title,
      template: `%s | ${localizedName}`,
    },
    description,
    keywords: isRu ? site.keywordsRu : site.keywordsEn,
    authors: [{ name: localizedName, url: site.url }],
    creator: localizedName,
    metadataBase: new URL(site.url),
    alternates: {
      canonical: `/${locale}`,
      languages: {
        ru: '/ru',
        en: '/en',
        'x-default': '/ru',
      },
    },
    openGraph: {
      title,
      description,
      url: `${site.url}/${locale}`,
      siteName: site.name,
      locale: isRu ? 'ru_RU' : 'en_US',
      type: 'website',
    },
    twitter: {
      card: 'summary_large_image',
      title,
      description,
    },
    verification: {
      ...(process.env.GOOGLE_VERIFICATION && { google: process.env.GOOGLE_VERIFICATION }),
      ...(process.env.YANDEX_VERIFICATION && { yandex: process.env.YANDEX_VERIFICATION }),
    },
    robots: {
      index: true,
      follow: true,
      googleBot: {
        index: true,
        follow: true,
        'max-video-preview': -1,
        'max-image-preview': 'large',
        'max-snippet': -1,
      },
    },
    manifest: '/site.webmanifest',
  };
}

export default async function LocaleLayout({
  children,
  params,
}: {
  children: React.ReactNode;
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;

  if (!routing.locales.includes(locale as 'ru' | 'en')) {
    notFound();
  }

  setRequestLocale(locale);

  const [messages, content] = await Promise.all([getMessages(), contentService.getHomepage()]);

  const isRu = locale === 'ru';
  const localizedName = isRu ? (site.nameRu ?? site.name) : site.name;

  const personJsonLd = {
    '@context': 'https://schema.org',
    '@type': 'Person',
    name: localizedName,
    url: site.url,
    jobTitle: isRu ? site.taglineRu : site.taglineEn,
    description: isRu ? site.descriptionRu : site.descriptionEn,
    email: `mailto:${site.email}`,
    image: `${site.url}/${locale}/opengraph-image`,
    knowsAbout: ['PHP', 'Laravel', 'Python', 'TypeScript', 'React', 'Next.js', 'Machine Learning', 'Docker', 'PostgreSQL'],
    sameAs: content.socialLinks
      ?.filter((l) => l.platform !== 'Email')
      .map((l) => l.url) || [],
  };

  const websiteJsonLd = {
    '@context': 'https://schema.org',
    '@type': 'WebSite',
    name: localizedName,
    url: site.url,
    description: isRu ? site.descriptionRu : site.descriptionEn,
    inLanguage: isRu ? 'ru-RU' : 'en-US',
    author: { '@type': 'Person', name: localizedName },
  };

  return (
    <html
      lang={locale}
      suppressHydrationWarning
      className={`${playfair.variable} ${dmSans.variable} ${jetbrainsMono.variable}`}
    >
      <head>
        <script dangerouslySetInnerHTML={{ __html: themeScript }} />
        <script
          type="application/ld+json"
          dangerouslySetInnerHTML={{ __html: JSON.stringify(personJsonLd).replace(/</g, '\\u003c') }}
        />
        <script
          type="application/ld+json"
          dangerouslySetInnerHTML={{ __html: JSON.stringify(websiteJsonLd).replace(/</g, '\\u003c') }}
        />
      </head>
      <body className="font-sans antialiased">
        <Umami />
        <UtmCleanup />
        <NextIntlClientProvider locale={locale} messages={messages}>
          <ThemeProvider>
            <CustomCursor />
            <div className="min-h-screen flex flex-col">
              <Header />
              <main className="flex-1">{children}</main>
              <Footer socialLinks={content.socialLinks} />
            </div>
          </ThemeProvider>
        </NextIntlClientProvider>
      </body>
    </html>
  );
}
