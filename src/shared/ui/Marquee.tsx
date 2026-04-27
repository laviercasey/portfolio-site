'use client';

import { useReducedMotion } from 'framer-motion';

interface MarqueeProps {
  items: string[];
  separator?: string;
  speed?: number;
  className?: string;
  itemClassName?: string;
  reverse?: boolean;
}

export default function Marquee({
  items,
  separator = '·',
  speed = 60,
  className,
  itemClassName,
  reverse = false,
}: MarqueeProps) {
  const shouldReduceMotion = useReducedMotion();
  const repeated = [...items, ...items, ...items];

  const estimatedWidth = items.length * 200;
  const duration = estimatedWidth / speed;

  return (
    <div
      className={`overflow-hidden whitespace-nowrap select-none ${className}`}
      aria-hidden="true"
    >
      <div
        className="inline-flex"
        style={
          shouldReduceMotion
            ? {}
            : {
                animation: `marquee-scroll ${duration}s linear infinite`,
                animationDirection: reverse ? 'reverse' : 'normal',
              }
        }
      >
        {repeated.map((item, i) => (
          <span
            key={i}
            className={`inline-flex items-center gap-5 px-3 ${itemClassName}`}
          >
            <span>{item}</span>
            <span className="opacity-30 text-primary">{separator}</span>
          </span>
        ))}
      </div>
    </div>
  );
}
