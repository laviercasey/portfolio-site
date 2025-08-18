export interface Project {
  id: string;
  title: string;
  shortDescription: string;
  imageUrl: string;
  technologies: string[];
  category: 'frontend' | 'backend' | 'fullstack' | 'devops' | 'other';
  demoUrl?: string;
  githubUrl?: string;
  featured: boolean;
}
