'use client';

import { useEffect, useState, useCallback } from 'react';
import { ThemeContext, type Theme } from '@/shared/lib';

function setThemeClass(theme: Theme) {
  if (theme === 'dark') document.documentElement.classList.add('dark');
  else document.documentElement.classList.remove('dark');
}

function animateThemeWipe(oldBg: string) {
  if (window.matchMedia('(prefers-reduced-motion: reduce)').matches) return;

  const overlay = document.createElement('div');
  overlay.style.cssText =
    'position:fixed;inset:0;z-index:99999;pointer-events:none;will-change:transform;background:' + oldBg;
  document.body.appendChild(overlay);

  requestAnimationFrame(() => {
    overlay.style.transition = 'transform 0.8s cubic-bezier(0.76, 0, 0.24, 1)';
    overlay.style.transform  = 'translateX(100%)';
    const cleanup = () => overlay.remove();
    overlay.addEventListener('transitionend', cleanup, { once: true });
    setTimeout(cleanup, 1000);
  });
}

interface ThemeProviderProps {
  children: React.ReactNode;
}

export function ThemeProvider({ children }: ThemeProviderProps) {
  const [theme, setThemeState] = useState<Theme>('dark');
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    const stored = localStorage.getItem('theme') as Theme | null;
    const initial: Theme = stored === 'light' ? 'light' : 'dark';
    setThemeState(initial);
    setThemeClass(initial);
    setMounted(true);
  }, []);

  const setTheme = useCallback((t: Theme) => {
    localStorage.setItem('theme', t);
    const oldBg = getComputedStyle(document.body).backgroundColor;
    setThemeState(t);
    setThemeClass(t);
    animateThemeWipe(oldBg);
  }, []);

  if (!mounted) {
    return (
      <ThemeContext.Provider value={{ theme: 'dark', resolvedTheme: 'dark', setTheme: () => {} }}>
        {children}
      </ThemeContext.Provider>
    );
  }

  return (
    <ThemeContext.Provider value={{ theme, resolvedTheme: theme, setTheme }}>
      {children}
    </ThemeContext.Provider>
  );
}
