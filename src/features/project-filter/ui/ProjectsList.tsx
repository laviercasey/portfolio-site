'use client';

import { useMemo } from 'react';
import { useSearchParams } from 'next/navigation';
import { useLocale } from 'next-intl';
import { ProjectCard, type Project } from '@/entities/project';

interface ProjectsListProps {
  projects: Project[];
}

export default function ProjectsList({ projects }: ProjectsListProps) {
  const searchParams = useSearchParams();
  const locale = useLocale();

  const filtered = useMemo(() => {
    const category = searchParams.get('category');
    const status = searchParams.get('status');
    const tag = searchParams.get('tag');

    return projects
      .filter((p) => {
        if (category && category !== 'all' && p.category !== category) return false;
        if (status && p.status !== status) return false;
        if (tag && !p.tags.includes(tag)) return false;
        return true;
      })
      .sort((a, b) => a.order - b.order);
  }, [projects, searchParams]);

  if (filtered.length === 0) {
    return (
      <div className="text-center py-24 text-muted-foreground glass-card rounded-2xl">
        {locale === 'ru' ? 'Проекты не найдены' : 'No projects found'}
      </div>
    );
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
      {filtered.map((project) => (
        <ProjectCard key={project.id} project={project} />
      ))}
    </div>
  );
}
