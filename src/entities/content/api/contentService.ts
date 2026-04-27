import type { HomepageContent, HomepageVisibility, ContactPageConfig, HowIWorkConfig, ContactCTAConfig, Achievement, SocialLink } from '../model/types';
import { serverApi } from '@/shared/api';

export interface ContentSection {
  section: string;
  data: unknown;
  updatedAt: string;
}

const BASE_VISIBILITY: HomepageVisibility = {
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
    hero: hero ?? {
      name: { en: '', ru: '' },
      titleAnimated: { en: [], ru: [] },
      subtitleEn: '',
      subtitleRu: '',
      availableForWork: false,
    },
    about: about ?? {
      bioEn: '',
      bioRu: '',
      stats: { projects: 0, yearsExperience: 0, certificates: 0, githubStars: 0, habrArticles: 0 },
    },
    marqueeItems: marqueeItems ?? [],
    howIWork: howIWork ?? { steps: [], philosophyTitleRu: '', philosophyTitleEn: '', philosophyTextRu: '', philosophyTextEn: '', payment1Ru: '', payment1En: '', payment1DescRu: '', payment1DescEn: '', payment2Ru: '', payment2En: '', payment2DescRu: '', payment2DescEn: '' },
    contactCTA: contactCTA ?? { headingRu: '', headingEn: '', subtitleRu: '', subtitleEn: '', bgWordRu: '', bgWordEn: '' },
    achievements: achievements ?? [],
    socialLinks: socialLinks ?? [],
    visibility: { ...BASE_VISIBILITY, ...visibility },
  };
}

function extractContact(sections: ContentSection[]): ContactPageConfig {
  const data = findSection<ContactPageConfig>(sections, 'contact');
  if (data) return data;
  return {
    heading: { en: '', ru: '' },
    subtitle: { en: '', ru: '' },
    howIWork: { en: '', ru: '' },
    formFields: [],
    submitTextRu: '',
    submitTextEn: '',
  };
}

export const contentService = {
  async getAll(): Promise<ContentSection[]> {
    const data = await serverApi.get<ContentSection[]>('/api/content');
    return Array.isArray(data) ? data : [];
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
