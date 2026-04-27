import { getTranslations } from 'next-intl/server';
import { Link } from '@/shared/config';
import { Button } from '@/shared/ui';

export default async function NotFound() {
  const t = await getTranslations('notFound');
  return (
    <div className="container py-32 text-center">
      <h1 className="text-8xl font-bold text-muted-foreground/30 mb-4">404</h1>
      <h2 className="text-2xl font-bold mb-4">{t('title')}</h2>
      <p className="text-muted-foreground mb-8">{t('description')}</p>
      <Button asChild>
        <Link href="/">{t('goHome')}</Link>
      </Button>
    </div>
  );
}
