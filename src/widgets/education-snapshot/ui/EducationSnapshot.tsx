'use client';

import { useLocale } from 'next-intl';
import Image from 'next/image';
import { motion, useReducedMotion } from 'framer-motion';
import { GraduationCap } from 'lucide-react';
import { ClipReveal } from '@/shared/ui';
import type { CareerContent } from '@/entities/career';

interface EducationSnapshotProps {
  career: CareerContent;
}

export default function EducationSnapshot({ career }: EducationSnapshotProps) {
  const locale = useLocale();
  const shouldReduceMotion = useReducedMotion();
  const topEducation = career.education.slice(0, 3);

  return (
    <section className="section-py">
      <div className="container">
        <ClipReveal direction="left" className="mb-2">
          <span className="font-mono text-xs text-primary uppercase tracking-[0.25em]">
            {locale === 'ru' ? '— образование' : '— education'}
          </span>
        </ClipReveal>
        <ClipReveal direction="up" delay={0.1} className="mb-12">
          <h2 className="font-heading text-4xl md:text-5xl font-black">
            {locale === 'ru' ? 'Образование' : 'Education'}
          </h2>
        </ClipReveal>

        <div className="relative">
          <div className="absolute left-[19px] top-2 bottom-2 w-px bg-gradient-to-b from-primary/60 via-primary/20 to-transparent hidden md:block" />

          <div className="space-y-5">
            {topEducation.map((edu, i) => {
              const institution = locale === 'ru' ? edu.institution.ru : edu.institution.en;
              const degree = locale === 'ru' ? edu.degree.ru : edu.degree.en;
              const field = locale === 'ru' ? edu.field.ru : edu.field.en;

              return (
                <motion.div
                  key={edu.id}
                  className="flex gap-6"
                  initial={{ opacity: 0, x: shouldReduceMotion ? 0 : -28 }}
                  whileInView={{ opacity: 1, x: 0 }}
                  viewport={{ once: true, margin: '-60px' }}
                  transition={{ duration: 0.45, delay: shouldReduceMotion ? 0 : i * 0.1 }}
                >
                  <div className="hidden md:flex flex-col items-center flex-shrink-0">
                    <div className="w-10 h-10 rounded-xl bg-primary/10 border border-primary/20 flex items-center justify-center z-10">
                      {edu.logoUrl ? (
                        <Image src={edu.logoUrl} alt={institution} width={20} height={20} className="object-contain rounded" />
                      ) : (
                        <GraduationCap className="h-4 w-4 text-primary" strokeWidth={1.5} />
                      )}
                    </div>
                  </div>

                  <div className="flex-1 glass-card p-5 hover:border-primary/20 transition-all duration-300">
                    <div className="flex items-start justify-between gap-3">
                      <h3 className="font-heading font-bold text-sm md:text-base leading-snug">{institution}</h3>
                      <span className="text-xs font-mono text-muted-foreground bg-muted px-2 py-0.5 rounded shrink-0">
                        {edu.startYear}–{edu.endYear}
                      </span>
                    </div>
                    <p className="text-sm md:text-base text-muted-foreground mt-1">{degree}, {field}</p>
                  </div>
                </motion.div>
              );
            })}
          </div>
        </div>

        {career.certificates.length > 0 && (
          <motion.div
            className="mt-12"
            initial={{ opacity: 0, y: shouldReduceMotion ? 0 : 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true, margin: '-60px' }}
            transition={{ duration: 0.5 }}
          >
            <p className="font-mono text-xs text-muted-foreground uppercase tracking-[0.25em] mb-5 text-center">
              {locale === 'ru' ? '— сертификаты' : '— certificates'}
            </p>
            <div className="flex flex-wrap justify-center gap-3">
              {career.certificates.map((cert) => (
                <div
                  key={cert.id}
                  className="px-4 py-2 rounded-xl glass-card text-xs font-mono text-muted-foreground hover:border-primary/25 hover:text-primary transition-all duration-200"
                >
                  {locale === 'ru' ? cert.title.ru : cert.title.en}
                </div>
              ))}
            </div>
          </motion.div>
        )}
      </div>
    </section>
  );
}
