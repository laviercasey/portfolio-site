'use client';

import { useLocale } from 'next-intl';
import { Briefcase } from 'lucide-react';
import { motion, useReducedMotion } from 'framer-motion';
import { Badge, ClipReveal } from '@/shared/ui';
import type { CareerContent } from '@/entities/career';
import { formatMonthYear } from '@/shared/lib';

interface WorkHistoryTimelineProps {
  career: CareerContent;
}

export default function WorkHistoryTimeline({ career }: WorkHistoryTimelineProps) {
  const locale = useLocale();
  const shouldReduceMotion = useReducedMotion();

  return (
    <section className="section-py">
      <div className="container">
        <ClipReveal direction="left" className="mb-2">
          <span className="font-mono text-xs text-primary uppercase tracking-[0.25em]">
            {locale === 'ru' ? '— карьера' : '— career'}
          </span>
        </ClipReveal>
        <ClipReveal direction="up" delay={0.1} className="mb-12">
          <h2 className="font-heading text-4xl md:text-5xl font-black">
            {locale === 'ru' ? 'Опыт работы' : 'Work Experience'}
          </h2>
        </ClipReveal>

        <div className="relative">
          <div className="absolute left-[19px] top-2 bottom-2 w-px bg-gradient-to-b from-primary/70 via-primary/30 to-transparent hidden md:block" />

          <div className="space-y-5">
            {career.workHistory.map((job, i) => {
              const company = locale === 'ru' ? job.company.ru : job.company.en;
              const position = locale === 'ru' ? job.position.ru : job.position.en;
              const description = locale === 'ru' ? job.description.ru : job.description.en;

              return (
                <motion.div
                  key={job.id}
                  className="flex gap-6"
                  initial={{ opacity: 0, x: shouldReduceMotion ? 0 : -32 }}
                  whileInView={{ opacity: 1, x: 0 }}
                  viewport={{ once: true, margin: '-60px' }}
                  transition={{ duration: 0.5, delay: shouldReduceMotion ? 0 : i * 0.09 }}
                >
                  <div className="hidden md:flex flex-col items-center flex-shrink-0">
                    <div className={`w-10 h-10 rounded-xl border flex items-center justify-center z-10 transition-colors ${
                      job.current
                        ? 'bg-primary/20 border-primary/50'
                        : 'bg-muted border-border'
                    }`}>
                      <Briefcase
                        className={`h-4 w-4 ${job.current ? 'text-primary' : 'text-muted-foreground'}`}
                        strokeWidth={1.5}
                      />
                    </div>
                  </div>

                  <div className={`flex-1 glass-card p-6 transition-all duration-300 ${
                    job.current
                      ? 'border-primary/25 hover:border-primary/40'
                      : 'hover:border-primary/20'
                  }`}>
                    <div className="flex flex-wrap items-start justify-between gap-3 mb-3">
                      <div>
                        <h3 className="font-heading font-bold text-lg leading-tight">{position}</h3>
                        <p className="text-sm text-muted-foreground mt-0.5">{company}</p>
                      </div>
                      <div className="flex items-center gap-2 shrink-0">
                        {job.current && (
                          <Badge variant="success" className="text-xs gap-1.5 font-mono">
                            <span
                              className="w-1.5 h-1.5 rounded-full bg-emerald-400"
                              style={{ animation: 'pulse-dot 1.8s ease-in-out infinite' }}
                            />
                            {locale === 'ru' ? 'Сейчас' : 'Current'}
                          </Badge>
                        )}
                        <span className="text-xs text-muted-foreground font-mono">
                          {formatMonthYear(job.startDate, locale)} — {job.endDate ? formatMonthYear(job.endDate, locale) : (locale === 'ru' ? 'наст. вр.' : 'present')}
                        </span>
                      </div>
                    </div>
                    <p className="text-sm text-muted-foreground leading-relaxed mb-4">{description}</p>
                    {job.technologies && job.technologies.length > 0 && (
                      <div className="flex flex-wrap gap-1.5">
                        {job.technologies.map((tech) => (
                          <span
                            key={tech}
                            className="px-2 py-0.5 rounded text-xs font-mono text-muted-foreground bg-muted border border-border"
                          >
                            {tech}
                          </span>
                        ))}
                      </div>
                    )}
                  </div>
                </motion.div>
              );
            })}
          </div>
        </div>
      </div>
    </section>
  );
}
