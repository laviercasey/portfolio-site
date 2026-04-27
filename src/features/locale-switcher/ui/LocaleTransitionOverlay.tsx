'use client';

import { useEffect, useMemo, useState } from 'react';
import { createPortal } from 'react-dom';
import { motion } from 'framer-motion';

const WORD = 'TRANSLATION';

interface Props {
  targetLocale: string;
  targetPathname: string;
}

function Overlay({ targetLocale, targetPathname }: Props) {
  const [phase, setPhase] = useState<'enter' | 'hold' | 'exit'>('enter');

  const letters = useMemo(
    () =>
      WORD.split('').map(() => ({
        fromX:     -(320 + Math.random() * 280),
        fromY:     (Math.random() - 0.5) * 180,
        fromR:     (Math.random() - 0.5) * 70,
        toX:        320 + Math.random() * 280,
        toY:       (Math.random() - 0.5) * 180,
        toR:       (Math.random() - 0.5) * 70,
      })),
    [],
  );

  useEffect(() => {
    const t1 = setTimeout(() => setPhase('hold'), 650);
    const t2 = setTimeout(() => setPhase('exit'), 900);
    const t3 = setTimeout(() => {
      const href = `/${targetLocale}${targetPathname === '/' ? '' : targetPathname}`;
      window.location.href = href;
    }, 1590);

    return () => {
      clearTimeout(t1);
      clearTimeout(t2);
      clearTimeout(t3);
    };
  }, [targetPathname, targetLocale]);

  return (
    <motion.div
      className="fixed inset-0 z-[9999] flex flex-col items-center justify-center select-none"
      style={{ background: 'rgba(8, 6, 4, 0.92)', backdropFilter: 'blur(32px)' }}
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      transition={{ duration: 0.12 }}
    >
      <div className="flex items-end gap-[2px] md:gap-1">
        {WORD.split('').map((char, i) => (
          <motion.span
            key={i}
            className="font-heading font-black leading-none"
            style={{
              fontSize: 'clamp(2.5rem, 9vw, 7rem)',
              background: 'linear-gradient(135deg, #d4a843 0%, #e07040 55%, #c45a30 100%)',
              WebkitBackgroundClip: 'text',
              WebkitTextFillColor: 'transparent',
              backgroundClip: 'text',
            }}
            initial={{
              x: letters[i].fromX,
              y: letters[i].fromY,
              rotate: letters[i].fromR,
              opacity: 0,
            }}
            animate={
              phase === 'enter' || phase === 'hold'
                ? { x: 0, y: 0, rotate: 0, opacity: 1 }
                : {
                    x: letters[i].toX,
                    y: letters[i].toY,
                    rotate: letters[i].toR,
                    opacity: 0,
                  }
            }
            transition={
              phase === 'enter'
                ? {
                    delay: i * 0.042,
                    type: 'spring',
                    stiffness: 110,
                    damping: 13,
                  }
                : {
                    delay: i * 0.028,
                    duration: 0.38,
                    ease: [0.4, 0, 1, 1],
                  }
            }
          >
            {char}
          </motion.span>
        ))}
      </div>

      <motion.p
        className="mt-5 font-mono text-xs tracking-[0.35em] uppercase text-white/35"
        initial={{ opacity: 0, y: 8 }}
        animate={{ opacity: phase === 'hold' ? 1 : 0, y: phase === 'hold' ? 0 : 8 }}
        transition={{ duration: 0.35 }}
      >
        switching to {targetLocale}
      </motion.p>

      <motion.div
        className="absolute bottom-0 left-0 h-0.5 bg-gradient-to-r from-transparent via-primary to-transparent"
        initial={{ width: '0%', left: '50%' }}
        animate={{ width: '100%', left: '0%' }}
        transition={{ duration: 0.9, ease: 'easeInOut' }}
      />
    </motion.div>
  );
}

export default function LocaleTransitionOverlay(props: Props) {
  if (typeof document === 'undefined') return null;
  return createPortal(<Overlay {...props} />, document.body);
}
