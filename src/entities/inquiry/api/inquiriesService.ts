import type { Inquiry } from '../model/types';
import { serverApi } from '@/shared/api';

export const inquiriesService = {
  list(token: string, status?: string) {
    const q = status ? `?status=${status}` : '';
    return serverApi.withToken(token).get<Inquiry[]>(`/api/inquiries${q}`);
  },

  getById(token: string, id: string) {
    return serverApi.withToken(token).get<Inquiry>(`/api/inquiries/${id}`);
  },

  create(data: {
    name: string;
    email: string;
    company?: string;
    type: string;
    budget?: string;
    message: string;
  }) {
    return serverApi.post<Inquiry>('/api/inquiries', data);
  },

  updateStatus(token: string, id: string, status: string, adminNotes?: string) {
    return serverApi
      .withToken(token)
      .patch<Inquiry>(`/api/inquiries/${id}`, { status, adminNotes });
  },
};
