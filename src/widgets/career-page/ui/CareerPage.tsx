'use client';

import { useState } from 'react';
import { motion, useReducedMotion, AnimatePresence } from 'framer-motion';
import { useTranslations } from 'next-intl';
import { GraduationCap, Briefcase, Award, BookOpen, ExternalLink, ChevronDown, FolderOpen } from 'lucide-react';
import { ClipReveal, Badge, Button } from '@/shared/ui';
import { Link } from '@/shared/config';
import type { CareerContent } from '@/entities/career';
import type { Project } from '@/entities/project';
import { formatMonthYear } from '@/shared/lib';

interface Props {
  career: CareerContent;
  locale: string;
  projects?: Project[];
}

const fadeUp = (i: number, reduce: boolean | null) =>
  reduce
    ? {}
    : {
        initial: { opacity: 0, y: 32 },
        whileInView: { opacity: 1, y: 0 },
        viewport: { once: true, margin: '-40px' },
        transition: { duration: 0.5, delay: i * 0.07, ease: [0.16, 1, 0.3, 1] as [number, number, number, number] },
      };

export default function CareerPage({ career, locale, projects = [] }: Props) {
  const t = useTranslations('career');
  const isRu = locale === 'ru';
  const reduce = useReducedMotion();
  const [expandedJobs, setExpandedJobs] = useState<Record<string, boolean>>({});
  const projectBySlug = new Map(projects.map((p) => [p.slug, p]));

  const earliestYear = career.workHistory.reduce((min, job) => {
    const y = parseInt(job.startDate.split('-')[0]);
    return y < min ? y : min;
  }, new Date().getFullYear());
  const yearsExp = new Date().getFullYear() - earliestYear;

  return (
    <div className="container py-20 md:py-28 max-w-6xl">
      <section className="mb-20 text-center">
        <ClipReveal>
          <h1 className="font-heading text-4xl sm:text-5xl md:text-6xl font-bold gradient-text mb-4">
            {t('title')}
          </h1>
        </ClipReveal>
        <ClipReveal delay={0.15}>
          <p className="text-lg text-muted-foreground mb-8">{t('subtitle')}</p>
        </ClipReveal>

        <div className="flex justify-center">
          <div className="flex gap-3 overflow-x-auto pb-2 px-1 scrollbar-none max-w-full">
            <StatPill value={`${yearsExp}+`} label={t('yearsExp')} />
            <StatPill value={String(career.certificates.length)} label={t('certsCount')} />
            <StatPill value={String(career.publications.length)} label={t('pubsCount')} />
          </div>
        </div>
      </section>

      <section className="mb-24">
        <SectionHeading icon={Briefcase} label={t('work')} />

        <div className="relative mt-10">
          <div className="absolute left-[19px] md:left-[23px] top-0 bottom-0 w-px bg-gradient-to-b from-primary/60 via-primary/20 to-transparent" />

          <div className="space-y-8">
            {career.workHistory.map((job, i) => (
              <motion.div key={job.id} {...fadeUp(i, reduce)} className="relative flex gap-5 md:gap-7">
                <div className="relative z-10 mt-6 flex-shrink-0">
                  <div
                    className={`w-[10px] h-[10px] md:w-3 md:h-3 rounded-full border-2 ${
                      job.current
                        ? 'border-primary bg-primary'
                        : 'border-primary/50 bg-background'
                    }`}
                    style={job.current ? { animation: 'pulse-dot 1.8s ease-in-out infinite' } : undefined}
                  />
                </div>

                <div
                  className={`flex-1 glass-card p-6 md:p-7 transition-all duration-300 ${
                    job.current ? 'glow-gold border-primary/25' : 'hover:border-primary/20'
                  }`}
                >
                  <div className="flex flex-wrap items-start justify-between gap-3 mb-3">
                    <div>
                      <h3 className="font-heading font-bold text-base md:text-lg">
                        {isRu ? job.position.ru : job.position.en}
                      </h3>
                      <p className="text-sm text-muted-foreground mt-0.5">
                        {isRu ? job.company.ru : job.company.en}
                      </p>
                    </div>
                    <div className="flex items-center gap-2 shrink-0">
                      {job.current && (
                        <Badge variant="success" className="text-xs gap-1.5">
                          <span
                            className="w-1.5 h-1.5 rounded-full bg-emerald-400"
                            style={{ animation: 'pulse-dot 1.8s ease-in-out infinite' }}
                          />
                          {t('currentRole')}
                        </Badge>
                      )}
                      <span className="text-xs font-mono text-muted-foreground bg-muted px-2 py-1 rounded">
                        {formatMonthYear(job.startDate, locale)} — {job.endDate ? formatMonthYear(job.endDate, locale) : t('present')}
                      </span>
                    </div>
                  </div>

                  <p className="text-sm text-muted-foreground leading-relaxed mb-4">
                    {isRu ? job.description.ru : job.description.en}
                  </p>

                  {job.technologies && (
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

                  {(() => {
                    const hasAch = (job.achievements?.length ?? 0) > 0;
                    const fullDesc = job.fullDescription && (isRu ? job.fullDescription.ru : job.fullDescription.en);
                    if (!hasAch && !fullDesc) return null;
                    const isOpen = !!expandedJobs[job.id];
                    return (
                      <>
                        <button
                          type="button"
                          onClick={() => setExpandedJobs((prev) => ({ ...prev, [job.id]: !prev[job.id] }))}
                          className="mt-4 flex items-center gap-1.5 text-xs font-mono uppercase tracking-wider text-primary hover:text-primary/80 transition-colors"
                          aria-expanded={isOpen}
                        >
                          <ChevronDown
                            className={`h-3.5 w-3.5 transition-transform duration-200 ${isOpen ? 'rotate-180' : ''}`}
                          />
                          {isOpen
                            ? (isRu ? 'Свернуть' : 'Collapse')
                            : (isRu ? 'Подробности о роли' : 'More about this role')}
                        </button>
                        <AnimatePresence initial={false}>
                          {isOpen && (
                            <motion.div
                              initial={{ height: 0, opacity: 0 }}
                              animate={{ height: 'auto', opacity: 1 }}
                              exit={{ height: 0, opacity: 0 }}
                              transition={{ duration: 0.28, ease: 'easeOut' }}
                              className="overflow-hidden"
                            >
                              <div className="pt-4 mt-4 border-t border-border/50 space-y-4">
                                {fullDesc && (
                                  <p className="text-sm text-muted-foreground leading-relaxed whitespace-pre-line">
                                    {fullDesc}
                                  </p>
                                )}
                                {hasAch && (
                                  <ul className="space-y-2">
                                    {job.achievements!.map((ach, k) => (
                                      <li key={k} className="flex gap-2 text-sm text-muted-foreground leading-relaxed">
                                        <span className="text-primary shrink-0 mt-1.5">—</span>
                                        <span>{isRu ? ach.ru : ach.en}</span>
                                      </li>
                                    ))}
                                  </ul>
                                )}
                              </div>
                            </motion.div>
                          )}
                        </AnimatePresence>
                      </>
                    );
                  })()}
                </div>
              </motion.div>
            ))}
          </div>
        </div>
      </section>

      <section className="mb-24">
        <SectionHeading icon={GraduationCap} label={t('education')} />

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-5 mt-10">
          {career.education.map((edu, i) => (
            <motion.div
              key={edu.id}
              {...fadeUp(i, reduce)}
              className="gradient-border glass-card p-6 md:p-7 hover:scale-[1.01] transition-transform duration-300"
            >
              <div className="flex items-start gap-4 mb-3">
                <div className="p-2.5 rounded-xl bg-primary/10 border border-primary/20 shrink-0">
                  <GraduationCap className="h-5 w-5 text-primary" strokeWidth={1.75} />
                </div>
                <div className="flex-1 min-w-0">
                  <h3 className="font-heading font-bold text-base">
                    {isRu ? edu.institution.ru : edu.institution.en}
                  </h3>
                  <p className="text-sm text-muted-foreground mt-0.5">
                    {isRu ? edu.degree.ru : edu.degree.en},{' '}
                    {isRu ? edu.field.ru : edu.field.en}
                  </p>
                </div>
                <Badge variant="outline" className="font-mono text-xs shrink-0">
                  {edu.startYear}–{edu.endYear}
                </Badge>
              </div>
              {edu.description && (
                <p className="text-sm text-muted-foreground leading-relaxed pl-[52px]">
                  {isRu ? edu.description.ru : edu.description.en}
                </p>
              )}

              {(() => {
                const slugs = edu.relatedProjectSlugs ?? [];
                const linked = slugs.map((s) => projectBySlug.get(s)).filter(Boolean) as Project[];
                if (linked.length === 0) return null;
                return (
                  <div className="mt-4 pt-4 border-t border-border/50 pl-[52px]">
                    <p className="text-[10px] uppercase tracking-[0.2em] text-muted-foreground mb-2 flex items-center gap-1.5">
                      <FolderOpen className="h-3 w-3 text-primary" />
                      {isRu ? 'Проекты во время обучения' : 'Projects during studies'}
                    </p>
                    <ul className="space-y-1.5">
                      {linked.map((p) => (
                        <li key={p.slug}>
                          <Link
                            href={`/projects/${p.slug}`}
                            className="group flex items-start gap-2 text-sm hover:text-primary transition-colors"
                          >
                            <span className="text-primary/40 group-hover:text-primary shrink-0 mt-0.5">→</span>
                            <span className="min-w-0">
                              <span className="font-medium">{isRu ? p.title.ru : p.title.en}</span>
                              <span className="text-xs text-muted-foreground ml-2">
                                {isRu ? p.shortDescription.ru : p.shortDescription.en}
                              </span>
                            </span>
                          </Link>
                        </li>
                      ))}
                    </ul>
                  </div>
                );
              })()}
            </motion.div>
          ))}
        </div>
      </section>

      <section className="mb-24">
        <SectionHeading icon={Award} label={t('certificates')} />

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mt-10">
          {career.certificates.map((cert, i) => (
            <motion.div
              key={cert.id}
              {...fadeUp(i, reduce)}
              className="glass-card p-5 flex items-start justify-between gap-3 hover:glow-gold hover:border-primary/20 transition-all duration-300 group"
            >
              <div className="min-w-0">
                <h3 className="font-semibold text-sm group-hover:text-primary transition-colors">
                  {isRu ? cert.title.ru : cert.title.en}
                </h3>
                <p className="text-xs text-muted-foreground mt-1 font-mono">
                  {isRu ? cert.issuer.ru : cert.issuer.en} · {formatMonthYear(cert.date, locale)}
                </p>
                {cert.credentialId && (
                  <p className="text-xs text-muted-foreground mt-0.5 font-mono">
                    ID: {cert.credentialId}
                  </p>
                )}
              </div>
              {cert.url && (
                <Button asChild variant="ghost" size="sm" className="flex-shrink-0 h-8 w-8 p-0">
                  <a href={cert.url} target="_blank" rel="noopener noreferrer" aria-label={t('viewCertificate')}>
                    <ExternalLink className="h-3.5 w-3.5" />
                  </a>
                </Button>
              )}
            </motion.div>
          ))}
        </div>
      </section>

      {career.publications.length > 0 && (
        <section>
          <SectionHeading icon={BookOpen} label={t('publications')} />

          <div className="space-y-5 mt-10">
            {career.publications.map((pub, i) => (
              <motion.div
                key={pub.id}
                {...fadeUp(i, reduce)}
                className="glass-card p-6 md:p-7 hover:border-primary/20 transition-all duration-300 flex gap-6"
              >
                <span className="hidden sm:flex items-start justify-center text-4xl font-heading font-bold gradient-text opacity-60 w-12 shrink-0 pt-1">
                  {String(i + 1).padStart(2, '0')}
                </span>

                <div className="flex-1 min-w-0">
                  <div className="flex flex-wrap items-start justify-between gap-3 mb-2">
                    <h3 className="font-heading font-bold text-sm md:text-base">
                      {isRu ? pub.title.ru : pub.title.en}
                    </h3>
                    <Badge variant="outline" className="font-mono text-xs shrink-0">
                      {pub.year}
                    </Badge>
                  </div>

                  {pub.journal && (
                    <p className="text-sm text-muted-foreground italic mb-3">
                      {isRu ? pub.journal.ru : pub.journal.en}
                    </p>
                  )}

                  {pub.abstract && (
                    <p className="text-sm text-muted-foreground leading-relaxed mb-3">
                      {isRu ? pub.abstract.ru : pub.abstract.en}
                    </p>
                  )}

                  {pub.doi && (
                    <a
                      href={`https://doi.org/${pub.doi}`}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="inline-flex items-center gap-1.5 text-xs text-primary hover:underline font-mono"
                    >
                      <ExternalLink className="h-3 w-3" />
                      DOI: {pub.doi}
                    </a>
                  )}
                </div>
              </motion.div>
            ))}
          </div>
        </section>
      )}
    </div>
  );
}


function StatPill({ value, label }: { value: string; label: string }) {
  return (
    <div className="flex items-center gap-2 px-4 py-2 rounded-full glass-card border border-primary/15 whitespace-nowrap">
      <span className="text-lg font-bold gradient-text">{value}</span>
      <span className="text-xs text-muted-foreground">{label}</span>
    </div>
  );
}

function SectionHeading({
  icon: Icon,
  label,
}: {
  icon: React.ComponentType<{ className?: string; strokeWidth?: number }>;
  label: string;
}) {
  return (
    <ClipReveal>
      <div className="flex items-center gap-3">
        <div className="p-2.5 rounded-xl bg-primary/10 border border-primary/20">
          <Icon className="h-5 w-5 text-primary" strokeWidth={1.75} />
        </div>
        <h2 className="font-heading text-2xl md:text-3xl font-bold">{label}</h2>
      </div>
    </ClipReveal>
  );
}
