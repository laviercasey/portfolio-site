'use client';

import { useLocale, useTranslations } from 'next-intl';
import { motion, useReducedMotion } from 'framer-motion';
import { Wrench } from 'lucide-react';
import type { TechChoice } from '@/entities/project';

interface Props {
  choices: TechChoice[];
}

export default function ProjectTechChoices({ choices }: Props) {
  const locale = useLocale();
  const t = useTranslations('projects');
  const shouldReduce = useReducedMotion();
  const isRu = locale === 'ru';

  if (!choices || choices.length === 0) return null;

  return (
    <section className="glass-card p-6 rounded-xl">
      <h2 className="font-heading text-lg lg:text-xl font-bold mb-5 flex items-center gap-2">
        <Wrench className="h-4 w-4 text-primary" />
        {t('techChoicesTitle')}
      </h2>
      <ul className="space-y-3">
        {choices.map((c, i) => {
          const reason = c.reason && (isRu ? c.reason.ru : c.reason.en);
          return (
            <motion.li
              key={`${i}-${c.tech}`}
              className="flex flex-col gap-1.5 sm:flex-row sm:items-start sm:gap-3"
              initial={shouldReduce ? {} : { opacity: 0, x: -8 }}
              whileInView={{ opacity: 1, x: 0 }}
              viewport={{ once: true, margin: '-40px' }}
              transition={{ duration: 0.35, delay: i * 0.05 }}
            >
              <span className="font-mono text-sm md:text-base px-2 py-0.5 md:px-2.5 rounded bg-primary/10 text-primary border border-primary/20 self-start break-all sm:shrink-0 sm:break-normal">
                {c.tech}
              </span>
              {reason && (
                <span className="text-sm md:text-base lg:text-lg text-muted-foreground leading-relaxed break-words min-w-0">
                  <span className="hidden sm:inline">— </span>{reason}
                </span>
              )}
            </motion.li>
          );
        })}
      </ul>
    </section>
  );
}
