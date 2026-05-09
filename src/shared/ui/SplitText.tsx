'use client';

import { motion, useReducedMotion } from 'framer-motion';

interface SplitTextProps {
  text: string;
  className?: string;
  delay?: number;
  as?: 'h1' | 'h2' | 'h3' | 'span' | 'p';
  animate?: boolean;
}

export default function SplitText({
  text,
  className,
  delay = 0,
  as: Tag = 'span',
  animate = true,
}: SplitTextProps) {
  const shouldReduceMotion = useReducedMotion();

  if (shouldReduceMotion || !animate) {
    return <Tag className={className}>{text}</Tag>;
  }

  const words = text.split(' ');
  let charCount = 0;

  return (
    <Tag className={className} aria-label={text}>
      {words.map((word, wi) => (
        <span
          key={wi}
          className="inline-block overflow-hidden"
          style={wi < words.length - 1 ? { marginRight: '0.22em' } : undefined}
        >
          {word.split('').map((char) => {
            const index = charCount++;
            return (
              <motion.span
                key={index}
                className="inline-block"
                initial={{ y: '110%' }}
                animate={{ y: '0%' }}
                transition={{
                  delay: delay + index * 0.028,
                  duration: 0.55,
                  ease: [0.22, 1, 0.36, 1],
                }}
              >
                {char}
              </motion.span>
            );
          })}
        </span>
      ))}
    </Tag>
  );
}
