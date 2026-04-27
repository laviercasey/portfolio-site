import type { CareerContent } from '../model/types';
import { serverApi } from '@/shared/api';

function normalizeCareer(raw: Partial<CareerContent> | null | undefined): CareerContent {
  if (!raw) {
    return { education: [], workHistory: [], certificates: [], publications: [] };
  }
  return {
    education: Array.isArray(raw.education) ? raw.education : [],
    workHistory: Array.isArray(raw.workHistory) ? raw.workHistory : [],
    certificates: Array.isArray(raw.certificates) ? raw.certificates : [],
    publications: Array.isArray(raw.publications) ? raw.publications : [],
  };
}

export const careerService = {
  async getAll(): Promise<CareerContent> {
    const data = await serverApi.get<Partial<CareerContent>>('/api/career');
    return normalizeCareer(data);
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
