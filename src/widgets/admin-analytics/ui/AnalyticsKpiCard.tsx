import type { LucideIcon } from 'lucide-react';
import { TrendingDown, TrendingUp, Minus } from 'lucide-react';
import { Card, CardContent, Badge } from '@/shared/ui';
import { cn } from '@/shared/lib';

export interface AnalyticsKpiCardProps {
  label: string;
  value: string;
  icon: LucideIcon;
  delta?: number;
  deltaLabel?: string;
  higherIsBetter?: boolean;
}

function formatDeltaPercent(delta: number): string {
  const pct = delta * 100;
  const sign = pct > 0 ? '+' : '';
  return `${sign}${pct.toFixed(1)}%`;
}

function pickDeltaVariant(
  delta: number,
  higherIsBetter: boolean,
): 'success' | 'destructive' | 'secondary' {
  if (Math.abs(delta) < 0.005) return 'secondary';
  const isUp = delta > 0;
  const isGood = higherIsBetter ? isUp : !isUp;
  return isGood ? 'success' : 'destructive';
}

function DeltaIcon({ delta }: { delta: number }) {
  if (Math.abs(delta) < 0.005) return <Minus className="h-3 w-3" aria-hidden="true" />;
  if (delta > 0) return <TrendingUp className="h-3 w-3" aria-hidden="true" />;
  return <TrendingDown className="h-3 w-3" aria-hidden="true" />;
}

export function AnalyticsKpiCard({
  label,
  value,
  icon: Icon,
  delta,
  deltaLabel,
  higherIsBetter = true,
}: AnalyticsKpiCardProps) {
  const hasDelta = typeof delta === 'number' && Number.isFinite(delta);

  return (
    <Card>
      <CardContent className="p-5">
        <div className="flex items-start justify-between gap-3">
          <div className="min-w-0">
            <p className="text-sm text-muted-foreground">{label}</p>
            <p className="text-3xl font-bold truncate">{value}</p>
          </div>
          <Icon className="h-8 w-8 shrink-0 text-primary/30" aria-hidden="true" />
        </div>
        {hasDelta ? (
          <div className="mt-3 flex items-center gap-2">
            <Badge
              variant={pickDeltaVariant(delta, higherIsBetter)}
              className={cn('gap-1 text-xs font-medium')}
            >
              <DeltaIcon delta={delta} />
              {formatDeltaPercent(delta)}
            </Badge>
            {deltaLabel ? (
              <span className="text-xs text-muted-foreground">{deltaLabel}</span>
            ) : null}
          </div>
        ) : null}
      </CardContent>
    </Card>
  );
}
