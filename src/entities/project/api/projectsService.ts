import type { Project } from '../model/types';
import { serverApi, ApiError } from '@/shared/api';
import { normalizeI18n, normalizeI18nOptional } from '@/shared/lib';

type RawProject = Omit<Partial<Project>, 'shortDescription' | 'goalDescription' | 'order'> & {
  id?: string;
  slug?: string;
  shortDescription?: unknown;
  shortDesc?: unknown;
  goalDescription?: unknown;
  goalDesc?: unknown;
  order?: number;
  sortOrder?: number;
};

function normalizeProject(raw: RawProject): Project {
  const shortDescriptionRaw = raw.shortDescription ?? raw.shortDesc;
  const goalDescriptionRaw = raw.goalDescription ?? raw.goalDesc;
  const order = typeof raw.order === 'number'
    ? raw.order
    : typeof raw.sortOrder === 'number'
      ? raw.sortOrder
      : 0;

  return {
    id: raw.id ?? '',
    slug: raw.slug ?? '',
    title: normalizeI18n(raw.title),
    shortDescription: normalizeI18n(shortDescriptionRaw),
    description: normalizeI18n(raw.description),
    category: (raw.category as Project['category']) ?? 'other',
    status: (raw.status as Project['status']) ?? 'completed',
    tags: Array.isArray(raw.tags) ? raw.tags : [],
    techStack: Array.isArray(raw.techStack) ? raw.techStack : [],
    goalDescription: normalizeI18nOptional(goalDescriptionRaw),
    githubUrl: raw.githubUrl,
    demoUrl: raw.demoUrl,
    siteUrl: raw.siteUrl,
    videoUrl: raw.videoUrl,
    stars: raw.stars,
    thumbnailUrl: raw.thumbnailUrl,
    images: raw.images,
    featured: raw.featured ?? false,
    order,
    problem: normalizeI18nOptional(raw.problem),
    approach: normalizeI18nOptional(raw.approach),
    outcome: normalizeI18nOptional(raw.outcome),
    techChoices: raw.techChoices,
    highlights: raw.highlights,
    timelineStarted: raw.timelineStarted,
    timelineShipped: raw.timelineShipped,
    demoCredentials: raw.demoCredentials,
    createdAt: raw.createdAt ?? new Date(0).toISOString(),
  };
}

function isBuildTimeFetchFailure(error: unknown): boolean {
  if (error instanceof ApiError) return error.status >= 500;
  return true;
}

export const projectsService = {
  async list(category?: string): Promise<Project[]> {
    const q = category ? `?category=${category}` : '';
    try {
      const data = await serverApi.get<RawProject[]>(`/api/projects${q}`);
      return Array.isArray(data) ? data.map(normalizeProject) : [];
    } catch (error: unknown) {
      if (isBuildTimeFetchFailure(error)) return [];
      throw error;
    }
  },

  async getBySlug(slug: string): Promise<Project | null> {
    try {
      const raw = await serverApi.get<RawProject>(`/api/projects/${slug}`);
      return normalizeProject(raw);
    } catch (error: unknown) {
      if (error instanceof ApiError && error.status === 404) return null;
      if (isBuildTimeFetchFailure(error)) return null;
      throw error;
    }
  },

  create(token: string, data: Partial<Project>) {
    return serverApi.withToken(token).post<Project>('/api/projects', data);
  },

  update(token: string, id: string, data: Partial<Project>) {
    return serverApi.withToken(token).put<Project>(`/api/projects/${id}`, data);
  },

  delete(token: string, id: string) {
    return serverApi.withToken(token).delete(`/api/projects/${id}`);
  },
};
