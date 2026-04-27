'use client';

import { useEffect } from 'react';
import { useLocale } from 'next-intl';

interface ErrorProps {
  error: Error & { digest?: string };
  reset: () => void;
}

export default function LocaleError({ error, reset }: ErrorProps) {
  const locale = useLocale();
  const isRu = locale === 'ru';

  useEffect(() => {
    if (typeof window !== 'undefined' && process.env.NODE_ENV !== 'production') {
      console.error(error);
    }
  }, [error]);

  return (
    <div className="container py-32 text-center">
      <h1 className="font-heading text-4xl font-bold mb-4 gradient-text">
        {isRu ? 'Что-то пошло не так' : 'Something went wrong'}
      </h1>
      <p className="text-muted-foreground mb-8 max-w-md mx-auto">
        {isRu
          ? 'Не удалось загрузить данные. Попробуйте обновить страницу через минуту.'
          : 'Failed to load content. Please retry in a moment.'}
      </p>
      <button
        type="button"
        onClick={() => reset()}
        className="px-6 py-3 rounded-xl border border-border/60 hover:border-border transition-colors"
      >
        {isRu ? 'Повторить' : 'Try again'}
      </button>
      {error.digest && (
        <p className="mt-8 text-xs text-muted-foreground/60 font-mono">{error.digest}</p>
      )}
    </div>
  );
}
