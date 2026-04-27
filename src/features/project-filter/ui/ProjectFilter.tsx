'use client';

import { useLocale } from 'next-intl';
import { useRouter, useSearchParams, usePathname } from 'next/navigation';
import { Button } from '@/shared/ui';

const categories = [
  { value: 'all', ru: 'Все', en: 'All' },
  { value: 'web', ru: 'Веб', en: 'Web' },
  { value: 'mobile', ru: 'Мобильные', en: 'Mobile' },
  { value: 'data', ru: 'Data', en: 'Data' },
  { value: 'research', ru: 'Исследования', en: 'Research' },
] as const;

export default function ProjectFilter() {
  const locale = useLocale();
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const currentCategory = searchParams.get('category') ?? 'all';

  const updateFilter = (category: string) => {
    const params = new URLSearchParams(searchParams.toString());
    if (category === 'all') {
      params.delete('category');
    } else {
      params.set('category', category);
    }
    const query = params.toString();
    router.push(query ? `${pathname}?${query}` : pathname);
  };

  return (
    <div className="flex flex-wrap gap-2 mb-8">
      {categories.map((cat) => (
        <Button
          key={cat.value}
          variant={currentCategory === cat.value ? 'default' : 'outline'}
          size="sm"
          onClick={() => updateFilter(cat.value)}
        >
          {locale === 'ru' ? cat.ru : cat.en}
        </Button>
      ))}
    </div>
  );
}
