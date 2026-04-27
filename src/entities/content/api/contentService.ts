import type { HomepageContent, HomepageVisibility, ContactPageConfig, HowIWorkConfig, ContactCTAConfig, Achievement, SocialLink } from '../model/types';
import { serverApi, ApiError } from '@/shared/api';

export interface ContentSection {
  section: string;
  data: unknown;
  updatedAt: string;
}

const DEFAULT_VISIBILITY: HomepageVisibility = {
  showAboutGif: true,
  showAboutBio: true,
  showStatProjects: true,
  showStatYears: true,
  showStatCertificates: true,
  showStatGithubStars: true,
  showStatHabrArticles: true,
  showProjectStars: true,
  showProjectStatus: true,
  showProjectTechStack: true,
  showMarquee: true,
  showHowIWork: true,
  showAchievements: true,
  showContactCTA: true,
};

function findSection<T>(sections: ContentSection[], name: string): T | null {
  const found = sections.find((s) => s.section === name);
  return found ? (found.data as T) : null;
}

const DEFAULT_HOMEPAGE: HomepageContent = {
  hero: {
    name: { en: 'Portfolio', ru: 'Портфолио' },
    titleAnimated: { en: ['Developer'], ru: ['Разработчик'] },
    subtitleEn: 'Welcome to my portfolio',
    subtitleRu: 'Добро пожаловать в моё портфолио',
    availableForWork: true,
  },
  about: {
    bioEn: '',
    bioRu: '',
    stats: { projects: 0, yearsExperience: 0, certificates: 0, githubStars: 0, habrArticles: 0 },
  },
  marqueeItems: [],
  howIWork: { steps: [], philosophyTitleRu: '', philosophyTitleEn: '', philosophyTextRu: '', philosophyTextEn: '', payment1Ru: '', payment1En: '', payment1DescRu: '', payment1DescEn: '', payment2Ru: '', payment2En: '', payment2DescRu: '', payment2DescEn: '' },
  contactCTA: { headingRu: '', headingEn: '', subtitleRu: '', subtitleEn: '', bgWordRu: '', bgWordEn: '' },
  achievements: [],
  socialLinks: [],
  visibility: DEFAULT_VISIBILITY,
};

const DEFAULT_CONTACT: ContactPageConfig = {
  heading: { en: 'Contact', ru: 'Контакты' },
  subtitle: { en: 'Get in touch', ru: 'Свяжитесь со мной' },
  howIWork: { en: '', ru: '' },
  formFields: [],
  submitTextRu: 'Отправить',
  submitTextEn: 'Send',
};

function extractHomepage(sections: ContentSection[]): HomepageContent {
  const hero = findSection<HomepageContent['hero']>(sections, 'hero');
  const about = findSection<HomepageContent['about']>(sections, 'about');
  const marqueeItems = findSection<string[]>(sections, 'marqueeItems');
  const howIWork = findSection<HowIWorkConfig>(sections, 'howIWork');
  const contactCTA = findSection<ContactCTAConfig>(sections, 'contactCTA');
  const visibility = findSection<HomepageVisibility>(sections, 'visibility');
  const achievements = findSection<Achievement[]>(sections, 'achievements');
  const socialLinks = findSection<SocialLink[]>(sections, 'socialLinks');

  return {
    hero: {
      ...DEFAULT_HOMEPAGE.hero,
      ...(hero ?? {}),
      titleAnimated: {
        ...DEFAULT_HOMEPAGE.hero.titleAnimated,
        ...((hero as Partial<HomepageContent['hero']> | null)?.titleAnimated ?? {}),
      },
    },
    about: about ?? DEFAULT_HOMEPAGE.about,
    marqueeItems: marqueeItems ?? DEFAULT_HOMEPAGE.marqueeItems,
    howIWork: howIWork ?? DEFAULT_HOMEPAGE.howIWork,
    contactCTA: contactCTA ?? DEFAULT_HOMEPAGE.contactCTA,
    achievements: achievements ?? DEFAULT_HOMEPAGE.achievements,
    socialLinks: socialLinks ?? DEFAULT_HOMEPAGE.socialLinks,
    visibility: { ...DEFAULT_VISIBILITY, ...visibility },
  };
}

function extractContact(sections: ContentSection[]): ContactPageConfig {
  const data = findSection<ContactPageConfig>(sections, 'contact');
  if (!data) return { ...DEFAULT_CONTACT };
  return data;
}

function isBuildTimeFetchFailure(error: unknown): boolean {
  if (error instanceof ApiError) return error.status >= 500;
  return true;
}

export const contentService = {
  async getAll(): Promise<ContentSection[]> {
    try {
      const data = await serverApi.get<ContentSection[]>('/api/content');
      return Array.isArray(data) ? data : [];
    } catch (error: unknown) {
      if (isBuildTimeFetchFailure(error)) return [];
      throw error;
    }
  },

  async getHomepage(): Promise<HomepageContent> {
    const sections = await this.getAll();
    return extractHomepage(sections);
  },

  async getContact(): Promise<ContactPageConfig> {
    const sections = await this.getAll();
    return extractContact(sections);
  },

  async getHomepageAndContact(): Promise<{ homepage: HomepageContent; contact: ContactPageConfig }> {
    const sections = await this.getAll();
    return {
      homepage: extractHomepage(sections),
      contact: extractContact(sections),
    };
  },

  update(token: string, section: string, data: unknown) {
    return serverApi
      .withToken(token)
      .put<ContentSection>(`/api/content/${section}`, data);
  },
};
