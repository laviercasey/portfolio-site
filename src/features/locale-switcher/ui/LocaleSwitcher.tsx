'use client';

import { useEffect, useState } from 'react';
import { motion } from 'framer-motion';
import { useLocale } from 'next-intl';
import { usePathname } from '@/shared/config';
import { cn } from '@/shared/lib';
import LocaleTransitionOverlay from './LocaleTransitionOverlay';

const locales = ['ru', 'en'] as const;

export default function LocaleSwitcher() {
  const locale   = useLocale();
  const pathname = usePathname();
  const [pending, setPending] = useState<{ locale: string; pathname: string } | null>(null);

  useEffect(() => {
    setPending(null);
  }, [locale]);

  const handleSwitch = (loc: typeof locales[number]) => {
    if (loc === locale || pending) return;
    setPending({ locale: loc, pathname });
  };

  return (
    <>
      <div className="relative flex items-center rounded-xl border border-white/10 p-0.5">
        {locales.map((loc) => (
          <button
            key={loc}
            onClick={() => handleSwitch(loc)}
            disabled={!!pending}
            className={cn(
              'relative flex items-center justify-center h-7 w-9',
              'text-xs font-mono font-semibold rounded-lg z-10',
              'transition-colors duration-200',
              locale === loc || pending?.locale === loc
                ? 'text-primary-foreground pointer-events-none'
                : 'text-muted-foreground hover:text-foreground',
            )}
          >
            {(locale === loc || pending?.locale === loc) && (
              <motion.span
                layoutId="locale-pill"
                className="absolute inset-0 rounded-lg bg-primary -z-10"
                transition={{ type: 'spring', stiffness: 400, damping: 32 }}
              />
            )}
            {loc.toUpperCase()}
          </button>
        ))}
      </div>

      {pending && (
        <LocaleTransitionOverlay
          targetLocale={pending.locale}
          targetPathname={pending.pathname}
        />
      )}
    </>
  );
}
