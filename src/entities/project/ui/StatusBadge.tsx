import { useLocale } from 'next-intl';
import { Badge } from '@/shared/ui';
import type { ProjectStatus } from '../model/types';

const statusVariants: Record<ProjectStatus, 'success' | 'warning' | 'info'> = {
  completed: 'success',
  in_development: 'warning',
  open_source: 'info',
};

interface StatusBadgeProps {
  status: ProjectStatus;
}

export default function StatusBadge({ status }: StatusBadgeProps) {
  const locale = useLocale();
  const labels: Record<ProjectStatus, Record<string, string>> = {
    completed: { ru: 'Завершён', en: 'Completed' },
    in_development: { ru: 'В разработке', en: 'In Development' },
    open_source: { ru: 'Open Source', en: 'Open Source' },
  };
  return (
    <Badge variant={statusVariants[status]}>
      {labels[status][locale] ?? labels[status].en}
    </Badge>
  );
}
