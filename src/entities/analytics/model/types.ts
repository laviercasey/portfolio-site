export type AnalyticsRange = '7d' | '30d';

export type AnalyticsMetric = 'pageviews' | 'visitors';

export interface AnalyticsSummary {
  range: AnalyticsRange;
  pageviews: number;
  uniqueVisitors: number;
  bounceRate: number;
  avgSessionSeconds: number;
  deltaPageviews: number;
  previousPageviews: number;
}

export interface TopItem {
  key: string;
  views: number;
  uniques?: number;
}

export interface TopPage {
  path: string;
  views: number;
  uniques: number;
}

export interface TopReferrer {
  referrer: string;
  views: number;
}

export interface TopCountry {
  country: string;
  views: number;
}

export type UtmType = 'source' | 'medium' | 'campaign';

export interface TopUtm {
  value: string;
  views: number;
}

export interface TimePoint {
  date: string;
  value: number;
}
