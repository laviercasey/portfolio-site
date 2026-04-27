'use client';

import { motion, useReducedMotion } from 'framer-motion';

interface ClipRevealProps {
  children: React.ReactNode;
  direction?: 'up' | 'down' | 'left' | 'right';
  delay?: number;
  duration?: number;
  className?: string;
  once?: boolean;
}

const clips = {
  up:    { hidden: 'inset(100% 0% 0% 0%)',  visible: 'inset(0% 0% 0% 0%)' },
  down:  { hidden: 'inset(0% 0% 100% 0%)',  visible: 'inset(0% 0% 0% 0%)' },
  left:  { hidden: 'inset(0% 100% 0% 0%)',  visible: 'inset(0% 0% 0% 0%)' },
  right: { hidden: 'inset(0% 0% 0% 100%)',  visible: 'inset(0% 0% 0% 0%)' },
};

export default function ClipReveal({
  children,
  direction = 'up',
  delay = 0,
  duration = 0.8,
  className,
  once = true,
}: ClipRevealProps) {
  const shouldReduceMotion = useReducedMotion();
  const { hidden, visible } = clips[direction];

  if (shouldReduceMotion) {
    return <div className={className}>{children}</div>;
  }

  return (
    <motion.div
      className={className}
      initial={{ clipPath: hidden, opacity: 0 }}
      whileInView={{ clipPath: visible, opacity: 1 }}
      viewport={{ once, margin: '-60px' }}
      transition={{ duration, delay, ease: [0.16, 1, 0.3, 1] }}
    >
      {children}
    </motion.div>
  );
}
