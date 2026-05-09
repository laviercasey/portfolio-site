'use client';

import { motion, useReducedMotion } from 'framer-motion';
import { useLocale } from 'next-intl';
import { Mail, ArrowRight } from 'lucide-react';
import { Link } from '@/shared/config';
import { Button, MagneticButton } from '@/shared/ui';

export default function ServicesCTA() {
  const locale = useLocale();
  const isRu = locale === 'ru';
  const shouldReduceMotion = useReducedMotion();

  return (
    <section className="relative py-32 overflow-hidden">
      {!shouldReduceMotion && (
        <>
          <motion.div
            className="absolute pointer-events-none rounded-full"
            style={{
              width: 500, height: 500,
              top: '-25%', left: '50%', transform: 'translateX(-50%)',
              background: 'radial-gradient(circle, rgba(212,168,67,0.18) 0%, transparent 65%)',
              filter: 'blur(70px)',
            }}
            animate={{ scale: [1, 1.08, 1] }}
            transition={{ duration: 8, repeat: Infinity, ease: 'easeInOut' }}
          />
          <motion.div
            className="absolute pointer-events-none rounded-full"
            style={{
              width: 380, height: 380,
              bottom: '-20%', left: '15%',
              background: 'radial-gradient(circle, rgba(180,90,55,0.13) 0%, transparent 65%)',
              filter: 'blur(60px)',
            }}
            animate={{ scale: [1, 1.12, 1] }}
            transition={{ duration: 10, repeat: Infinity, ease: 'easeInOut' }}
          />
        </>
      )}

      <div className="container max-w-3xl text-center relative">
        <motion.div
          initial={{ opacity: 0, y: 24 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true, margin: '-50px' }}
          transition={{ duration: 0.7, ease: [0.22, 1, 0.36, 1] }}
        >
          <div className="flex items-center justify-center gap-3 mb-8 text-xs font-mono uppercase tracking-[0.3em] text-muted-foreground">
            <span className="w-8 h-px bg-primary/60" />
            <span className="text-primary">{isRu ? 'НАЧНЁМ' : "LET'S START"}</span>
            <span className="w-8 h-px bg-primary/60" />
          </div>

          <h2 className="font-heading text-5xl md:text-7xl font-black mb-6 leading-[1.05]">
            {isRu ? (
              <>
                Обсудим
                <br />
                <span className="italic gradient-text">ваш проект?</span>
              </>
            ) : (
              <>
                Want to discuss
                <br />
                <span className="italic gradient-text">your project?</span>
              </>
            )}
          </h2>

          <p className="text-foreground/70 text-lg md:text-xl leading-relaxed mb-12 max-w-xl mx-auto">
            {isRu
              ? 'Бесплатный созвон 20–30 минут с оценкой задачи, без обязательств. Отвечаю в течение дня.'
              : 'Free 20–30 min intro call with a rough estimate, no commitment. I reply within the day.'}
          </p>

          <div className="flex flex-wrap items-center justify-center gap-4">
            <MagneticButton>
              <Button asChild size="lg" className="gap-2 group shadow-lg" style={{ boxShadow: '0 4px 24px var(--glow-primary)' }}>
                <Link href="/contact">
                  <Mail className="h-4 w-4" />
                  {isRu ? 'Написать' : 'Get in touch'}
                  <ArrowRight className="h-4 w-4 group-hover:translate-x-0.5 transition-transform" />
                </Link>
              </Button>
            </MagneticButton>
            <MagneticButton>
              <Button asChild variant="outline" size="lg" className="glass-card gap-2 border-[var(--glass-border)] hover:border-primary/50">
                <Link href="/projects">
                  {isRu ? 'Посмотреть проекты' : 'See projects'}
                </Link>
              </Button>
            </MagneticButton>
          </div>
        </motion.div>
      </div>
    </section>
  );
}
