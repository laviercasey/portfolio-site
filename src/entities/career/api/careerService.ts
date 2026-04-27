import type { CareerContent } from '../model/types';
import { serverApi, ApiError } from '@/shared/api';

const EMPTY_CAREER: CareerContent = {
  education: [],
  workHistory: [],
  certificates: [],
  publications: [],
};

function isBuildTimeFetchFailure(error: unknown): boolean {
  if (error instanceof ApiError) return error.status >= 500;
  return true;
}

function normalizeCareer(raw: Partial<CareerContent> | null | undefined): CareerContent {
  if (!raw) return { ...EMPTY_CAREER };
  return {
    education: Array.isArray(raw.education) ? raw.education : [],
    workHistory: Array.isArray(raw.workHistory) ? raw.workHistory : [],
    certificates: Array.isArray(raw.certificates) ? raw.certificates : [],
    publications: Array.isArray(raw.publications) ? raw.publications : [],
  };
}

export const careerService = {
  async getAll(): Promise<CareerContent> {
    try {
      const data = await serverApi.get<Partial<CareerContent>>('/api/career');
      return normalizeCareer(data);
    } catch (error: unknown) {
      if (isBuildTimeFetchFailure(error)) return { ...EMPTY_CAREER };
      throw error;
    }
  },

  create(token: string, type: string, data: unknown) {
    return serverApi.withToken(token).post(`/api/career/${type}`, data);
  },

  update(token: string, type: string, id: string, data: unknown) {
    return serverApi.withToken(token).put(`/api/career/${type}/${id}`, data);
  },

  delete(token: string, type: string, id: string) {
    return serverApi.withToken(token).delete(`/api/career/${type}/${id}`);
  },
};
