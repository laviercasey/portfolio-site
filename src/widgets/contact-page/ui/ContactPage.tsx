'use client';

import { motion } from 'framer-motion';
import { Mail, Clock, Globe2, Briefcase, Languages } from 'lucide-react';
import { useLocale, useTranslations } from 'next-intl';
import { ClipReveal } from '@/shared/ui';
import { socialIconMap, isSafeUrl, safeHref } from '@/shared/lib';
import type { ContactPageConfig, SocialLink } from '@/entities/content';

const PRIMARY_ICONS = new Set(['mail', 'send']);

function formatHandle(link: SocialLink): string {
  if (link.url.startsWith('mailto:')) return link.url.slice('mailto:'.length);
  if (link.url.startsWith('tel:')) return link.url.slice('tel:'.length);
  const tmeMatch = link.url.match(/^https?:\/\/t\.me\/([^/?#]+)/i);
  if (tmeMatch) return `@${tmeMatch[1]}`;
  return link.url.replace(/^https?:\/\//, '');
}

interface Props {
  config: ContactPageConfig;
  socialLinks: SocialLink[];
}

const stagger = {
  hidden: {},
  visible: { transition: { staggerChildren: 0.08 } },
};

const fadeUp = {
  hidden: { opacity: 0, y: 24 },
  visible: { opacity: 1, y: 0, transition: { duration: 0.5, ease: [0.16, 1, 0.3, 1] as [number, number, number, number] } },
};

interface InfoRowProps {
  icon: React.ComponentType<{ className?: string }>;
  label: string;
  value: string;
}

function InfoRow({ icon: Icon, label, value }: InfoRowProps) {
  return (
    <div className="flex items-start gap-3">
      <span className="mt-0.5 p-2 rounded-lg bg-primary/10 text-primary flex-shrink-0">
        <Icon className="h-4 w-4" />
      </span>
      <div className="min-w-0">
        <p className="font-mono text-[11px] uppercase tracking-[0.18em] text-muted-foreground/70 mb-1">
          {label}
        </p>
        <p className="text-sm text-foreground/90 leading-snug">{value}</p>
      </div>
    </div>
  );
}

export default function ContactPage({ config, socialLinks }: Props) {
  const locale = useLocale();
  const isRu = locale === 'ru';
  const t = useTranslations('contact');

  const safeLinks = socialLinks.filter((l) => isSafeUrl(l.url));
  const primaryContacts = safeLinks.filter((l) => PRIMARY_ICONS.has(l.icon));
  const secondaryContacts = safeLinks.filter((l) => !PRIMARY_ICONS.has(l.icon));

  return (
    <div className="container py-24 max-w-6xl">
      <ClipReveal direction="up" className="mb-14">
        <h1 className="font-heading text-5xl md:text-6xl font-bold mb-3 gradient-text">
          {isRu ? config.heading.ru : config.heading.en}
        </h1>
        <p className="text-lg text-muted-foreground max-w-xl">
          {isRu ? config.subtitle.ru : config.subtitle.en}
        </p>
      </ClipReveal>

      <motion.div
        variants={stagger}
        initial="hidden"
        animate="visible"
        className="grid grid-cols-1 lg:grid-cols-12 gap-5 items-start"
      >
        <motion.div variants={fadeUp} className="lg:col-span-7 space-y-5">
          <div className="glass-card p-6">
            <div className="flex items-center gap-2 mb-4">
              <div className="h-1.5 w-1.5 rounded-full bg-primary animate-pulse" />
              <h2 className="font-heading font-semibold text-sm uppercase tracking-wider text-muted-foreground">
                {t('howIWork')}
              </h2>
            </div>
            <p className="text-sm text-muted-foreground leading-relaxed">
              {isRu ? config.howIWork.ru : config.howIWork.en}
            </p>

            {primaryContacts.length > 0 && (
              <div className="mt-5 pt-5 border-t border-white/8">
                <p className="font-mono text-[11px] uppercase tracking-[0.2em] text-muted-foreground/70 mb-3">
                  {t('quickContacts')}
                </p>
                <div className="flex flex-col gap-2">
                  {primaryContacts.map((link) => {
                    const Icon = socialIconMap[link.icon] ?? Mail;
                    return (
                      <a
                        key={link.platform}
                        href={safeHref(link.url)}
                        {...(link.url.startsWith('http') ? { target: '_blank', rel: 'noopener noreferrer' } : {})}
                        className="flex items-center gap-3 px-3 py-2.5 rounded-lg bg-white/[0.03] border border-white/8 hover:border-primary/30 hover:bg-primary/5 transition-all group"
                      >
                        <span className="p-1.5 rounded-md bg-primary/10 text-primary group-hover:bg-primary/20 transition-colors">
                          <Icon className="h-3.5 w-3.5" />
                        </span>
                        <span className="font-mono text-sm text-foreground/90 truncate">
                          {formatHandle(link)}
                        </span>
                      </a>
                    );
                  })}
                </div>
              </div>
            )}
          </div>

          {secondaryContacts.length > 0 && (
            <div className="grid grid-cols-2 sm:grid-cols-3 gap-3">
              {secondaryContacts.map((link, i) => {
                const Icon = socialIconMap[link.icon] ?? Mail;
                return (
                  <motion.a
                    key={link.platform}
                    href={safeHref(link.url)}
                    {...(link.url.startsWith('http') ? { target: '_blank', rel: 'noopener noreferrer' } : {})}
                    className="glass-card p-4 flex flex-col gap-3 group hover:border-primary/30 hover:bg-primary/5 transition-all duration-300"
                    whileHover={{ y: -4 }}
                    initial={{ opacity: 0, scale: 0.9 }}
                    animate={{ opacity: 1, scale: 1 }}
                    transition={{ delay: 0.15 + i * 0.05, duration: 0.4 }}
                  >
                    <div className="p-2 rounded-lg bg-primary/10 text-primary w-fit group-hover:bg-primary/20 transition-colors">
                      <Icon className="h-4 w-4" />
                    </div>
                    <div>
                      <p className="font-medium text-sm">{link.platform}</p>
                      <p className="text-xs text-muted-foreground truncate font-mono">
                        {link.url.replace(/^https?:\/\//, '').replace(/^mailto:/, '')}
                      </p>
                    </div>
                  </motion.a>
                );
              })}
            </div>
          )}
        </motion.div>

        <motion.div variants={fadeUp} className="lg:col-span-5">
          <div className="glass-card p-6">
            <h2 className="font-heading font-semibold text-sm uppercase tracking-wider text-muted-foreground mb-5">
              {t('atGlance')}
            </h2>
            <div className="space-y-5">
              <InfoRow
                icon={Clock}
                label={t('info.responseLabel')}
                value={t('info.responseValue')}
              />
              <InfoRow
                icon={Globe2}
                label={t('info.timezoneLabel')}
                value={t('info.timezoneValue')}
              />
              <InfoRow
                icon={Briefcase}
                label={t('info.openToLabel')}
                value={t('info.openToValue')}
              />
              <InfoRow
                icon={Languages}
                label={t('info.languagesLabel')}
                value={t('info.languagesValue')}
              />
            </div>
          </div>
        </motion.div>
      </motion.div>
    </div>
  );
}
