export type {
  Project,
  ProjectStatus,
  ProjectCategory,
  TechChoice,
  Highlight,
  DemoCredential,
  CredentialRole,
  I18nString,
} from './model/types';
export { projectsService } from './api/projectsService';
export { default as ProjectCard } from './ui/ProjectCard';
export { default as StatusBadge } from './ui/StatusBadge';
