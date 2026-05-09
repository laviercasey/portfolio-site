'use client';

import { useEffect, useState } from 'react';
import { useTranslations, useLocale } from 'next-intl';
import {
  motion,
  useReducedMotion,
  useMotionValue,
  useSpring,
  useTransform,
} from 'framer-motion';
import { ArrowRight, Globe, Mail } from 'lucide-react';
import { Link } from '@/shared/config';
import { Button, MagneticButton, SplitText } from '@/shared/ui';
import { socialIconMap, isSafeUrl, safeHref } from '@/shared/lib';
import type { HomepageContent } from '@/entities/content';

interface HeroSectionProps {
  content: HomepageContent;
}

export default function HeroSection({ content }: HeroSectionProps) {
  const t = useTranslations('hero');
  const locale = useLocale();
  const shouldReduceMotion = useReducedMotion();

  const titles = locale === 'ru' ? content.hero.titleAnimated.ru : content.hero.titleAnimated.en;
  const [titleIndex, setTitleIndex] = useState(0);
  const [displayed, setDisplayed] = useState('');
  const [typing, setTyping] = useState(true);

  useEffect(() => {
    if (shouldReduceMotion) { setDisplayed(titles[titleIndex]); return; }
    const current = titles[titleIndex];
    if (typing) {
      if (displayed.length < current.length) {
        const t = setTimeout(() => setDisplayed(current.slice(0, displayed.length + 1)), 65);
        return () => clearTimeout(t);
      } else {
        const t = setTimeout(() => setTyping(false), 2400);
        return () => clearTimeout(t);
      }
    } else {
      if (displayed.length > 0) {
        const t = setTimeout(() => setDisplayed(displayed.slice(0, -1)), 32);
        return () => clearTimeout(t);
      } else {
        setTitleIndex((i) => (i + 1) % titles.length);
        setTyping(true);
      }
    }
  }, [displayed, typing, titleIndex, titles, shouldReduceMotion]);

  const mouseX = useMotionValue(0);
  const mouseY = useMotionValue(0);

  const orb1x = useSpring(useTransform(mouseX, [-0.5, 0.5], [-40, 40]), { stiffness: 60, damping: 30 });
  const orb1y = useSpring(useTransform(mouseY, [-0.5, 0.5], [-25, 25]), { stiffness: 60, damping: 30 });
  const orb2x = useSpring(useTransform(mouseX, [-0.5, 0.5], [30, -30]), { stiffness: 40, damping: 25 });
  const orb2y = useSpring(useTransform(mouseY, [-0.5, 0.5], [20, -20]), { stiffness: 40, damping: 25 });
  const orb3x = useSpring(useTransform(mouseX, [-0.5, 0.5], [-20, 20]), { stiffness: 30, damping: 20 });
  const orb3y = useSpring(useTransform(mouseY, [-0.5, 0.5], [15, -15]), { stiffness: 30, damping: 20 });

  useEffect(() => {
    if (shouldReduceMotion) return;
    const handler = (e: MouseEvent) => {
      mouseX.set((e.clientX / window.innerWidth) - 0.5);
      mouseY.set((e.clientY / window.innerHeight) - 0.5);
    };
    window.addEventListener('mousemove', handler, { passive: true });
    return () => window.removeEventListener('mousemove', handler);
  }, [mouseX, mouseY, shouldReduceMotion]);

  const subtitle = locale === 'ru' ? content.hero.subtitleRu : content.hero.subtitleEn;

  const stagger = {
    hidden: {},
    visible: { transition: { staggerChildren: shouldReduceMotion ? 0 : 0.1, delayChildren: 0.3 } },
  };
  const fadeUp = {
    hidden: { opacity: 1, y: shouldReduceMotion ? 0 : 16 },
    visible: { opacity: 1, y: 0, transition: { duration: 0.6 } },
  };

  return (
    <section className="relative min-h-screen flex items-center overflow-hidden pt-24">
      {!shouldReduceMotion && (
        <>
          <motion.div
            className="pointer-events-none absolute rounded-full hidden md:block"
            style={{
              x: orb1x, y: orb1y,
              width: 700, height: 700,
              top: -200, right: -150,
              background: 'radial-gradient(circle, rgba(212,168,67,0.14) 0%, transparent 65%)',
              filter: 'blur(60px)',
            }}
          />
          <motion.div
            className="pointer-events-none absolute rounded-full hidden md:block"
            style={{
              x: orb2x, y: orb2y,
              width: 500, height: 500,
              bottom: 50, left: -100,
              background: 'radial-gradient(circle, rgba(180,90,55,0.12) 0%, transparent 65%)',
              filter: 'blur(70px)',
            }}
          />
          <motion.div
            className="pointer-events-none absolute rounded-full hidden md:block"
            style={{
              x: orb3x, y: orb3y,
              width: 300, height: 300,
              top: '45%', left: '50%',
              background: 'radial-gradient(circle, rgba(212,168,67,0.07) 0%, transparent 65%)',
              filter: 'blur(50px)',
            }}
          />
        </>
      )}

      <div
        className="absolute inset-0 -z-10"
        style={{
          backgroundImage:
            'linear-gradient(rgba(255,220,160,0.04) 1px, transparent 1px), linear-gradient(90deg, rgba(255,220,160,0.04) 1px, transparent 1px)',
          backgroundSize: '80px 80px',
        }}
      />

      <div className="container py-16 relative z-10">
        <motion.div
          className="max-w-4xl"
          variants={stagger}
          initial="hidden"
          animate="visible"
        >
          {content.hero.availableForWork && (
            <motion.div variants={fadeUp} className="mb-8">
              <span className="inline-flex items-center gap-2.5 glass-card px-4 py-2 text-sm font-mono text-foreground/80">
                <span
                  className="w-2 h-2 rounded-full bg-emerald-400"
                  style={{ animation: 'pulse-dot 1.8s ease-in-out infinite' }}
                />
                {t('availableBadge')}
              </span>
            </motion.div>
          )}

          <div className="mb-4 overflow-hidden">
            <SplitText
              text={locale === 'ru' ? content.hero.name.ru : content.hero.name.en}
              as="h1"
              delay={0.35}
              className="font-heading text-[clamp(3.5rem,9vw,7rem)] font-black leading-none tracking-tight gradient-text"
            />
          </div>

          <motion.div variants={fadeUp} className="mb-7 h-10 md:h-12 flex items-center">
            <span className="font-heading italic text-2xl md:text-3xl text-foreground/70">
              {displayed}
              <span
                className="ml-1 inline-block w-0.5 h-6 bg-primary align-middle"
                style={{ animation: 'pulse-dot 0.9s step-end infinite' }}
              />
            </span>
          </motion.div>

          <motion.p
            variants={fadeUp}
            className="text-lg md:text-xl text-muted-foreground max-w-lg leading-relaxed mb-10"
          >
            {subtitle}
          </motion.p>

          <motion.div variants={fadeUp} className="flex flex-wrap gap-4 mb-12">
            <MagneticButton>
              <Button
                asChild
                size="lg"
                className="gap-2 font-medium shadow-lg"
                style={{ boxShadow: '0 4px 24px var(--glow-primary)' }}
              >
                <Link href="/projects">
                  {t('ctaProjects')}
                  <ArrowRight className="h-4 w-4" />
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
                <Link href="/contact">
                  <Mail className="h-4 w-4" />
                  {t('ctaContact')}
                </Link>
              </Button>
            </MagneticButton>
          </motion.div>

          <motion.div variants={fadeUp} className="flex items-center gap-1">
            <span className="text-xs font-mono text-muted-foreground/50 mr-3 uppercase tracking-widest">
              {t('socialLabel')}
            </span>
            {content.socialLinks.filter((l) => isSafeUrl(l.url)).map((link) => {
              const Icon = socialIconMap[link.icon] ?? Globe;
              const isExternal = link.url.startsWith('http');
              return (
                <a
                  key={link.platform}
                  href={safeHref(link.url)}
                  target={isExternal ? '_blank' : undefined}
                  rel={isExternal ? 'noopener noreferrer' : undefined}
                  aria-label={link.platform}
                  className="flex items-center justify-center w-9 h-9 rounded-xl text-muted-foreground hover:text-primary border border-transparent hover:border-[var(--glass-border)] hover:bg-primary/8 transition-all duration-200"
                >
                  <Icon className="h-4 w-4" />
                </a>
              );
            })}
          </motion.div>
        </motion.div>

        {!shouldReduceMotion && (
          <motion.div
            className="absolute right-8 bottom-12 hidden lg:flex flex-col items-center gap-2 text-muted-foreground/30"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 1.8, duration: 0.8 }}
            style={{ writingMode: 'vertical-rl' }}
          >
            <span className="text-xs font-mono uppercase tracking-[0.3em]">{t('scrollLabel')}</span>
            <div
              className="w-px h-16 bg-gradient-to-b from-muted-foreground/20 to-transparent"
              style={{ animation: 'bounce-arrow 2.2s ease-in-out infinite' }}
            />
          </motion.div>
        )}
      </div>

      <div className="absolute bottom-0 left-0 right-0 h-32 bg-gradient-to-t from-background to-transparent pointer-events-none" />
    </section>
  );
}
