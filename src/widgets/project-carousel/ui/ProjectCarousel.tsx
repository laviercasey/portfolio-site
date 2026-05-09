'use client';

import { useState, useEffect, useCallback } from 'react';
import { useLocale } from 'next-intl';
import { motion, AnimatePresence, useReducedMotion } from 'framer-motion';
import { Link } from '@/shared/config';
import Image from 'next/image';
import { ArrowRight, ChevronLeft, ChevronRight, Star } from 'lucide-react';
import { Button, Badge, ClipReveal } from '@/shared/ui';
import type { Project } from '@/entities/project';
import type { HomepageVisibility } from '@/entities/content';

interface ProjectCarouselProps {
  projects: Project[];
  visibility?: HomepageVisibility;
}

export default function ProjectCarousel({ projects, visibility }: ProjectCarouselProps) {
  const locale = useLocale();
  const shouldReduceMotion = useReducedMotion();
  const [current, setCurrent] = useState(0);
  const [paused, setPaused] = useState(false);

  const next = useCallback(() => setCurrent((i) => (i + 1) % projects.length), [projects.length]);
  const prev = () => setCurrent((i) => (i - 1 + projects.length) % projects.length);

  useEffect(() => {
    if (paused || projects.length <= 1 || shouldReduceMotion) return;
    const interval = setInterval(next, 4500);
    return () => clearInterval(interval);
  }, [paused, next, projects.length, shouldReduceMotion]);

  if (projects.length === 0) return null;

  const project = projects[current];
  const title = locale === 'ru' ? project.title.ru : project.title.en;
  const description = locale === 'ru' ? project.shortDescription.ru : project.shortDescription.en;

  const statusLabel: Record<string, string> = {
    completed:      locale === 'ru' ? 'Завершён' : 'Completed',
    in_development: locale === 'ru' ? 'В разработке' : 'In Development',
    open_source:    'Open Source',
  };
  const statusVariant: Record<string, 'success' | 'warning' | 'info'> = {
    completed:      'success',
    in_development: 'warning',
    open_source:    'info',
  };

  return (
    <section className="section-py">
      <div className="container">
        <div className="flex items-end justify-between mb-10 gap-4">
          <div>
            <ClipReveal direction="left">
              <span className="font-mono text-xs text-primary uppercase tracking-[0.25em] block mb-2">
                {locale === 'ru' ? '— проекты' : '— projects'}
              </span>
            </ClipReveal>
            <ClipReveal direction="up" delay={0.1}>
              <h2 className="font-heading text-4xl md:text-5xl font-black">
                {locale === 'ru' ? 'Проекты' : 'Projects'}
              </h2>
            </ClipReveal>
          </div>
          <ClipReveal direction="right" delay={0.15}>
            <Button asChild variant="ghost" className="gap-2 text-muted-foreground hover:text-primary font-mono text-sm">
              <Link href="/projects">
                {locale === 'ru' ? 'Все проекты' : 'All projects'}
                <ArrowRight className="h-3.5 w-3.5" />
              </Link>
            </Button>
          </ClipReveal>
        </div>

        <div
          className="relative glass-card overflow-hidden hover:border-primary/20 transition-all duration-300"
          onMouseEnter={() => setPaused(true)}
          onMouseLeave={() => setPaused(false)}
        >
          <AnimatePresence mode="wait">
            <motion.div
              key={current}
              initial={{ opacity: 0, x: shouldReduceMotion ? 0 : 50 }}
              animate={{ opacity: 1, x: 0 }}
              exit={{ opacity: 0, x: shouldReduceMotion ? 0 : -50 }}
              transition={{ duration: 0.35, ease: [0.25, 0.1, 0.25, 1] }}
              className="grid grid-cols-1 md:grid-cols-2"
            >
              <div className="relative h-56 md:h-80 overflow-hidden bg-gradient-to-br from-primary/10 to-muted">
                {project.thumbnailUrl ? (
                  <Image src={project.thumbnailUrl} alt={title} fill className="object-cover" />
                ) : (
                  <div className="w-full h-full flex items-center justify-center">
                    <span className="font-heading italic font-black text-6xl text-primary/15">
                      {project.category.slice(0, 2).toUpperCase()}
                    </span>
                  </div>
                )}
                <div className="absolute inset-0 bg-gradient-to-r from-transparent to-[var(--background)]/20" />
              </div>

              <div className="p-8 md:p-10 flex flex-col justify-between">
                <div>
                  <div className="flex items-center gap-2 mb-4">
                    {visibility?.showProjectStatus !== false && (
                      <Badge variant={statusVariant[project.status]} className="font-mono text-xs">
                        {statusLabel[project.status]}
                      </Badge>
                    )}
                    {visibility?.showProjectStars !== false && (project.stars ?? 0) > 0 && (
                      <span className="flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-mono text-yellow-400 bg-yellow-400/10 border border-yellow-400/20">
                        <Star className="h-3 w-3 fill-yellow-400" />
                        {project.stars}
                      </span>
                    )}
                  </div>
                  <h3 className="font-heading text-2xl lg:text-3xl font-bold mb-3 leading-tight">{title}</h3>
                  <p className="text-muted-foreground text-sm md:text-base lg:text-lg leading-relaxed mb-6">{description}</p>
                  {visibility?.showProjectTechStack !== false && (
                    <div className="flex flex-wrap gap-1.5 mb-6">
                      {project.techStack.slice(0, 5).map((tech) => (
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
                <Button
                  asChild
                  className="self-start gap-2 shadow-md"
                  style={{ boxShadow: '0 2px 16px var(--glow-primary)' }}
                >
                  <Link href={`/projects/${project.slug}`}>
                    {locale === 'ru' ? 'Подробнее' : 'View project'}
                    <ArrowRight className="h-4 w-4" />
                  </Link>
                </Button>
              </div>
            </motion.div>
          </AnimatePresence>

          <button
            onClick={prev}
            className="absolute left-3 top-1/2 -translate-y-1/2 glass-card w-9 h-9 flex items-center justify-center hover:border-primary/40 transition-all"
            aria-label="Previous"
          >
            <ChevronLeft className="h-4 w-4" />
          </button>
          <button
            onClick={next}
            className="absolute right-3 top-1/2 -translate-y-1/2 glass-card w-9 h-9 flex items-center justify-center hover:border-primary/40 transition-all"
            aria-label="Next"
          >
            <ChevronRight className="h-4 w-4" />
          </button>

          <div className="absolute bottom-4 left-1/2 -translate-x-1/2 flex gap-2">
            {projects.map((_, i) => (
              <button
                key={i}
                onClick={() => setCurrent(i)}
                className={`rounded-full transition-all duration-300 ${
                  i === current ? 'w-6 h-1.5 bg-primary' : 'w-1.5 h-1.5 bg-muted-foreground/30'
                }`}
                aria-label={`Slide ${i + 1}`}
              />
            ))}
          </div>
        </div>
      </div>
    </section>
  );
}
