export type ServiceIconKey = 'bot' | 'code' | 'layers' | 'database' | 'workflow';
export type ServiceVisualKey = 'terminal' | 'browser' | 'editor' | 'dataframe' | 'pipeline';

export interface I18nText {
  ru: string;
  en: string;
}

export interface ServiceBullets {
  ru: string[];
  en: string[];
}

export interface ServiceCaseProject {
  slug: string;
  name: string;
}

export interface Service {
  id: string;
  slug: string;
  num: string;
  iconKey: ServiceIconKey;
  visualKey: ServiceVisualKey;
  accent: string;
  title: I18nText;
  lead: I18nText;
  bullets: ServiceBullets;
  stack: string;
  timeline: I18nText;
  priceRu: string;
  priceEn: string;
  caseProjects: ServiceCaseProject[];
  order: number;
  createdAt: string;
  updatedAt: string;
}

export interface ServiceFaq {
  id: string;
  question: I18nText;
  answer: I18nText;
  order: number;
  createdAt: string;
  updatedAt: string;
}

export interface ServiceProcessStep {
  id: string;
  num: string;
  title: I18nText;
  description: I18nText;
  order: number;
  createdAt: string;
  updatedAt: string;
}

export interface ServicesPageData {
  services: Service[];
  faqs: ServiceFaq[];
  processSteps: ServiceProcessStep[];
}

export interface CreateServiceInput {
  slug: string;
  num: string;
  iconKey: ServiceIconKey;
  visualKey: ServiceVisualKey;
  accent: string;
  title: I18nText;
  lead: I18nText;
  bullets: ServiceBullets;
  stack: string;
  timeline: I18nText;
  priceRu: string;
  priceEn: string;
  caseProjects: ServiceCaseProject[];
  order: number;
}

export type UpdateServiceInput = Partial<CreateServiceInput>;

export interface CreateServiceFaqInput {
  question: I18nText;
  answer: I18nText;
  order: number;
}

export type UpdateServiceFaqInput = Partial<CreateServiceFaqInput>;

export interface CreateServiceProcessStepInput {
  num: string;
  title: I18nText;
  description: I18nText;
  order: number;
}

export type UpdateServiceProcessStepInput = Partial<CreateServiceProcessStepInput>;
