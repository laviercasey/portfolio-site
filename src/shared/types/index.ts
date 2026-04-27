export interface I18nString {
  en: string;
  ru: string;
}

export interface SiteConfig {
  name: string;
  nameRu?: string;
  taglineEn: string;
  taglineRu: string;
  descriptionEn: string;
  descriptionRu: string;
  url: string;
  author: string;
  email: string;
  ogImage?: string;
  keywordsEn?: string[];
  keywordsRu?: string[];
}
