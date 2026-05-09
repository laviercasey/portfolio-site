'use client';

import { useEffect, useMemo, useRef, useState } from 'react';
import { useInView, useReducedMotion } from 'framer-motion';
import { useLocale } from 'next-intl';
import type { Highlight } from '@/entities/project';

interface Props {
  highlights: Highlight[];
}

export default function ProjectHighlights({ highlights }: Props) {
  if (!highlights || highlights.length === 0) return null;

  return (
    <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
      {highlights.map((h, i) => (
        <HighlightCard key={`${i}-${h.value}`} highlight={h} delay={i * 0.08} />
      ))}
    </div>
  );
}

function HighlightCard({ highlight, delay }: { highlight: Highlight; delay: number }) {
  const locale = useLocale();
  const isRu = locale === 'ru';
  const ref = useRef<HTMLDivElement>(null);
  const inView = useInView(ref, { once: true, margin: '-40px' });
  const label = isRu ? highlight.label.ru : highlight.label.en;

  return (
    <div
      ref={ref}
      className="glass-card rounded-lg p-4 transition-all"
      style={{
        opacity: inView ? 1 : 0,
        transform: inView ? 'translateY(0)' : 'translateY(12px)',
        transition: `opacity 0.5s ease ${delay}s, transform 0.5s ease ${delay}s`,
      }}
    >
      <div className="font-heading text-2xl md:text-3xl font-black text-primary leading-none">
        <AnimatedValue value={highlight.value} inView={inView} delay={delay} />
      </div>
      <p className="mt-1.5 text-xs md:text-sm text-muted-foreground leading-snug">{label}</p>
    </div>
  );
}

function parseValue(raw: string): { numeric: number | null; suffix: string; decimals: number } {
  const match = raw.match(/^([\d.]+)(.*)$/);
  if (!match) return { numeric: null, suffix: '', decimals: 0 };
  const numeric = parseFloat(match[1]);
  if (!Number.isFinite(numeric)) return { numeric: null, suffix: '', decimals: 0 };
  const decimals = match[1].includes('.') ? 1 : 0;
  return { numeric, suffix: match[2] ?? '', decimals };
}

function AnimatedValue({
  value,
  inView,
  delay,
}: {
  value: string;
  inView: boolean;
  delay: number;
}) {
  const shouldReduce = useReducedMotion();
  const parsed = useMemo(() => parseValue(value), [value]);
  const [displayed, setDisplayed] = useState<string>(() =>
    parsed.numeric === null ? value : `0${parsed.suffix}`,
  );

  useEffect(() => {
    if (!inView) return;

    if (shouldReduce || parsed.numeric === null) {
      setDisplayed(value);
      return;
    }

    const duration = 900;
    const target = parsed.numeric;
    const suffix = parsed.suffix;
    const decimals = parsed.decimals;
    const startAt = performance.now() + delay * 1000;
    let raf = 0;

    const tick = (now: number) => {
      const progress = Math.min(1, Math.max(0, (now - startAt) / duration));
      const eased = 1 - Math.pow(1 - progress, 3);
      const current = target * eased;
      setDisplayed(`${current.toFixed(decimals)}${suffix}`);
      if (progress < 1) {
        raf = requestAnimationFrame(tick);
      }
    };
    raf = requestAnimationFrame(tick);

    return () => {
      if (raf) cancelAnimationFrame(raf);
    };
  }, [inView, shouldReduce, parsed, value, delay]);

  return <span>{displayed}</span>;
}
