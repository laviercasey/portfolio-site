'use client';

import { motion, AnimatePresence } from 'framer-motion';
import { Sun, Moon } from 'lucide-react';
import { useTheme, cn } from '@/shared/lib';

const modes = [
  { value: 'dark'  as const, icon: Moon, label: 'Тёмная'  },
  { value: 'light' as const, icon: Sun,  label: 'Светлая' },
];

export default function ThemeToggle() {
  const { theme, setTheme } = useTheme();

  const toggle = () => {
    setTheme(theme === 'dark' ? 'light' : 'dark');
  };

  const current = modes.find((m) => m.value === theme) ?? modes[0];
  const Icon    = current.icon;

  return (
    <button
      onClick={toggle}
      aria-label={`Тема: ${current.label}`}
      title={`Тема: ${current.label}`}
      className={cn(
        'relative overflow-hidden flex items-center justify-center w-9 h-9 rounded-xl',
        'text-muted-foreground hover:text-foreground',
        'bg-transparent hover:bg-white/10',
        'border border-transparent hover:border-white/15',
        'transition-colors duration-200'
      )}
    >
      <AnimatePresence mode="wait">
        <motion.div
          key={current.value}
          initial={{ y: 10, opacity: 0, scale: 0.6 }}
          animate={{ y: 0,  opacity: 1, scale: 1   }}
          exit=  {{ y: -10, opacity: 0, scale: 0.6 }}
          transition={{ duration: 0.18, ease: [0.25, 0.1, 0.25, 1] }}
        >
          <Icon className="h-4 w-4" strokeWidth={1.75} />
        </motion.div>
      </AnimatePresence>
    </button>
  );
}
