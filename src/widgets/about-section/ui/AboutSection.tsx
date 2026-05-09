'use client';

import { useTranslations, useLocale } from 'next-intl';
import { useReducedMotion } from 'framer-motion';
import {
  Briefcase,
  GraduationCap,
  Star,
  FolderKanban,
  BookOpen,
} from 'lucide-react';
import { ClipReveal } from '@/shared/ui';
import type { HomepageContent } from '@/entities/content';
import type { Project } from '@/entities/project';

interface AboutSectionProps {
  content: HomepageContent;
  projects?: Project[];
}

export default function AboutSection({ content, projects = [] }: AboutSectionProps) {
  const t = useTranslations('about');
  const locale = useLocale();
  const shouldReduceMotion = useReducedMotion();

  const bio = locale === 'ru' ? content.about.bioRu : content.about.bioEn;
  const bioFull = locale === 'ru' ? content.about.bioFullRu : content.about.bioFullEn;
  const hasFullBio = Boolean(bioFull && bioFull.trim().length > 0 && bioFull.trim() !== bio.trim());
  const { stats } = content.about;
  const vis = content.visibility;
  const githubStars = projects.reduce((sum, p) => sum + (p.stars ?? 0), 0);

  const allStats = [
    { key: 'showStatProjects'      as const, icon: FolderKanban,   value: `${stats.projects}+`,       label: t('projects'),     color: 'text-primary',     bg: 'bg-primary/12' },
    { key: 'showStatYears'         as const, icon: Briefcase,      value: `${stats.yearsExperience}+`, label: t('experience'),   color: 'text-amber-400',   bg: 'bg-amber-500/12' },
    { key: 'showStatCertificates'  as const, icon: GraduationCap,  value: `${stats.certificates}+`,    label: t('certificates'), color: 'text-emerald-400', bg: 'bg-emerald-500/12' },
    { key: 'showStatGithubStars'   as const, icon: Star,           value: `${githubStars}`,            label: t('githubStars'),  color: 'text-yellow-400',  bg: 'bg-yellow-500/12' },
    { key: 'showStatHabrArticles'  as const, icon: BookOpen,       value: `${stats.habrArticles}`,     label: t('habrArticles'), color: 'text-cyan-400',    bg: 'bg-cyan-500/12' },
  ];
  const statItems = allStats.filter(s => vis?.[s.key] !== false);

  return (
    <section id="about" className="section-py">
      <div className="container">
        <ClipReveal direction="left" delay={0} className="mb-8">
          <span className="font-mono text-xs text-primary uppercase tracking-[0.25em]">
            {t('label')}
          </span>
        </ClipReveal>

        <div className="grid grid-cols-12 gap-4 md:gap-5 md:items-start">

          {vis?.showAboutGif !== false && (
            <ClipReveal
              direction="left"
              delay={0.1}
              className="col-span-12 md:col-span-5"
            >
              <div className="relative h-[360px] md:h-[440px] rounded-2xl overflow-hidden group border border-white/10 bg-gradient-to-br from-primary/5 via-background to-secondary/5">
                {content.about.gifUrl ? (
                  (() => {
                    const url = content.about.gifUrl;
                    const isSvg = /\.svg$/i.test(url);
                    const staticUrl = isSvg ? url.replace(/\.svg$/i, '-static.svg') : url;
                    const alt = locale === 'ru' ? 'Кейси Лавьер' : 'Casey Laviere';
                    return (
                      <img
                        src={staticUrl}
                        alt={alt}
                        loading="lazy"
                        decoding="async"
                        className="w-full h-full object-cover"
                      />
                    );
                  })()
                ) : (
                  <div className="w-full h-full flex flex-col items-center justify-center gap-4">
                    <div className="relative">
                      <div className="w-24 h-24 rounded-2xl bg-primary/10 border border-primary/20 flex items-center justify-center">
                        <span className="font-heading italic font-black text-4xl gradient-text">CL</span>
                      </div>
                      <div className="absolute inset-0 rounded-2xl border border-primary/30 animate-ping opacity-20" />
                    </div>
                    <span className="font-mono text-xs text-muted-foreground/40 uppercase tracking-wider">
                      GIF placeholder
                    </span>
                  </div>
                )}

                <div className="absolute bottom-0 left-0 right-0 p-5 bg-gradient-to-t from-black/70 via-black/30 to-transparent">
                  <p className="font-mono text-xs text-white/60 uppercase tracking-widest">{locale === 'ru' ? 'Кейси Лавьер' : 'Casey Laviere'}</p>
                  <p className="font-mono text-xs text-primary">IT Developer / Freelance</p>
                </div>
              </div>
            </ClipReveal>
          )}

          {vis?.showAboutBio !== false && (
            <ClipReveal direction="up" delay={0.15} className={`col-span-12 ${vis?.showAboutGif !== false ? 'md:col-span-7' : 'md:col-span-12'}`}>
              <div className="glass-card p-7 flex flex-col h-[360px] md:h-[440px]">
                <h2 className="font-heading text-3xl md:text-4xl font-bold mb-4 leading-tight shrink-0">
                  {t('title')}
                </h2>
                <div className="flex-1 overflow-y-auto pr-2 -mr-2">
                  <p className="text-sm md:text-base lg:text-lg text-muted-foreground leading-relaxed whitespace-pre-line">{bio}</p>
                  {hasFullBio && (
                    <div className="pt-4 mt-4 border-t border-white/8">
                      <p className="text-sm md:text-base lg:text-lg text-muted-foreground leading-relaxed whitespace-pre-line">
                        {bioFull}
                      </p>
                    </div>
                  )}
                </div>
              </div>
            </ClipReveal>
          )}

          <div className="col-span-12">
            <div className="grid grid-cols-3 sm:grid-cols-[repeat(auto-fit,minmax(140px,1fr))] gap-3 md:gap-4">
              {statItems.map(({ icon: Icon, value, label, color, bg }, i) => (
                <ClipReveal
                  key={label}
                  direction="up"
                  delay={shouldReduceMotion ? 0 : 0.2 + i * 0.06}
                >
                  <div className="glass-card p-4 md:p-5 flex flex-col items-center justify-center text-center h-full hover:border-primary/30 transition-all duration-300 group">
                    <div className={`p-2 rounded-xl ${bg} mb-2.5`}>
                      <Icon className={`h-4 w-4 ${color} group-hover:scale-110 transition-transform`} strokeWidth={1.5} />
                    </div>
                    <span className="font-heading text-2xl md:text-3xl font-black gradient-text leading-none">
                      {value}
                    </span>
                    <span className="text-[10px] md:text-xs text-muted-foreground mt-1.5 font-mono uppercase tracking-wider leading-tight">
                      {label}
                    </span>
                  </div>
                </ClipReveal>
              ))}
            </div>
          </div>

        </div>
      </div>
    </section>
  );
}
