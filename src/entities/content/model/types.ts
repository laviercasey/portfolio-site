export interface SocialLink {
  platform: string;
  url: string;
  icon: string;
}

export interface Achievement {
  id: string;
  title: { en: string; ru: string };
  description: { en: string; ru: string };
  year: string;
  type: 'publication' | 'award' | 'certificate' | 'achievement';
  url?: string;
}

export interface HowIWorkStep {
  titleRu: string;
  titleEn: string;
  descRu: string;
  descEn: string;
  gifUrl?: string;
}

export interface HowIWorkConfig {
  steps: HowIWorkStep[];
  philosophyTitleRu: string;
  philosophyTitleEn: string;
  philosophyTextRu: string;
  philosophyTextEn: string;
  payment1Ru: string;
  payment1En: string;
  payment1DescRu: string;
  payment1DescEn: string;
  payment2Ru: string;
  payment2En: string;
  payment2DescRu: string;
  payment2DescEn: string;
}

export interface ContactCTAConfig {
  headingRu: string;
  headingEn: string;
  subtitleRu: string;
  subtitleEn: string;
  bgWordRu: string;
  bgWordEn: string;
}

export interface HomepageVisibility {
  showAboutGif: boolean;
  showAboutBio: boolean;
  showStatProjects: boolean;
  showStatYears: boolean;
  showStatCertificates: boolean;
  showStatGithubStars: boolean;
  showStatHabrArticles: boolean;
  showProjectStars: boolean;
  showProjectStatus: boolean;
  showProjectTechStack: boolean;
  showMarquee: boolean;
  showHowIWork: boolean;
  showAchievements: boolean;
  showContactCTA: boolean;
}

export interface HomepageContent {
  hero: {
    name: { en: string; ru: string };
    titleAnimated: { en: string[]; ru: string[] };
    subtitleEn: string;
    subtitleRu: string;
    availableForWork: boolean;
  };
  about: {
    gifUrl?: string;
    bioEn: string;
    bioRu: string;
    bioFullEn?: string;
    bioFullRu?: string;
    stats: {
      projects: number;
      yearsExperience: number;
      certificates: number;
      githubStars: number;
      habrArticles: number;
    };
  };
  marqueeItems: string[];
  howIWork: HowIWorkConfig;
  contactCTA: ContactCTAConfig;
  achievements: Achievement[];
  socialLinks: SocialLink[];
  visibility?: HomepageVisibility;
}

export interface ContactFormField {
  id: string;
  type: 'text' | 'email' | 'textarea' | 'select';
  labelRu: string;
  labelEn: string;
  placeholderRu?: string;
  placeholderEn?: string;
  required: boolean;
  options?: Array<{ value: string; labelRu: string; labelEn: string }>;
  gridCol?: 1 | 2;
}

export interface ContactPageConfig {
  heading: { en: string; ru: string };
  subtitle: { en: string; ru: string };
  howIWork: { en: string; ru: string };
  formFields: ContactFormField[];
  submitTextRu: string;
  submitTextEn: string;
}
