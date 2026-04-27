export interface Education {
  id: string;
  institution: { en: string; ru: string };
  degree: { en: string; ru: string };
  field: { en: string; ru: string };
  startYear: string;
  endYear: string;
  logoUrl?: string;
  description?: { en: string; ru: string };
  relatedProjectSlugs?: string[];
}

export interface WorkExperience {
  id: string;
  company: { en: string; ru: string };
  position: { en: string; ru: string };
  startDate: string;
  endDate?: string;
  current: boolean;
  description: { en: string; ru: string };
  technologies?: string[];
  logoUrl?: string;
  achievements?: Array<{ en: string; ru: string }>;
  fullDescription?: { en: string; ru: string };
}

export interface Certificate {
  id: string;
  title: { en: string; ru: string };
  issuer: { en: string; ru: string };
  date: string;
  url?: string;
  imageUrl?: string;
  credentialId?: string;
}

export interface Publication {
  id: string;
  title: { en: string; ru: string };
  journal?: { en: string; ru: string };
  year: string;
  doi?: string;
  url?: string;
  abstract?: { en: string; ru: string };
}

export interface CareerContent {
  education: Education[];
  workHistory: WorkExperience[];
  certificates: Certificate[];
  publications: Publication[];
}
