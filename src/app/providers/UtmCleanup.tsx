'use client';

import { useEffect } from 'react';

const UTM_PARAM_RE = /^utm_/i;

export function UtmCleanup() {
  useEffect(() => {
    const timer = setTimeout(() => {
      try {
        const url = new URL(window.location.href);
        let changed = false;
        const params = Array.from(url.searchParams.keys());
        for (const key of params) {
          if (UTM_PARAM_RE.test(key)) {
            url.searchParams.delete(key);
            changed = true;
          }
        }
        if (changed) {
          const cleaned =
            url.pathname +
            (url.searchParams.toString() ? '?' + url.searchParams.toString() : '') +
            url.hash;
          window.history.replaceState(window.history.state, '', cleaned);
        }
      } catch {}
    }, 1500);
    return () => clearTimeout(timer);
  }, []);

  return null;
}
