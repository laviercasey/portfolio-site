'use client';

import { ReactNode } from 'react';
import { motion, useReducedMotion } from 'framer-motion';

interface Props {
  label: string;
  title: string;
  accent?: 'problem' | 'approach' | 'outcome';
  children: ReactNode;
  index?: number;
}

const ACCENTS = {
  problem: 'text-red-300/90',
  approach: 'text-primary',
  outcome: 'text-emerald-300',
};

export default function ProjectStoryBlock({ label, title, accent = 'approach', children, index = 0 }: Props) {
  const shouldReduce = useReducedMotion();

  return (
    <motion.section
      initial={shouldReduce ? {} : { clipPath: 'inset(0 100% 0 0)', opacity: 0 }}
      whileInView={{ clipPath: 'inset(0 0% 0 0)', opacity: 1 }}
      viewport={{ once: true, margin: '-80px' }}
      transition={{ duration: 0.7, delay: index * 0.08, ease: [0.16, 1, 0.3, 1] }}
      className="relative"
    >
      <div className={`font-mono text-[10px] tracking-[0.3em] uppercase mb-2 ${ACCENTS[accent]}`}>
        {label}
      </div>
      <h3 className="font-heading text-2xl md:text-3xl font-black mb-4 leading-tight">{title}</h3>
      <div className="text-muted-foreground leading-relaxed whitespace-pre-line">
        {children}
      </div>
    </motion.section>
  );
}
