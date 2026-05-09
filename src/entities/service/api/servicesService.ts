import type {
  Service,
  ServiceFaq,
  ServiceProcessStep,
  ServicesPageData,
  CreateServiceInput,
  UpdateServiceInput,
  CreateServiceFaqInput,
  UpdateServiceFaqInput,
  CreateServiceProcessStepInput,
  UpdateServiceProcessStepInput,
} from '../model/types';
import { serverApi } from '@/shared/api';

const EMPTY_PAGE: ServicesPageData = {
  services: [],
  faqs: [],
  processSteps: [],
};

export const servicesService = {
  async getPageData(): Promise<ServicesPageData> {
    try {
      const data = await serverApi.get<ServicesPageData>('/api/services');
      return {
        services: data?.services ?? [],
        faqs: data?.faqs ?? [],
        processSteps: data?.processSteps ?? [],
      };
    } catch {
      return EMPTY_PAGE;
    }
  },

  listServices(): Promise<Service[]> {
    return serverApi.get<Service[]>('/api/services/list').then((d) => d ?? []);
  },

  listFaqs(): Promise<ServiceFaq[]> {
    return serverApi.get<ServiceFaq[]>('/api/services/faqs').then((d) => d ?? []);
  },

  listProcessSteps(): Promise<ServiceProcessStep[]> {
    return serverApi.get<ServiceProcessStep[]>('/api/services/process').then((d) => d ?? []);
  },

  createService(token: string, input: CreateServiceInput) {
    return serverApi.withToken(token).post<Service>('/api/services', input);
  },

  updateService(token: string, id: string, input: UpdateServiceInput) {
    return serverApi.withToken(token).put<Service>(`/api/services/${id}`, input);
  },

  deleteService(token: string, id: string) {
    return serverApi.withToken(token).delete(`/api/services/${id}`);
  },

  createFaq(token: string, input: CreateServiceFaqInput) {
    return serverApi.withToken(token).post<ServiceFaq>('/api/services/faqs', input);
  },

  updateFaq(token: string, id: string, input: UpdateServiceFaqInput) {
    return serverApi.withToken(token).put<ServiceFaq>(`/api/services/faqs/${id}`, input);
  },

  deleteFaq(token: string, id: string) {
    return serverApi.withToken(token).delete(`/api/services/faqs/${id}`);
  },

  createProcessStep(token: string, input: CreateServiceProcessStepInput) {
    return serverApi.withToken(token).post<ServiceProcessStep>('/api/services/process', input);
  },

  updateProcessStep(token: string, id: string, input: UpdateServiceProcessStepInput) {
    return serverApi.withToken(token).put<ServiceProcessStep>(`/api/services/process/${id}`, input);
  },

  deleteProcessStep(token: string, id: string) {
    return serverApi.withToken(token).delete(`/api/services/process/${id}`);
  },
};
