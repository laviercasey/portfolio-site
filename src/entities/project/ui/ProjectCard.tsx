'use client';

import { useRef } from 'react';
import { Link } from '@/shared/config';
import Image from 'next/image';
import { useLocale, useTranslations } from 'next-intl';
import { Github, ArrowRight, Star, ExternalLink } from 'lucide-react';
import {
  motion,
  useMotionValue,
  useSpring,
  useTransform,
  useReducedMotion,
} from 'framer-motion';
import { Button } from '@/shared/ui';
import StatusBadge from './StatusBadge';
import type { Project } from '../model/types';

interface ProjectCardProps {
  project: Project;
  index?: number;
}

const TILT_MAX = 6;
const PARALLAX_MAX = 10;
const SPRING_CONFIG = { stiffness: 220, damping: 22, mass: 0.4 };

export default function ProjectCard({ project, index = 0 }: ProjectCardProps) {
  const locale = useLocale();
  const t = useTranslations('projects');
  const shouldReduceMotion = useReducedMotion();
  const ref = useRef<HTMLDivElement>(null);

  const mouseX = useMotionValue(0.5);
  const mouseY = useMotionValue(0.5);

  const tiltX = useSpring(
    useTransform(mouseY, [0, 1], [TILT_MAX, -TILT_MAX]),
    SPRING_CONFIG,
  );
  const tiltY = useSpring(
    useTransform(mouseX, [0, 1], [-TILT_MAX, TILT_MAX]),
    SPRING_CONFIG,
  );
  const parallaxX = useSpring(
    useTransform(mouseX, [0, 1], [-PARALLAX_MAX, PARALLAX_MAX]),
    SPRING_CONFIG,
  );
  const parallaxY = useSpring(
    useTransform(mouseY, [0, 1], [-PARALLAX_MAX, PARALLAX_MAX]),
    SPRING_CONFIG,
  );

  const handleMouseMove = (e: React.MouseEvent<HTMLDivElement>) => {
    if (shouldReduceMotion) return;
    const rect = ref.current?.getBoundingClientRect();
    if (!rect) return;
    const px = (e.clientX - rect.left) / rect.width;
    const py = (e.clientY - rect.top) / rect.height;
    mouseX.set(px);
    mouseY.set(py);
    ref.current?.style.setProperty('--spot-x', `${px * 100}%`);
    ref.current?.style.setProperty('--spot-y', `${py * 100}%`);
  };

  const handleMouseLeave = () => {
    mouseX.set(0.5);
    mouseY.set(0.5);
  };

  const title = locale === 'ru' ? project.title.ru : project.title.en;
  const description = locale === 'ru' ? project.shortDescription.ru : project.shortDescription.en;
  const liveUrl = project.siteUrl || project.demoUrl;

  return (
    <motion.div
      ref={ref}
      onMouseMove={handleMouseMove}
      onMouseLeave={handleMouseLeave}
      className="group relative overflow-hidden rounded-xl flex flex-col"
      style={{
        rotateX: shouldReduceMotion ? 0 : tiltX,
        rotateY: shouldReduceMotion ? 0 : tiltY,
        transformStyle: 'preserve-3d',
        transformPerspective: 1000,
      }}
      initial={{ opacity: 0, y: shouldReduceMotion ? 0 : 24 }}
      whileInView={{ opacity: 1, y: 0 }}
      viewport={{ once: true, margin: '-60px' }}
      transition={{ duration: 0.45, delay: shouldReduceMotion ? 0 : index * 0.07 }}
    >
      <div
        className="absolute inset-0 rounded-xl glass-card group-hover:border-primary/25 transition-colors duration-300"
        aria-hidden
      />
      <div
        className="pointer-events-none absolute inset-0 rounded-xl opacity-0 group-hover:opacity-100 transition-opacity duration-300"
        style={{
          background:
            'radial-gradient(400px circle at var(--spot-x, 50%) var(--spot-y, 50%), rgba(255,255,255,0.08), transparent 60%)',
        }}
        aria-hidden
      />

      <div className="relative h-48 overflow-hidden rounded-t-xl bg-gradient-to-br from-primary/20 to-muted">
        {project.thumbnailUrl ? (
          <motion.div
            className="absolute inset-[-10px]"
            style={{
              x: shouldReduceMotion ? 0 : parallaxX,
              y: shouldReduceMotion ? 0 : parallaxY,
            }}
          >
            <Image
              src={project.thumbnailUrl}
              alt={title}
              fill
              className="object-cover group-hover:scale-105 transition-transform duration-500"
            />
          </motion.div>
        ) : (
          <div className="w-full h-full flex items-center justify-center font-heading text-4xl font-bold text-primary/20">
            {project.category.slice(0, 2).toUpperCase()}
          </div>
        )}
        <div className="absolute top-3 left-3 flex items-center gap-2 z-10">
          <StatusBadge status={project.status} />
          {(project.stars ?? 0) > 0 && (
            <span className="flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-mono bg-black/50 backdrop-blur-sm text-yellow-400 border border-yellow-400/20">
              <Star className="h-3 w-3 fill-yellow-400" />
              {project.stars}
            </span>
          )}
        </div>
      </div>

      <div className="relative p-5 flex flex-col flex-1">
        <h3 className="font-heading font-bold text-base mb-2 leading-snug group-hover:text-primary transition-colors">
          {title}
        </h3>
        <p className="text-sm text-muted-foreground mb-3 flex-1 leading-relaxed">{description}</p>

        <div className="flex flex-wrap gap-1.5 mb-3">
          {project.techStack.slice(0, 3).map((tech) => (
            <span
              key={tech}
              className="px-2 py-0.5 rounded text-xs font-mono text-muted-foreground bg-muted border border-border"
            >
              {tech}
            </span>
          ))}
          {project.techStack.length > 3 && (
            <span className="px-2 py-0.5 rounded text-xs font-mono text-muted-foreground bg-muted border border-border">
              +{project.techStack.length - 3}
            </span>
          )}
        </div>

        <div className="flex gap-2">
          <Button asChild size="sm" className="flex-1 gap-1.5">
            <Link href={`/projects/${project.slug}`}>
              <ArrowRight className="h-3.5 w-3.5" />
              {t('viewProject')}
            </Link>
          </Button>
          {liveUrl && (
            <Button
              asChild
              variant="outline"
              size="sm"
              className="glass-card border-white/15 hover:border-primary/30"
              title={locale === 'ru' ? 'Открыть живое демо' : 'Open live demo'}
            >
              <a href={liveUrl} target="_blank" rel="noopener noreferrer" aria-label="Live demo">
                <ExternalLink className="h-4 w-4" />
              </a>
            </Button>
          )}
          {project.githubUrl && (
            <Button
              asChild
              variant="outline"
              size="sm"
              className="glass-card border-white/15 hover:border-primary/30"
              title="GitHub"
            >
              <a href={project.githubUrl} target="_blank" rel="noopener noreferrer" aria-label="GitHub">
                <Github className="h-4 w-4" />
              </a>
            </Button>
          )}
        </div>
      </div>
    </motion.div>
  );
}
