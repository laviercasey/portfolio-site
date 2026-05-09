'use client';

import { useRef } from 'react';
import { motion, useScroll, useTransform } from 'framer-motion';
import { useLocale } from 'next-intl';
import type { ServiceProcessStep } from '@/entities/service';

interface ServicesProcessProps {
  steps: ServiceProcessStep[];
}

export default function ServicesProcess({ steps }: ServicesProcessProps) {
  const locale = useLocale();
  const isRu = locale === 'ru';
  const ref = useRef<HTMLDivElement>(null);

  const { scrollYProgress } = useScroll({
    target: ref,
    offset: ['start 0.85', 'end 0.4'],
  });
  const lineWidth = useTransform(scrollYProgress, [0, 1], ['0%', '100%']);

  if (steps.length === 0) return null;

  return (
    <section className="py-24 relative">
      <div className="container">
        <header className="mb-20 max-w-2xl">
          <div className="flex items-center gap-3 mb-5 text-xs font-mono uppercase tracking-[0.3em] text-muted-foreground">
            <span className="w-8 h-px bg-primary/60" />
            <span className="text-primary">PROCESS</span>
          </div>
          <h2 className="font-heading text-4xl md:text-6xl font-black leading-[1.05] mb-5">
            {isRu ? (
              <>
                Как я <span className="italic gradient-text">работаю</span>
              </>
            ) : (
              <>
                How I <span className="italic gradient-text">work</span>
              </>
            )}
          </h2>
          <p className="text-foreground/70 text-lg leading-relaxed">
            {isRu
              ? 'Прозрачный процесс без сюрпризов: от первого созвона до сдачи и поддержки.'
              : 'A transparent process with no surprises: from the first call to delivery and support.'}
          </p>
        </header>

        <div ref={ref} className="relative">
          <div className="absolute top-[18px] left-0 right-0 h-px bg-white/10 hidden md:block" />
          <motion.div
            className="absolute top-[18px] left-0 h-px bg-gradient-to-r from-primary/80 to-primary/30 hidden md:block"
            style={{ width: lineWidth }}
          />

          <div
            className="grid grid-cols-1 md:gap-6 gap-10 md:[grid-template-columns:var(--cols)]"
            style={{ ['--cols' as string]: `repeat(${steps.length}, minmax(0, 1fr))` }}
          >
            {steps.map((step, i) => (
              <motion.div
                key={step.id}
                initial={{ opacity: 0, y: 20 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true, margin: '-50px' }}
                transition={{ duration: 0.55, delay: i * 0.12, ease: [0.22, 1, 0.36, 1] }}
                className="relative"
              >
                <div className="flex md:block items-center gap-4 mb-4 md:mb-5">
                  <div
                    className="w-9 h-9 rounded-full flex items-center justify-center text-xs font-mono font-bold relative z-10 shrink-0"
                    style={{
                      background: 'rgba(212,168,67,0.12)',
                      border: '1px solid rgba(212,168,67,0.5)',
                      color: '#d4a843',
                    }}
                  >
                    {step.num}
                  </div>
                  <h3 className="font-heading text-xl md:text-2xl font-bold leading-tight md:mt-5">
                    {step.title[isRu ? 'ru' : 'en']}
                  </h3>
                </div>
                <p className="text-foreground/65 leading-relaxed text-[15px]">
                  {step.description[isRu ? 'ru' : 'en']}
                </p>
              </motion.div>
            ))}
          </div>
        </div>
      </div>
    </section>
  );
}
