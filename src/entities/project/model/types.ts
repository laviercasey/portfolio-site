import type { I18nString } from '@/shared/types';

export type { I18nString };
export type ProjectStatus = 'completed' | 'in_development' | 'open_source';
export type ProjectCategory = 'web' | 'mobile' | 'data' | 'research' | 'other';

export interface TechChoice {
  id?: string;
  tech: string;
  reason: I18nString;
}

export interface Highlight {
  id?: string;
  label: I18nString;
  value: string;
}

export type CredentialRole = 'admin' | 'manager' | 'user' | 'custom';

export interface DemoCredential {
  id?: string;
  role: CredentialRole;
  label: I18nString;
  login: string;
  password: string;
  note: I18nString;
}

export interface Project {
  id: string;
  slug: string;
  title: I18nString;
  shortDescription: I18nString;
  description: I18nString;
  category: ProjectCategory;
  status: ProjectStatus;
  tags: string[];
  techStack: string[];
  goalDescription?: I18nString;
  githubUrl?: string;
  demoUrl?: string;
  siteUrl?: string;
  videoUrl?: string;
  stars?: number;
  thumbnailUrl?: string;
  images?: string[];
  featured: boolean;
  order: number;
  problem?: I18nString;
  approach?: I18nString;
  outcome?: I18nString;
  techChoices?: TechChoice[];
  highlights?: Highlight[];
  timelineStarted?: string;
  timelineShipped?: string;
  demoCredentials?: DemoCredential[];
  createdAt: string;
}
