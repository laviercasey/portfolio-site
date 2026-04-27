'use client';

import { useLocale, useTranslations } from 'next-intl';
import { Link } from '@/shared/config';
import Image from 'next/image';
import { motion, useReducedMotion } from 'framer-motion';
import { Badge, Button } from '@/shared/ui';
import { StatusBadge } from '@/entities/project';
import {
  Github,
  ExternalLink,
  ArrowLeft,
  Star,
  Target,
  Calendar,
  Images,
  Sparkles,
  Layers,
  Play,
} from 'lucide-react';
import type { Project } from '@/entities/project';
import { formatMonthYear } from '@/shared/lib';
import ProjectCredentials from './ProjectCredentials';
import ProjectHighlights from './ProjectHighlights';
import ProjectTechChoices from './ProjectTechChoices';
import ProjectGallery from './ProjectGallery';
import ProjectStoryBlock from './ProjectStoryBlock';

interface Props {
  project: Project;
}

function hasContent(i18n?: { en?: string; ru?: string }): boolean {
  if (!i18n) return false;
  return Boolean((i18n.en ?? '').trim() || (i18n.ru ?? '').trim());
}

function isDirectVideoFile(url: string): boolean {
  if (url.startsWith('/')) return true;
  return /\.(mp4|webm|mov|ogg|m4v)(\?.*)?$/i.test(url);
}

export default function ProjectDetail({ project }: Props) {
  const locale = useLocale();
  const t = useTranslations('projects');
  const shouldReduce = useReducedMotion();
  const isRu = locale === 'ru';

  const title = isRu ? project.title.ru : project.title.en;
  const shortDesc = isRu ? project.shortDescription.ru : project.shortDescription.en;
  const description = isRu ? project.description.ru : project.description.en;
  const goalDesc = project.goalDescription
    ? isRu ? project.goalDescription.ru : project.goalDescription.en
    : null;

  const hasProblem = hasContent(project.problem);
  const hasApproach = hasContent(project.approach);
  const hasOutcome = hasContent(project.outcome);
  const hasStory = hasProblem || hasApproach || hasOutcome;
  const hasHighlights = (project.highlights?.length ?? 0) > 0;
  const hasTechChoices = (project.techChoices?.length ?? 0) > 0;
  const hasCreds = (project.demoCredentials?.length ?? 0) > 0;
  const hasGallery = (project.images?.length ?? 0) > 0;
  const liveUrl = project.siteUrl || project.demoUrl;

  const timelineStarted = project.timelineStarted ? formatMonthYear(project.timelineStarted, locale) : null;
  const timelineShipped = project.timelineShipped
    ? formatMonthYear(project.timelineShipped, locale)
    : project.status === 'in_development'
      ? t('timelineInDev')
      : null;

  const fadeUp = shouldReduce
    ? {}
    : { initial: { opacity: 0, y: 16 }, animate: { opacity: 1, y: 0 } };

  return (
    <div className="container py-10 md:py-14 max-w-6xl">
      <motion.div {...fadeUp} transition={{ duration: 0.3 }}>
        <Button
          asChild
          variant="ghost"
          size="sm"
          className="-ml-3 mb-6 text-muted-foreground hover:text-primary"
        >
          <Link href="/projects">
            <ArrowLeft className="mr-1.5 h-4 w-4" />
            {t('backToProjects')}
          </Link>
        </Button>
      </motion.div>

      <motion.section
        className="glass-card rounded-2xl overflow-hidden mb-10 lg:mb-14"
        {...fadeUp}
        transition={{ duration: 0.4 }}
      >
        <div className="grid grid-cols-1 lg:grid-cols-[minmax(0,1.1fr)_minmax(0,1fr)]">
          <div className="relative aspect-video lg:aspect-auto lg:min-h-[360px] bg-gradient-to-br from-primary/10 to-muted">
            {project.videoUrl && isDirectVideoFile(project.videoUrl) ? (
              <video
                src={project.videoUrl}
                poster={project.thumbnailUrl}
                controls
                preload="none"
                className="w-full h-full object-cover"
              />
            ) : project.videoUrl ? (
              <a
                href={project.videoUrl}
                target="_blank"
                rel="noopener noreferrer"
                aria-label={`${title} — ${t('watchVideo')}`}
                className="group relative block w-full h-full"
              >
                {project.thumbnailUrl ? (
                  <Image src={project.thumbnailUrl} alt={title} fill className="object-cover" priority />
                ) : (
                  <div className="w-full h-full flex items-center justify-center">
                    <span className="font-heading italic font-black text-8xl text-primary/10">
                      {project.category.slice(0, 2).toUpperCase()}
                    </span>
                  </div>
                )}
                <div className="absolute inset-0 flex items-center justify-center bg-black/30 group-hover:bg-black/45 transition-colors">
                  <span className="flex items-center justify-center w-16 h-16 md:w-20 md:h-20 rounded-full bg-primary/95 text-primary-foreground shadow-lg backdrop-blur-sm group-hover:scale-110 transition-transform">
                    <Play className="w-7 h-7 md:w-9 md:h-9 ml-1" fill="currentColor" />
                  </span>
                </div>
              </a>
            ) : project.thumbnailUrl ? (
              <Image src={project.thumbnailUrl} alt={title} fill className="object-cover" priority />
            ) : (
              <div className="w-full h-full flex items-center justify-center">
                <span className="font-heading italic font-black text-8xl text-primary/10">
                  {project.category.slice(0, 2).toUpperCase()}
                </span>
              </div>
            )}
          </div>

          <div className="p-6 md:p-8 lg:p-10 flex flex-col justify-between gap-6">
            <div>
              <div className="flex items-center gap-2 mb-4 flex-wrap">
                <StatusBadge status={project.status} />
                <Badge variant="outline" className="font-mono text-xs">
                  {project.category}
                </Badge>
                {(project.stars ?? 0) > 0 && (
                  <span className="flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-mono text-yellow-400 bg-yellow-400/10 border border-yellow-400/20">
                    <Star className="h-3.5 w-3.5 fill-yellow-400" />
                    {project.stars} {t('stars')}
                  </span>
                )}
              </div>

              <h1 className="font-heading text-3xl md:text-4xl lg:text-5xl font-black mb-4 leading-[1.05]">
                {title}
              </h1>
              <p className="text-muted-foreground text-base md:text-lg leading-relaxed">
                {shortDesc}
              </p>
            </div>

            <div className="flex flex-wrap gap-3">
              {liveUrl && (
                <Button asChild className="gap-2">
                  <a href={liveUrl} target="_blank" rel="noopener noreferrer">
                    <ExternalLink className="h-4 w-4" />
                    {t('visitSite')}
                  </a>
                </Button>
              )}
              {project.githubUrl && (
                <Button
                  asChild
                  variant="outline"
                  className="gap-2 glass-card border-white/15 hover:border-primary/30"
                >
                  <a href={project.githubUrl} target="_blank" rel="noopener noreferrer">
                    <Github className="h-4 w-4" />
                    {t('viewSource')}
                  </a>
                </Button>
              )}
            </div>
          </div>
        </div>
      </motion.section>

      {hasHighlights && (
        <section className="mb-10 lg:mb-14 space-y-4">
          <h2 className="font-heading text-lg font-bold flex items-center gap-2">
            <Sparkles className="h-4 w-4 text-primary" />
            {t('highlightsTitle')}
          </h2>
          <ProjectHighlights highlights={project.highlights!} />
        </section>
      )}

      <div className="grid grid-cols-1 lg:grid-cols-[minmax(0,2fr)_minmax(0,1fr)] gap-8 lg:gap-12">
        <main className="min-w-0 space-y-10 lg:space-y-14">
          {hasProblem && (
            <ProjectStoryBlock
              label={t('problemLabel')}
              title={t('problemTitle')}
              accent="problem"
              index={0}
            >
              {isRu ? project.problem!.ru : project.problem!.en}
            </ProjectStoryBlock>
          )}

          {hasApproach && (
            <ProjectStoryBlock
              label={t('approachLabel')}
              title={t('approachTitle')}
              accent="approach"
              index={1}
            >
              {isRu ? project.approach!.ru : project.approach!.en}
            </ProjectStoryBlock>
          )}

          {hasTechChoices && <ProjectTechChoices choices={project.techChoices!} />}

          {hasOutcome && (
            <ProjectStoryBlock
              label={t('outcomeLabel')}
              title={t('outcomeTitle')}
              accent="outcome"
              index={2}
            >
              {isRu ? project.outcome!.ru : project.outcome!.en}
            </ProjectStoryBlock>
          )}

          {!hasStory && (
            <section className="glass-card p-6 rounded-xl">
              <h2 className="font-heading text-lg font-bold mb-4 flex items-center gap-2">
                <span className="w-1.5 h-1.5 rounded-full bg-primary" />
                {t('fullDescription')}
              </h2>
              <p className="text-muted-foreground leading-relaxed whitespace-pre-line">
                {description}
              </p>
            </section>
          )}

          {goalDesc && !hasProblem && (
            <section className="glass-card p-6 rounded-xl">
              <h2 className="font-heading text-lg font-bold mb-4 flex items-center gap-2">
                <Target className="h-4 w-4 text-primary" />
                {t('goal')}
              </h2>
              <p className="text-muted-foreground leading-relaxed">{goalDesc}</p>
            </section>
          )}

          {hasGallery && (
            <section className="space-y-4">
              <h2 className="font-heading text-lg font-bold flex items-center gap-2">
                <Images className="h-4 w-4 text-primary" />
                {t('galleryTitle')}
              </h2>
              <ProjectGallery images={project.images!} title={title} />
            </section>
          )}
        </main>

        <aside className="lg:sticky lg:top-24 lg:self-start space-y-4">
          {(timelineStarted || timelineShipped) && (
            <div className="glass-card rounded-xl p-5">
              <h3 className="font-mono text-[10px] uppercase tracking-[0.25em] text-muted-foreground mb-3 flex items-center gap-2">
                <Calendar className="h-3.5 w-3.5 text-primary" />
                Timeline
              </h3>
              <div className="font-mono text-sm">
                {timelineStarted && (
                  <div className="flex justify-between items-baseline">
                    <span className="text-muted-foreground text-xs">{t('timelineStarted')}</span>
                    <span className="text-foreground">{timelineStarted}</span>
                  </div>
                )}
                {timelineShipped && (
                  <div className="flex justify-between items-baseline mt-1.5">
                    <span className="text-muted-foreground text-xs">{t('timelineShipped')}</span>
                    <span className="text-foreground">{timelineShipped}</span>
                  </div>
                )}
              </div>
            </div>
          )}

          <div className="glass-card rounded-xl p-5">
            <h3 className="font-mono text-[10px] uppercase tracking-[0.25em] text-muted-foreground mb-3 flex items-center gap-2">
              <Layers className="h-3.5 w-3.5 text-primary" />
              {t('techStack')}
            </h3>
            <div className="flex flex-wrap gap-1.5">
              {project.techStack.map((tech) => (
                <Badge key={tech} variant="secondary" className="text-xs py-1 px-2 font-mono">
                  {tech}
                </Badge>
              ))}
            </div>
          </div>

          {hasCreds && (
            <div>
              <h3 className="font-mono text-[10px] uppercase tracking-[0.25em] text-muted-foreground mb-3 flex items-center gap-2">
                <ExternalLink className="h-3.5 w-3.5 text-primary" />
                {t('tryItTitle')}
              </h3>
              <ProjectCredentials credentials={project.demoCredentials!} />
            </div>
          )}
        </aside>
      </div>
    </div>
  );
}
