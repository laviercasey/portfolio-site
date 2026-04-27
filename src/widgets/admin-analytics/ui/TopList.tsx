import type { LucideIcon } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui';
import type { TopItem } from '@/entities/analytics';

export interface TopListProps {
  title: string;
  items: readonly TopItem[];
  icon?: LucideIcon;
  emptyLabel: string;
  renderKey?: (key: string) => React.ReactNode;
}

function formatViews(n: number): string {
  return new Intl.NumberFormat('en-US').format(n);
}

export function TopList({
  title,
  items,
  icon: Icon,
  emptyLabel,
  renderKey,
}: TopListProps) {
  return (
    <Card>
      <CardHeader className="pb-3">
        <CardTitle className="text-lg flex items-center gap-2">
          {Icon ? <Icon className="h-5 w-5 text-muted-foreground" aria-hidden="true" /> : null}
          {title}
        </CardTitle>
      </CardHeader>
      <CardContent>
        {items.length === 0 ? (
          <p className="text-muted-foreground text-sm py-4 text-center">{emptyLabel}</p>
        ) : (
          <ul className="space-y-2">
            {items.map((item) => (
              <li
                key={item.key}
                className="flex items-center justify-between gap-3 text-sm"
              >
                <span className="min-w-0 truncate text-foreground/90" title={item.key}>
                  {renderKey ? renderKey(item.key) : item.key}
                </span>
                <span className="shrink-0 tabular-nums text-muted-foreground">
                  {formatViews(item.views)}
                </span>
              </li>
            ))}
          </ul>
        )}
      </CardContent>
    </Card>
  );
}
