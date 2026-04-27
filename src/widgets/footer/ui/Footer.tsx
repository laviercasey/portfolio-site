'use client';

import { useTranslations, useLocale } from 'next-intl';
import { siteConfig as siteJson } from '@/shared/config';
import type { SiteConfig } from '@/shared/types';
import { Globe } from 'lucide-react';
import { socialIconMap, isSafeUrl, safeHref } from '@/shared/lib';
import type { SocialLink } from '@/entities/content';

const site = siteJson as SiteConfig;

interface FooterProps {
  socialLinks: SocialLink[];
}

export default function Footer({ socialLinks }: FooterProps) {
  const t = useTranslations('footer');
  const locale = useLocale();
  const siteName = locale === 'ru' ? (site.nameRu ?? site.name) : site.name;
  const year = new Date().getFullYear();

  return (
    <footer className="relative mt-auto">
      <div className="h-px w-full bg-gradient-to-r from-transparent via-primary/40 to-transparent" />

      <div className="container py-8">
        <div className="flex flex-col md:flex-row items-center justify-between gap-5">
          <div className="flex items-center gap-3">
            <span className="font-heading italic font-black text-lg gradient-text">CL</span>
            <span className="text-sm text-muted-foreground font-mono">
              &copy; {year} {siteName}
            </span>
          </div>

          <p className="text-xs text-muted-foreground/50 font-mono order-last md:order-none">
            {t('rights')}
          </p>

          <div className="flex items-center gap-1.5">
            {socialLinks.filter((l) => isSafeUrl(l.url)).map((link) => {
              const Icon = socialIconMap[link.icon] ?? Globe;
              const isExternal = link.url.startsWith('http');
              return (
                <a
                  key={link.platform}
                  href={safeHref(link.url)}
                  target={isExternal ? '_blank' : undefined}
                  rel={isExternal ? 'noopener noreferrer' : undefined}
                  aria-label={link.platform}
                  className="flex items-center justify-center w-8 h-8 rounded-xl text-muted-foreground hover:text-primary border border-transparent hover:border-[var(--glass-border)] hover:bg-primary/8 transition-all duration-200"
                >
                  <Icon className="h-3.5 w-3.5" />
                </a>
              );
            })}
          </div>
        </div>
      </div>
    </footer>
  );
}
