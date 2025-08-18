"use client";

import { motion, TargetAndTransition, VariantLabels, useReducedMotion } from 'framer-motion';
import { ReactNode } from 'react';

type CardProps = {
  children: ReactNode;
  className?: string;
  onClick?: () => void;
  whileHover?: TargetAndTransition | VariantLabels;
};

export const Card = ({
  children,
  className = '',
  onClick,
  whileHover = { scale: 1.05, y: -10 },
}: CardProps) => {
  const shouldReduceMotion = useReducedMotion();
  
  return (
    <motion.div
      className={`bg-white dark:bg-gray-800 rounded-xl shadow-lg overflow-hidden ${className}`}
      whileHover={shouldReduceMotion ? {} : whileHover}
      whileTap={shouldReduceMotion ? {} : { scale: 0.98 }}
      onClick={onClick}
      initial={{ opacity: 0, y: shouldReduceMotion ? 0 : 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: shouldReduceMotion ? 0.1 : 0.5 }}
      style={{ willChange: shouldReduceMotion ? 'auto' : 'transform' }}
    >
      {children}
    </motion.div>
  );
};
