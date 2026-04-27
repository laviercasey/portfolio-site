import { serverApi, ApiError } from '@/shared/api';
import type {
  AnalyticsMetric,
  AnalyticsRange,
  AnalyticsSummary,
  TimePoint,
  TopCountry,
  TopPage,
  TopReferrer,
  TopUtm,
  UtmType,
} from '../model/types';

interface AnalyticsEnvelope<T> {
  success: boolean;
  data?: T;
  error?: string;
}

const DISABLED_ERROR = 'analytics_not_configured';

function isDisabled(error: unknown): boolean {
  return error instanceof ApiError && error.status === 503 && error.message === DISABLED_ERROR;
}

function unwrap<T>(envelope: AnalyticsEnvelope<T>): T | undefined {
  if (!envelope.success) return undefined;
  return envelope.data;
}

function buildQuery(params: Record<string, string | number | undefined>): string {
  const entries = Object.entries(params).filter(([, v]) => v !== undefined);
  if (entries.length === 0) return '';
  const usp = new URLSearchParams();
  for (const [k, v] of entries) usp.set(k, String(v));
  return `?${usp.toString()}`;
}

export const analyticsService = {
  async summary(token: string, range: AnalyticsRange = '30d'): Promise<AnalyticsSummary | null> {
    try {
      const envelope = await serverApi
        .withToken(token)
        .get<AnalyticsEnvelope<AnalyticsSummary>>(`/api/analytics/summary${buildQuery({ range })}`);
      return unwrap(envelope) ?? null;
    } catch (error: unknown) {
      if (isDisabled(error)) return null;
      throw error;
    }
  },

  async topPages(
    token: string,
    range: AnalyticsRange = '30d',
    limit = 10,
  ): Promise<TopPage[]> {
    try {
      const envelope = await serverApi
        .withToken(token)
        .get<AnalyticsEnvelope<TopPage[]>>(`/api/analytics/top-pages${buildQuery({ range, limit })}`);
      return unwrap(envelope) ?? [];
    } catch (error: unknown) {
      if (isDisabled(error)) return [];
      throw error;
    }
  },

  async topReferrers(
    token: string,
    range: AnalyticsRange = '30d',
    limit = 10,
  ): Promise<TopReferrer[]> {
    try {
      const envelope = await serverApi
        .withToken(token)
        .get<AnalyticsEnvelope<TopReferrer[]>>(
          `/api/analytics/top-referrers${buildQuery({ range, limit })}`,
        );
      return unwrap(envelope) ?? [];
    } catch (error: unknown) {
      if (isDisabled(error)) return [];
      throw error;
    }
  },

  async topCountries(
    token: string,
    range: AnalyticsRange = '30d',
    limit = 10,
  ): Promise<TopCountry[]> {
    try {
      const envelope = await serverApi
        .withToken(token)
        .get<AnalyticsEnvelope<TopCountry[]>>(
          `/api/analytics/top-countries${buildQuery({ range, limit })}`,
        );
      return unwrap(envelope) ?? [];
    } catch (error: unknown) {
      if (isDisabled(error)) return [];
      throw error;
    }
  },

  async topUtm(
    token: string,
    range: AnalyticsRange = '30d',
    type: UtmType,
    limit = 10,
  ): Promise<TopUtm[]> {
    try {
      const envelope = await serverApi
        .withToken(token)
        .get<AnalyticsEnvelope<TopUtm[]>>(
          `/api/analytics/top-utm${buildQuery({ type, range, limit })}`,
        );
      return unwrap(envelope) ?? [];
    } catch (error: unknown) {
      if (isDisabled(error)) return [];
      throw error;
    }
  },

  async timeseries(
    token: string,
    range: AnalyticsRange = '30d',
    metric: AnalyticsMetric = 'pageviews',
  ): Promise<TimePoint[]> {
    try {
      const envelope = await serverApi
        .withToken(token)
        .get<AnalyticsEnvelope<TimePoint[]>>(
          `/api/analytics/timeseries${buildQuery({ range, metric })}`,
        );
      return unwrap(envelope) ?? [];
    } catch (error: unknown) {
      if (isDisabled(error)) return [];
      throw error;
    }
  },
};
