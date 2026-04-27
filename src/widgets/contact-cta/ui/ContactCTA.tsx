'use client';

import { useLocale } from 'next-intl';
import { useTranslations } from 'next-intl';
import { Link } from '@/shared/config';
import { Mail, ArrowRight } from 'lucide-react';
import { motion, useReducedMotion } from 'framer-motion';
import { Button, MagneticButton, SplitText } from '@/shared/ui';
import type { HomepageContent } from '@/entities/content';

interface ContactCTAProps {
  content: HomepageContent;
}

export default function ContactCTA({ content }: ContactCTAProps) {
  const t = useTranslations('hero');
  const locale = useLocale();
  const shouldReduceMotion = useReducedMotion();
  const isRu = locale === 'ru';

  const cta = content.contactCTA;
  const heading = isRu ? (cta?.headingRu || 'Давайте создадим что-то вместе') : (cta?.headingEn || "Let's build something together");
  const subtitle = isRu ? (cta?.subtitleRu || 'Открыта к фриланс-проектам, сотрудничеству и интересным предложениям.') : (cta?.subtitleEn || 'Open to freelance projects, collaborations, and interesting opportunities.');
  const bgWord = isRu ? (cta?.bgWordRu || 'ВМЕСТЕ') : (cta?.bgWordEn || 'TOGETHER');

  return (
    <section className="section-py">
      <div className="container">
        <motion.div
          className="relative overflow-hidden rounded-3xl glass-card p-10 md:p-20 text-center"
          initial={{ opacity: 0, y: shouldReduceMotion ? 0 : 40 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true, margin: '-80px' }}
          transition={{ duration: 0.7 }}
          style={{
            background: 'radial-gradient(ellipse at 50% 50%, rgba(212,168,67,0.08) 0%, rgba(180,90,55,0.05) 50%, transparent 80%)',
          }}
        >
          <span
            className="absolute inset-0 flex items-center justify-center font-heading font-black text-[clamp(4rem,15vw,12rem)] leading-none select-none pointer-events-none opacity-[0.025] whitespace-nowrap"
          >
            {bgWord}
          </span>

          <div className="absolute top-0 left-1/2 -translate-x-1/2 w-32 h-0.5 bg-gradient-to-r from-transparent via-primary to-transparent" />

          <div className="relative">
            {content.hero.availableForWork && (
              <div className="mb-8 flex justify-center">
                <span className="inline-flex items-center gap-2.5 glass-card px-4 py-2 text-sm font-mono text-foreground/70">
                  <span
                    className="w-2 h-2 rounded-full bg-emerald-400"
                    style={{ animation: 'pulse-dot 1.8s ease-in-out infinite' }}
                  />
                  {t('availableBadge')}
                </span>
              </div>
            )}

            <div className="mb-4">
              <SplitText
                text={heading}
                as="h2"
                delay={0.1}
                className="font-heading text-4xl md:text-6xl font-black gradient-text leading-tight"
              />
            </div>

            <p className="text-muted-foreground text-lg mb-12 max-w-xl mx-auto leading-relaxed">
              {subtitle}
            </p>

            <div className="flex flex-wrap gap-4 justify-center">
              <MagneticButton>
                <Button
                  asChild
                  size="lg"
                  className="gap-2 shadow-lg font-medium"
                  style={{ boxShadow: '0 4px 24px var(--glow-primary)' }}
                >
                  <Link href="/contact">
                    <Mail className="h-4 w-4" />
                    {t('ctaContact')}
                  </Link>
                </Button>
              </MagneticButton>
              <MagneticButton>
                <Button
                  asChild
                  variant="outline"
                  size="lg"
                  className="gap-2 glass-card border-[var(--glass-border)] hover:border-primary/50 hover:bg-primary/5 font-medium"
                >
                  <Link href="/projects">
                    {isRu ? 'Смотреть проекты' : 'View projects'}
                    <ArrowRight className="h-4 w-4" />
                  </Link>
                </Button>
              </MagneticButton>
            </div>
          </div>
        </motion.div>
      </div>
    </section>
  );
}
