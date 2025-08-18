"use client";

import Link from 'next/link';
import { motion } from 'framer-motion';
import { Button } from '@/shared/ui/Button';

export default function NotFound() {
  return (
    <div className="min-h-screen flex items-center justify-center px-4">
      <div className="text-center">
        <motion.h1
          className="text-9xl font-bold text-blue-600 dark:text-blue-400"
          initial={{ opacity: 0, y: -50 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5 }}
        >
          404
        </motion.h1>
        <motion.h2
          className="text-3xl md:text-4xl font-bold mt-4 mb-6"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ duration: 0.5, delay: 0.2 }}
        >
          Страница не найдена
        </motion.h2>
        <motion.p
          className="text-gray-600 dark:text-gray-300 mb-8 max-w-md mx-auto"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ duration: 0.5, delay: 0.4 }}
        >
          Извините, но страница, которую вы ищете, не существует или была перемещена.
        </motion.p>
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5, delay: 0.6 }}
        >
          <Button href="/">Вернуться на главную</Button>
        </motion.div>
      </div>
    </div>
  );
}
