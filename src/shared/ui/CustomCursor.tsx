'use client';

import { useEffect, useRef, useState } from 'react';
import { motion, useSpring, useMotionValue } from 'framer-motion';

export default function CustomCursor() {
  const cursorX = useMotionValue(-200);
  const cursorY = useMotionValue(-200);

  const springConfig = { stiffness: 600, damping: 32, mass: 0.4 };
  const ringX = useSpring(cursorX, springConfig);
  const ringY = useSpring(cursorY, springConfig);

  const [isPointer, setIsPointer] = useState(false);
  const [isHidden, setIsHidden] = useState(false);
  const rafRef = useRef<number>(0);

  useEffect(() => {
    if (typeof window === 'undefined') return;
    if (window.matchMedia('(pointer: coarse)').matches) return;

    document.body.classList.add('custom-cursor-active');

    const onMove = (e: MouseEvent) => {
      cursorX.set(e.clientX);
      cursorY.set(e.clientY);

      cancelAnimationFrame(rafRef.current);
      rafRef.current = requestAnimationFrame(() => {
        const el = document.elementFromPoint(e.clientX, e.clientY);
        if (el) {
          const style = window.getComputedStyle(el as Element);
          setIsPointer(style.cursor === 'pointer');
        }
      });
    };

    const onLeave = () => setIsHidden(true);
    const onEnter = () => setIsHidden(false);

    window.addEventListener('mousemove', onMove, { passive: true });
    document.addEventListener('mouseleave', onLeave);
    document.addEventListener('mouseenter', onEnter);

    return () => {
      document.body.classList.remove('custom-cursor-active');
      window.removeEventListener('mousemove', onMove);
      document.removeEventListener('mouseleave', onLeave);
      document.removeEventListener('mouseenter', onEnter);
      cancelAnimationFrame(rafRef.current);
    };
  }, [cursorX, cursorY]);

  return (
    <>
      <motion.div
        className="fixed top-0 left-0 z-[9999] pointer-events-none rounded-full"
        style={{
          x: cursorX,
          y: cursorY,
          translateX: '-50%',
          translateY: '-50%',
          width: 6,
          height: 6,
        }}
        animate={{
          scale: isPointer ? 0 : 1,
          opacity: isHidden ? 0 : 1,
          backgroundColor: '#d4a843',
        }}
        transition={{ duration: 0.15 }}
      />

      <motion.div
        className="fixed top-0 left-0 z-[9999] pointer-events-none rounded-full border"
        style={{
          x: ringX,
          y: ringY,
          translateX: '-50%',
          translateY: '-50%',
          width: 32,
          height: 32,
          borderColor: 'rgba(212,168,67,0.5)',
          mixBlendMode: 'difference',
        }}
        animate={{
          scale: isPointer ? 1.8 : 1,
          opacity: isHidden ? 0 : isPointer ? 0.9 : 0.5,
        }}
        transition={{ duration: 0.2 }}
      />
    </>
  );
}
