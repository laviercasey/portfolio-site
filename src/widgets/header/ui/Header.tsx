'use client';

import { useState, useEffect } from 'react';
import { useTranslations, useLocale } from 'next-intl';
import { siteConfig as siteJson } from '@/shared/config';
import type { SiteConfig } from '@/shared/types';

const site = siteJson as SiteConfig;
import { Link } from '@/shared/config';
import { Menu, X, Code2 } from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';
import { ThemeToggle } from '@/features/theme-toggle';
import { LocaleSwitcher } from '@/features/locale-switcher';
import { cn } from '@/shared/lib';

const navLinks = [
  { key: 'home', href: '/' },
  { key: 'projects', href: '/projects' },
  { key: 'career', href: '/career' },
  { key: 'contact', href: '/contact' },
] as const;

export default function Header() {
  const t = useTranslations('nav');
  const locale = useLocale();
  const siteName = locale === 'ru' ? (site.nameRu ?? site.name) : site.name;
  const [mobileOpen, setMobileOpen] = useState(false);
  const [scrolled, setScrolled] = useState(false);

  useEffect(() => {
    const handler = () => setScrolled(window.scrollY > 20);
    window.addEventListener('scroll', handler, { passive: true });
    return () => window.removeEventListener('scroll', handler);
  }, []);

  return (
    <header className="fixed top-4 left-4 right-4 z-50 flex flex-col items-center pointer-events-none">
      <div
        className={cn(
          'w-full max-w-5xl pointer-events-auto',
          'rounded-2xl border transition-all duration-300',
          scrolled
            ? 'bg-background/80 backdrop-blur-2xl border-white/15 shadow-lg shadow-black/10'
            : 'bg-background/60 backdrop-blur-xl border-white/10',
        )}
      >
        <div className="flex h-14 items-center justify-between px-4 md:px-6">
          <Link
            href="/"
            className="flex items-center gap-2 font-bold text-foreground hover:text-primary transition-colors"
          >
            <div className="flex items-center justify-center w-7 h-7 rounded-lg bg-primary/10 border border-primary/20">
              <Code2 className="h-4 w-4 text-primary" strokeWidth={2} />
            </div>
            <span className="font-heading text-sm tracking-tight">{siteName}</span>
          </Link>

          <nav className="hidden md:flex items-center gap-1">
            {navLinks.map(({ key, href }) => (
              <Link
                key={key}
                href={href}
                className="px-3 py-1.5 text-sm font-medium text-muted-foreground hover:text-foreground rounded-lg hover:bg-white/8 transition-all duration-200"
              >
                {t(key)}
              </Link>
            ))}
          </nav>

          <div className="flex items-center gap-2">
            <ThemeToggle />
            <LocaleSwitcher />

            <button
              className={cn(
                'md:hidden flex items-center justify-center w-9 h-9 rounded-xl',
                'text-muted-foreground hover:text-foreground',
                'hover:bg-white/10 border border-transparent hover:border-white/15',
                'transition-all duration-200'
              )}
              onClick={() => setMobileOpen(!mobileOpen)}
              aria-label="Toggle menu"
            >
              {mobileOpen ? <X className="h-4 w-4" /> : <Menu className="h-4 w-4" />}
            </button>
          </div>
        </div>

        <AnimatePresence>
          {mobileOpen && (
            <motion.div
              key="mobile-menu"
              initial={{ opacity: 0, height: 0 }}
              animate={{ opacity: 1, height: 'auto' }}
              exit={{ opacity: 0, height: 0 }}
              transition={{ duration: 0.2 }}
              className="overflow-hidden border-t border-white/10 md:hidden"
            >
              <nav className="flex flex-col gap-1 p-3">
                {navLinks.map(({ key, href }) => (
                  <Link
                    key={key}
                    href={href}
                    className="px-3 py-2 text-sm font-medium text-muted-foreground hover:text-foreground rounded-lg hover:bg-white/8 transition-all duration-200"
                    onClick={() => setMobileOpen(false)}
                  >
                    {t(key)}
                  </Link>
                ))}
              </nav>
            </motion.div>
          )}
        </AnimatePresence>
      </div>
    </header>
  );
}
