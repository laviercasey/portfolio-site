'use client';

import { useEffect, useState, useRef, useCallback } from 'react';
import { useLocale } from 'next-intl';
import {
  motion,
  AnimatePresence,
  useReducedMotion,
  useInView,
} from 'framer-motion';
import {
  ClipboardList,
  Phone,
  Code2,
  PackageCheck,
  Wallet,
} from 'lucide-react';
import { ClipReveal } from '@/shared/ui';
import type { HomepageContent } from '@/entities/content';
import type { LucideIcon } from 'lucide-react';

const STEP_DURATION = 6000;
const PAUSE_DURATION = 30000;

function StepMedia({
  gifUrl, isActive, alt, icon: Icon, bg, color, size = 'desktop',
}: {
  gifUrl?: string; isActive: boolean; alt: string;
  icon: LucideIcon; bg: string; color: string;
  size?: 'desktop' | 'mobile';
}) {
  const [firstFrame, setFirstFrame] = useState<string | null>(null);

  useEffect(() => {
    if (!gifUrl) return;
    const img = new Image();
    img.crossOrigin = 'anonymous';
    img.onload = () => {
      const canvas = document.createElement('canvas');
      canvas.width = img.naturalWidth;
      canvas.height = img.naturalHeight;
      const ctx = canvas.getContext('2d');
      if (ctx) {
        ctx.drawImage(img, 0, 0);
        try { setFirstFrame(canvas.toDataURL('image/webp', 0.85)); } catch {}
      }
    };
    img.src = gifUrl;
  }, [gifUrl]);

  if (isActive && gifUrl) {
    return <img key={`active-${gifUrl}`} src={gifUrl} alt={alt} className="w-full h-full object-cover" />;
  }

  if (!isActive && firstFrame) {
    return <img src={firstFrame} alt={alt} className="w-full h-full object-cover" />;
  }

  const iconSize = size === 'desktop' ? 'h-8 w-8' : 'h-6 w-6';
  const pad = size === 'desktop' ? 'p-3' : 'p-2.5';
  return (
    <div className="w-full h-full flex flex-col items-center justify-center gap-3 bg-gradient-to-br from-white/[0.03] to-transparent">
      <div className={`${pad} rounded-xl ${bg}`}>
        <Icon className={`${iconSize} ${color}`} strokeWidth={1.5} />
      </div>
    </div>
  );
}

const stepMeta = [
  { icon: ClipboardList, color: 'text-amber-400',   bg: 'bg-amber-500/15',   ring: 'ring-amber-500/30' },
  { icon: Phone,         color: 'text-sky-400',     bg: 'bg-sky-500/15',     ring: 'ring-sky-500/30' },
  { icon: Code2,         color: 'text-emerald-400', bg: 'bg-emerald-500/15', ring: 'ring-emerald-500/30' },
  { icon: PackageCheck,  color: 'text-primary',     bg: 'bg-primary/15',     ring: 'ring-primary/30' },
];

interface Props {
  content: HomepageContent;
}

export default function HowIWorkSection({ content }: Props) {
  const locale = useLocale();
  const isRu = locale === 'ru';
  const shouldReduceMotion = useReducedMotion();
  const sectionRef = useRef<HTMLElement>(null);
  const isInView = useInView(sectionRef, { once: false, margin: '-100px' });

  const [activeStep, setActiveStep] = useState(0);
  const [isPlaying, setIsPlaying] = useState(false);
  const [isPaused, setIsPaused] = useState(false);
  const pauseTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const hw = content.howIWork;
  const steps = hw?.steps ?? [];

  useEffect(() => {
    if (isInView && !isPlaying) { setIsPlaying(true); setActiveStep(0); }
    if (!isInView) setIsPlaying(false);
  }, [isInView, isPlaying]);

  useEffect(() => {
    if (!isPlaying || isPaused || steps.length === 0) return;
    const timer = setInterval(() => {
      setActiveStep((prev) => (prev + 1) % steps.length);
    }, STEP_DURATION);
    return () => clearInterval(timer);
  }, [isPlaying, isPaused, steps.length]);

  useEffect(() => {
    return () => {
      if (pauseTimeoutRef.current) clearTimeout(pauseTimeoutRef.current);
    };
  }, []);

  const handleStepClick = useCallback((index: number) => {
    setActiveStep(index);
    setIsPlaying(true);
    setIsPaused(true);
    if (pauseTimeoutRef.current) clearTimeout(pauseTimeoutRef.current);
    pauseTimeoutRef.current = setTimeout(() => {
      setIsPaused(false);
      pauseTimeoutRef.current = null;
    }, PAUSE_DURATION);
  }, []);

  if (steps.length === 0) return null;

  return (
    <section ref={sectionRef} className="section-py overflow-hidden">
      <div className="container">
        <div className="mb-12 md:mb-16">
          <ClipReveal direction="left">
            <span className="font-mono text-xs text-primary uppercase tracking-[0.25em] block mb-2">
              {isRu ? '— как я работаю' : '— how i work'}
            </span>
          </ClipReveal>
          <ClipReveal direction="up" delay={0.1}>
            <h2 className="font-heading text-4xl md:text-5xl font-black">
              {isRu ? 'Как я работаю' : 'How I work'}
            </h2>
          </ClipReveal>
        </div>

        <div className="relative">
          <div className="hidden md:block">
            <div className="relative mx-8 mb-8">
              <div className="h-0.5 w-full bg-white/8 rounded-full" />
              <motion.div
                className="absolute top-0 left-0 h-0.5 bg-gradient-to-r from-primary via-amber-400 to-emerald-400 rounded-full origin-left"
                initial={{ scaleX: 0 }}
                animate={{ scaleX: (activeStep + 1) / steps.length }}
                transition={{ duration: 0.6, ease: [0.16, 1, 0.3, 1] }}
              />
            </div>

            <div className="grid grid-cols-4 gap-6">
              {steps.map((step, i) => {
                const meta = stepMeta[i % stepMeta.length];
                const Icon = meta.icon;
                const isActive = i === activeStep;
                const isPast = i < activeStep;

                return (
                  <button key={i} onClick={() => handleStepClick(i)} className="text-left focus:outline-none group">
                    <motion.div
                      className={`
                        relative aspect-[4/3] rounded-2xl overflow-hidden mb-5 border transition-all duration-500
                        ${isActive
                          ? `border-white/20 shadow-lg shadow-primary/10 ${meta.ring} ring-1`
                          : 'border-white/8 opacity-50 group-hover:opacity-70'
                        }
                      `}
                      animate={isActive && !shouldReduceMotion ? { scale: 1.02 } : { scale: 1 }}
                      transition={{ type: 'spring', stiffness: 200, damping: 20 }}
                    >
                      <StepMedia
                        gifUrl={step.gifUrl}
                        isActive={isActive}
                        alt={isRu ? step.titleRu : step.titleEn}
                        icon={Icon}
                        bg={meta.bg}
                        color={meta.color}
                      />

                      {isActive && (
                        <motion.div
                          className="absolute inset-0 border-2 border-primary/20 rounded-2xl"
                          initial={{ opacity: 0 }}
                          animate={{ opacity: [0, 0.5, 0] }}
                          transition={{ duration: 2, repeat: Infinity }}
                        />
                      )}
                    </motion.div>

                    <div className="flex items-center gap-3 mb-3">
                      <div className={`
                        w-3 h-3 rounded-full border-2 transition-all duration-300
                        ${isActive
                          ? 'bg-primary border-primary shadow-md shadow-primary/30 scale-125'
                          : isPast
                            ? 'bg-primary/60 border-primary/60'
                            : 'bg-transparent border-white/20 group-hover:border-white/40'
                        }
                      `} />
                      <span className={`font-mono text-xs uppercase tracking-wider transition-colors duration-300 ${isActive ? 'text-primary' : 'text-muted-foreground/50'}`}>
                        {String(i + 1).padStart(2, '0')}
                      </span>
                    </div>

                    <h3 className={`font-heading font-bold text-base mb-2 transition-colors duration-300 ${isActive ? 'text-foreground' : 'text-muted-foreground/60 group-hover:text-muted-foreground'}`}>
                      {isRu ? step.titleRu : step.titleEn}
                    </h3>

                    <AnimatePresence mode="wait">
                      {isActive && (
                        <motion.p
                          key={`desc-${i}`}
                          className="text-sm text-muted-foreground leading-relaxed"
                          initial={{ opacity: 0, y: 8 }}
                          animate={{ opacity: 1, y: 0 }}
                          exit={{ opacity: 0, y: -4 }}
                          transition={{ duration: 0.35 }}
                        >
                          {isRu ? step.descRu : step.descEn}
                        </motion.p>
                      )}
                    </AnimatePresence>

                    {isActive && isPlaying && (
                      <div className="mt-3 h-0.5 w-full bg-white/5 rounded-full overflow-hidden">
                        <motion.div
                          className={`h-full rounded-full origin-left ${isPaused ? 'bg-amber-400/70' : 'bg-primary/50'}`}
                          initial={{ scaleX: 0 }}
                          animate={{ scaleX: 1 }}
                          transition={{ duration: (isPaused ? PAUSE_DURATION : STEP_DURATION) / 1000, ease: 'linear' }}
                          key={`timer-${activeStep}-${isPaused ? 'pause' : 'play'}`}
                        />
                      </div>
                    )}
                  </button>
                );
              })}
            </div>
          </div>

          <div className="md:hidden">
            <div className="relative pl-8">
              <div className="absolute left-[11px] top-0 bottom-0 w-0.5 bg-white/8 rounded-full" />
              <motion.div
                className="absolute left-[11px] top-0 w-0.5 bg-gradient-to-b from-primary via-amber-400 to-emerald-400 rounded-full origin-top"
                animate={{ height: `${((activeStep + 1) / steps.length) * 100}%` }}
                transition={{ duration: 0.6, ease: [0.16, 1, 0.3, 1] }}
              />

              <div className="space-y-8">
                {steps.map((step, i) => {
                  const meta = stepMeta[i % stepMeta.length];
                  const Icon = meta.icon;
                  const isActive = i === activeStep;

                  return (
                    <button key={i} onClick={() => handleStepClick(i)} className="relative text-left w-full focus:outline-none group">
                      <div className={`
                        absolute -left-8 top-1 w-3 h-3 rounded-full border-2 transition-all duration-300
                        ${isActive
                          ? 'bg-primary border-primary shadow-md shadow-primary/30 scale-125'
                          : i < activeStep
                            ? 'bg-primary/60 border-primary/60'
                            : 'bg-background border-white/20'
                        }
                      `} />

                      <div className={`
                        relative aspect-video rounded-xl overflow-hidden mb-3 border transition-all duration-500
                        ${isActive ? `border-white/20 ${meta.ring} ring-1` : 'border-white/8 opacity-50'}
                      `}>
                        <StepMedia
                          gifUrl={step.gifUrl}
                          isActive={isActive}
                          alt={isRu ? step.titleRu : step.titleEn}
                          icon={Icon}
                          bg={meta.bg}
                          color={meta.color}
                          size="mobile"
                        />
                      </div>

                      <div className="flex items-center gap-2 mb-1.5">
                        <span className={`font-mono text-xs ${isActive ? 'text-primary' : 'text-muted-foreground/40'}`}>
                          {String(i + 1).padStart(2, '0')}
                        </span>
                        <h3 className={`font-heading font-bold text-sm transition-colors ${isActive ? 'text-foreground' : 'text-muted-foreground/60'}`}>
                          {isRu ? step.titleRu : step.titleEn}
                        </h3>
                      </div>

                      <AnimatePresence mode="wait">
                        {isActive && (
                          <motion.p
                            key={`m-desc-${i}`}
                            className="text-sm text-muted-foreground leading-relaxed"
                            initial={{ opacity: 0, height: 0 }}
                            animate={{ opacity: 1, height: 'auto' }}
                            exit={{ opacity: 0, height: 0 }}
                            transition={{ duration: 0.3 }}
                          >
                            {isRu ? step.descRu : step.descEn}
                          </motion.p>
                        )}
                      </AnimatePresence>

                      {isActive && isPlaying && (
                        <div className="mt-2 h-0.5 w-full bg-white/5 rounded-full overflow-hidden">
                          <motion.div
                            className={`h-full rounded-full origin-left ${isPaused ? 'bg-amber-400/70' : 'bg-primary/50'}`}
                            initial={{ scaleX: 0 }}
                            animate={{ scaleX: 1 }}
                            transition={{ duration: (isPaused ? PAUSE_DURATION : STEP_DURATION) / 1000, ease: 'linear' }}
                            key={`m-timer-${activeStep}-${isPaused ? 'pause' : 'play'}`}
                          />
                        </div>
                      )}
                    </button>
                  );
                })}
              </div>
            </div>
          </div>
        </div>

        <ClipReveal direction="up" delay={0.3} className="mt-14 md:mt-16">
          <div className="glass-card p-7 md:p-9 relative overflow-hidden">
            <div className="absolute -right-10 -top-10 w-40 h-40 rounded-full bg-primary/5 blur-3xl pointer-events-none" />
            <div className="relative">
              <h3 className="font-heading text-xl md:text-2xl font-bold mb-3">
                {isRu ? (hw?.philosophyTitleRu || 'Философия работы') : (hw?.philosophyTitleEn || 'Work Philosophy')}
              </h3>
              <p className="text-muted-foreground leading-relaxed mb-8 max-w-2xl">
                {isRu ? (hw?.philosophyTextRu || '') : (hw?.philosophyTextEn || '')}
              </p>

              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div className="flex items-start gap-4 p-4 rounded-xl bg-white/[0.03] border border-white/8 hover:border-primary/25 transition-colors">
                  <div className="p-2.5 rounded-xl bg-primary/12 shrink-0">
                    <Wallet className="h-5 w-5 text-primary" strokeWidth={1.5} />
                  </div>
                  <div>
                    <p className="font-heading font-bold text-sm mb-0.5">{isRu ? (hw?.payment1Ru || '') : (hw?.payment1En || '')}</p>
                    <p className="text-xs text-muted-foreground">{isRu ? (hw?.payment1DescRu || '') : (hw?.payment1DescEn || '')}</p>
                  </div>
                </div>

                <div className="flex items-start gap-4 p-4 rounded-xl bg-white/[0.03] border border-white/8 hover:border-emerald-500/25 transition-colors">
                  <div className="p-2.5 rounded-xl bg-emerald-500/12 shrink-0">
                    <PackageCheck className="h-5 w-5 text-emerald-400" strokeWidth={1.5} />
                  </div>
                  <div>
                    <p className="font-heading font-bold text-sm mb-0.5">{isRu ? (hw?.payment2Ru || '') : (hw?.payment2En || '')}</p>
                    <p className="text-xs text-muted-foreground">{isRu ? (hw?.payment2DescRu || '') : (hw?.payment2DescEn || '')}</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </ClipReveal>
      </div>
    </section>
  );
}
