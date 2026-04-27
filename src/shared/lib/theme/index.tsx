'use client';

import { createContext, useContext } from 'react';

export type Theme = 'light' | 'dark';

export interface ThemeContextValue {
  theme: Theme;
  resolvedTheme: Theme;
  setTheme: (t: Theme) => void;
}

export const ThemeContext = createContext<ThemeContextValue>({
  theme: 'dark',
  resolvedTheme: 'dark',
  setTheme: () => {},
});

export function useTheme() {
  return useContext(ThemeContext);
}

export const themeScript = `
(function(){
  try {
    var t = localStorage.getItem('theme');
    if (t !== 'light') document.documentElement.classList.add('dark');
  } catch(e) { document.documentElement.classList.add('dark'); }
})();
`.trim();
