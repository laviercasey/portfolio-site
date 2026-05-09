'use client';

import { useEffect } from 'react';
import { motion, useMotionValue, useSpring, useTransform, useReducedMotion } from 'framer-motion';
import { useLocale } from 'next-intl';
import { SplitText } from '@/shared/ui';
import { SERVICE_ICON_MAP } from './service-icons';
import type { Service, ServiceIconKey } from '@/entities/service';

interface ServicesHeroProps {
  services: Service[];
}

export default function ServicesHero({ services }: ServicesHeroProps) {
  const locale = useLocale();
  const isRu = locale === 'ru';
  const shouldReduceMotion = useReducedMotion();

  const mouseX = useMotionValue(0);
  const mouseY = useMotionValue(0);
  const orb1x = useSpring(useTransform(mouseX, [-0.5, 0.5], [-30, 30]), { stiffness: 60, damping: 30 });
  const orb1y = useSpring(useTransform(mouseY, [-0.5, 0.5], [-20, 20]), { stiffness: 60, damping: 30 });
  const orb2x = useSpring(useTransform(mouseX, [-0.5, 0.5], [25, -25]), { stiffness: 40, damping: 25 });
  const orb2y = useSpring(useTransform(mouseY, [-0.5, 0.5], [15, -15]), { stiffness: 40, damping: 25 });

  useEffect(() => {
    if (shouldReduceMotion) return;
    const handler = (e: MouseEvent) => {
      mouseX.set(e.clientX / window.innerWidth - 0.5);
      mouseY.set(e.clientY / window.innerHeight - 0.5);
    };
    window.addEventListener('mousemove', handler, { passive: true });
    return () => window.removeEventListener('mousemove', handler);
  }, [mouseX, mouseY, shouldReduceMotion]);

  const stagger = {
    hidden: {},
    visible: { transition: { staggerChildren: 0.08, delayChildren: 0.4 } },
  };
  const fadeUp = {
    hidden: { opacity: 1, y: shouldReduceMotion ? 0 : 12 },
    visible: { opacity: 1, y: 0, transition: { duration: 0.55 } },
  };

  const scrollTo = (id: string) => {
    const el = document.getElementById(id);
    if (el) el.scrollIntoView({ behavior: 'smooth', block: 'start' });
  };

  const total = services.length || 5;
  const totalLabel = total < 10 ? `0${total}` : String(total);

  return (
    <section className="relative min-h-[88vh] flex items-center overflow-hidden pt-32 pb-20">
      {!shouldReduceMotion && (
        <>
          <motion.div
            className="pointer-events-none absolute rounded-full hidden md:block"
            style={{
              x: orb1x, y: orb1y,
              width: 720, height: 720,
              top: -180, right: -160,
              background: 'radial-gradient(circle, rgba(212,168,67,0.16) 0%, transparent 65%)',
              filter: 'blur(70px)',
            }}
          />
          <motion.div
            className="pointer-events-none absolute rounded-full hidden md:block"
            style={{
              x: orb2x, y: orb2y,
              width: 520, height: 520,
              bottom: 40, left: -120,
              background: 'radial-gradient(circle, rgba(180,90,55,0.14) 0%, transparent 65%)',
              filter: 'blur(75px)',
            }}
          />
        </>
      )}

      <div
        className="absolute inset-0 -z-10"
        style={{
          backgroundImage:
            'linear-gradient(rgba(255,220,160,0.03) 1px, transparent 1px), linear-gradient(90deg, rgba(255,220,160,0.03) 1px, transparent 1px)',
          backgroundSize: '80px 80px',
        }}
      />

      <div className="container relative z-10">
        <motion.div
          className="max-w-5xl"
          variants={stagger}
          initial="hidden"
          animate="visible"
        >
          <motion.div variants={fadeUp} className="mb-8 flex items-center gap-3 text-xs font-mono uppercase tracking-[0.3em] text-muted-foreground">
            <span className="w-8 h-px bg-primary/60" />
            <span className="text-primary">01—{totalLabel}</span>
            <span>·</span>
            <span>{isRu ? 'УСЛУГИ' : 'SERVICES'}</span>
          </motion.div>

          <div className="mb-6 overflow-hidden">
            <SplitText
              text={isRu ? 'Разработка' : 'Custom'}
              as="h1"
              delay={0.4}
              className="font-heading block text-[clamp(3.5rem,10vw,8rem)] font-black leading-[0.95] tracking-tight gradient-text"
            />
            <SplitText
              text={isRu ? 'на заказ' : 'development'}
              as="span"
              delay={0.55}
              className="font-heading block text-[clamp(3.5rem,10vw,8rem)] font-black leading-[0.95] tracking-tight italic text-foreground/85"
            />
          </div>

          <motion.p
            variants={fadeUp}
            className="text-lg md:text-2xl text-foreground/70 max-w-2xl leading-relaxed mb-12 font-light"
          >
            {isRu
              ? 'От ботов на 2 недели до ML-моделей в продакшене. Сайты на Laravel и Next.js, парсеры, интеграции — всё на удалёнке, по договору, с фикс-ценой после ТЗ.'
              : 'From 2-week bots to production ML. Laravel and Next.js sites, parsers, integrations — remote, contract-based, fixed price once the spec is agreed.'}
          </motion.p>

          {services.length > 0 && (
            <motion.div variants={fadeUp} className="flex flex-wrap gap-2.5 max-w-3xl">
              {services.map((s) => {
                const Icon = SERVICE_ICON_MAP[s.iconKey as ServiceIconKey] ?? SERVICE_ICON_MAP.bot;
                return (
                  <button
                    key={s.id}
                    onClick={() => scrollTo(s.slug)}
                    className="group inline-flex items-center gap-2 glass-card pl-2 pr-4 py-2 text-sm hover:border-primary/40 hover:bg-white/[0.04] transition-all duration-300"
                  >
                    <span
                      className="flex items-center justify-center w-7 h-7 rounded-lg shrink-0"
                      style={{
                        background: `${s.accent}1a`,
                        border: `1px solid ${s.accent}33`,
                      }}
                    >
                      <Icon className="h-3.5 w-3.5" style={{ color: s.accent }} strokeWidth={2.2} />
                    </span>
                    <span className="font-medium text-foreground/85 group-hover:text-foreground transition-colors">
                      {s.title[isRu ? 'ru' : 'en']}
                    </span>
                    <span className="text-muted-foreground/70 text-xs font-mono">
                      {isRu ? s.priceRu : s.priceEn}
                    </span>
                  </button>
                );
              })}
            </motion.div>
          )}
        </motion.div>
      </div>

      <div className="absolute bottom-0 left-0 right-0 h-24 bg-gradient-to-t from-background to-transparent pointer-events-none" />
    </section>
  );
}
