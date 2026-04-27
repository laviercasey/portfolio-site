import Link from 'next/link';
import { getTranslations } from 'next-intl/server';
import {
  Activity,
  Eye,
  Globe2,
  LinkIcon,
  MousePointer2,
  TrendingDown,
  Users,
  Timer,
} from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui';
import { cn } from '@/shared/lib';
import { analyticsService } from '@/entities/analytics';
import type { AnalyticsRange, TopItem, TopUtm } from '@/entities/analytics';
import { AnalyticsKpiCard } from './AnalyticsKpiCard';
import { TopList } from './TopList';

export interface AdminAnalyticsProps {
  token: string;
  pathname: string;
  range?: AnalyticsRange;
}

const DEFAULT_RANGE: AnalyticsRange = '30d';
const TOP_LIMIT = 8;

function isAnalyticsRange(value: unknown): value is AnalyticsRange {
  return value === '7d' || value === '30d';
}

function formatNumber(n: number): string {
  return new Intl.NumberFormat('en-US').format(n);
}

function formatDuration(seconds: number): string {
  if (!Number.isFinite(seconds) || seconds <= 0) return '0s';
  const m = Math.floor(seconds / 60);
  const s = Math.round(seconds % 60);
  if (m === 0) return `${s}s`;
  return `${m}m ${s}s`;
}

function formatPercent(fraction: number): string {
  if (!Number.isFinite(fraction)) return '0%';
  return `${(fraction * 100).toFixed(1)}%`;
}

function countryFlag(code: string): string {
  if (code.length !== 2 || !/^[A-Z]{2}$/.test(code.toUpperCase())) return '';
  const cc = code.toUpperCase();
  const base = 0x1f1e6;
  const a = 'A'.charCodeAt(0);
  return String.fromCodePoint(base + (cc.charCodeAt(0) - a), base + (cc.charCodeAt(1) - a));
}

export async function AdminAnalytics({ token, pathname, range }: AdminAnalyticsProps) {
  const selected: AnalyticsRange = isAnalyticsRange(range) ? range : DEFAULT_RANGE;
  const t = await getTranslations('analytics');

  const [summary, pages, referrers, countries, utmSource, utmMedium, utmCampaign] =
    await Promise.all([
      analyticsService.summary(token, selected),
      analyticsService.topPages(token, selected, TOP_LIMIT),
      analyticsService.topReferrers(token, selected, TOP_LIMIT),
      analyticsService.topCountries(token, selected, TOP_LIMIT),
      analyticsService.topUtm(token, selected, 'source', TOP_LIMIT),
      analyticsService.topUtm(token, selected, 'medium', TOP_LIMIT),
      analyticsService.topUtm(token, selected, 'campaign', TOP_LIMIT),
    ]);

  if (summary === null) {
    return (
      <Card>
        <CardHeader>
          <CardTitle className="text-xl flex items-center gap-2">
            <Activity className="h-5 w-5 text-muted-foreground" aria-hidden="true" />
            {t('disabled.title')}
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-3 text-sm text-muted-foreground">
          <p>{t('disabled.description')}</p>
          <p className="font-mono text-xs bg-muted/50 p-3 rounded border">
            {t('disabled.envHint')}
          </p>
        </CardContent>
      </Card>
    );
  }

  const deltaLabel = t('kpi.vsPrevious');

  const pageItems: TopItem[] = pages.map((p) => ({
    key: p.path || '/',
    views: p.views,
    uniques: p.uniques,
  }));
  const referrerItems: TopItem[] = referrers.map((r) => ({
    key: r.referrer,
    views: r.views,
  }));
  const countryItems: TopItem[] = countries.map((c) => ({
    key: c.country,
    views: c.views,
  }));
  const toUtmItems = (rows: readonly TopUtm[]): TopItem[] =>
    rows.map((u) => ({ key: u.value, views: u.views }));
  const utmSourceItems = toUtmItems(utmSource);
  const utmMediumItems = toUtmItems(utmMedium);
  const utmCampaignItems = toUtmItems(utmCampaign);

  return (
    <section aria-labelledby="admin-analytics-heading" className="space-y-4">
      <header className="flex flex-wrap items-center justify-between gap-3">
        <h2 id="admin-analytics-heading" className="text-2xl font-semibold">
          {t('title')}
        </h2>
        <div
          className="inline-flex items-center rounded-md border p-1 text-sm"
          role="tablist"
          aria-label={t('title')}
        >
          <RangeTab href={pathname} range="7d" label={t('range.7d')} active={selected === '7d'} />
          <RangeTab
            href={pathname}
            range="30d"
            label={t('range.30d')}
            active={selected === '30d'}
          />
        </div>
      </header>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <AnalyticsKpiCard
          label={t('kpi.pageviews')}
          value={formatNumber(summary.pageviews)}
          icon={Eye}
          delta={summary.deltaPageviews}
          deltaLabel={deltaLabel}
        />
        <AnalyticsKpiCard
          label={t('kpi.visitors')}
          value={formatNumber(summary.uniqueVisitors)}
          icon={Users}
        />
        <AnalyticsKpiCard
          label={t('kpi.session')}
          value={formatDuration(summary.avgSessionSeconds)}
          icon={Timer}
        />
        <AnalyticsKpiCard
          label={t('kpi.bounce')}
          value={formatPercent(summary.bounceRate)}
          icon={TrendingDown}
          higherIsBetter={false}
        />
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <TopList
          title={t('topPages')}
          items={pageItems}
          icon={MousePointer2}
          emptyLabel={t('empty')}
        />
        <TopList
          title={t('topReferrers')}
          items={referrerItems}
          icon={LinkIcon}
          emptyLabel={t('empty')}
        />
      </div>

      {countryItems.length > 0 ? (
        <TopList
          title={t('topCountries')}
          items={countryItems}
          icon={Globe2}
          emptyLabel={t('empty')}
          renderKey={(code) => {
            const flag = countryFlag(code);
            return (
              <span className="inline-flex items-center gap-2">
                {flag ? <span aria-hidden="true">{flag}</span> : null}
                <span>{code}</span>
              </span>
            );
          }}
        />
      ) : null}

      <section className="space-y-3">
        <h3 className="text-lg font-semibold">{t('utm.sectionTitle')}</h3>
        <p className="text-sm text-muted-foreground">{t('utm.description')}</p>
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
          <TopList
            title={t('utm.source')}
            items={utmSourceItems}
            emptyLabel={t('empty')}
          />
          <TopList
            title={t('utm.medium')}
            items={utmMediumItems}
            emptyLabel={t('empty')}
          />
          <TopList
            title={t('utm.campaign')}
            items={utmCampaignItems}
            emptyLabel={t('empty')}
          />
        </div>
      </section>
    </section>
  );
}

interface RangeTabProps {
  href: string;
  range: AnalyticsRange;
  label: string;
  active: boolean;
}

function RangeTab({ href, range, label, active }: RangeTabProps) {
  return (
    <Link
      href={`${href}?range=${range}`}
      role="tab"
      aria-selected={active}
      prefetch={false}
      className={cn(
        'px-3 py-1.5 rounded-sm transition-colors',
        active
          ? 'bg-primary text-primary-foreground'
          : 'text-muted-foreground hover:text-foreground hover:bg-muted',
      )}
    >
      {label}
    </Link>
  );
}

export function AnalyticsSkeleton() {
  return (
    <section aria-hidden="true" className="space-y-4 animate-pulse">
      <div className="flex items-center justify-between gap-3">
        <div className="h-7 w-32 rounded bg-muted" />
        <div className="h-8 w-32 rounded-md border bg-muted/40" />
      </div>
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {Array.from({ length: 4 }).map((_, i) => (
          <Card key={i}>
            <CardContent className="p-5 space-y-3">
              <div className="h-4 w-24 rounded bg-muted" />
              <div className="h-8 w-20 rounded bg-muted" />
              <div className="h-4 w-16 rounded bg-muted/60" />
            </CardContent>
          </Card>
        ))}
      </div>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {Array.from({ length: 2 }).map((_, i) => (
          <Card key={i}>
            <CardContent className="p-5 space-y-2">
              <div className="h-5 w-32 rounded bg-muted" />
              {Array.from({ length: 5 }).map((__, j) => (
                <div key={j} className="h-4 w-full rounded bg-muted/60" />
              ))}
            </CardContent>
          </Card>
        ))}
      </div>
    </section>
  );
}
