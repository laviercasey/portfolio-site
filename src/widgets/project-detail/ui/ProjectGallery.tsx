'use client';

import { useCallback, useEffect, useState } from 'react';
import Image from 'next/image';
import * as Dialog from '@radix-ui/react-dialog';
import { ChevronLeft, ChevronRight, X } from 'lucide-react';
import { motion, useReducedMotion } from 'framer-motion';

interface Props {
  images: string[];
  title: string;
}

export default function ProjectGallery({ images, title }: Props) {
  const [openIndex, setOpenIndex] = useState<number | null>(null);
  const shouldReduce = useReducedMotion();

  const close = useCallback(() => setOpenIndex(null), []);
  const prev = useCallback(
    () => setOpenIndex((i) => (i === null ? null : (i - 1 + images.length) % images.length)),
    [images.length],
  );
  const next = useCallback(
    () => setOpenIndex((i) => (i === null ? null : (i + 1) % images.length)),
    [images.length],
  );

  useEffect(() => {
    if (openIndex === null) return;
    const handler = (e: KeyboardEvent) => {
      if (e.key === 'ArrowLeft') prev();
      if (e.key === 'ArrowRight') next();
    };
    window.addEventListener('keydown', handler);
    return () => window.removeEventListener('keydown', handler);
  }, [openIndex, prev, next]);

  if (!images || images.length === 0) return null;

  return (
    <>
      <div className="grid grid-cols-2 md:grid-cols-3 gap-3">
        {images.map((src, i) => (
          <motion.button
            key={`${i}-${src}`}
            type="button"
            onClick={() => setOpenIndex(i)}
            className="relative aspect-[4/3] overflow-hidden rounded-lg group"
            initial={shouldReduce ? {} : { opacity: 0, y: 12 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true, margin: '-40px' }}
            transition={{ duration: 0.35, delay: i * 0.05 }}
          >
            <Image
              src={src}
              alt={`${title} — screenshot ${i + 1}`}
              fill
              className="object-cover transition-all duration-300 group-hover:scale-105 group-hover:brightness-110"
              sizes="(max-width: 768px) 50vw, 33vw"
            />
            <div className="absolute inset-0 bg-gradient-to-t from-black/30 to-transparent opacity-0 group-hover:opacity-100 transition-opacity" />
          </motion.button>
        ))}
      </div>

      <Dialog.Root open={openIndex !== null} onOpenChange={(o) => !o && close()}>
        <Dialog.Portal>
          <Dialog.Overlay className="fixed inset-0 z-50 bg-black/90 backdrop-blur-sm data-[state=open]:animate-in data-[state=open]:fade-in" />
          <Dialog.Content
            className="fixed inset-0 z-50 flex items-center justify-center p-4 md:p-8 outline-none"
            onOpenAutoFocus={(e) => e.preventDefault()}
          >
            <Dialog.Title className="sr-only">{title}</Dialog.Title>
            <Dialog.Description className="sr-only">
              Image viewer. Use arrow keys to navigate, Escape to close.
            </Dialog.Description>

            {openIndex !== null && (
              <div className="relative w-full h-full flex items-center justify-center">
                <Image
                  src={images[openIndex]}
                  alt={`${title} — screenshot ${openIndex + 1}`}
                  fill
                  className="object-contain"
                  sizes="100vw"
                  priority
                />
              </div>
            )}

            <Dialog.Close asChild>
              <button
                type="button"
                aria-label="Close"
                className="absolute top-4 right-4 p-2 rounded-full bg-white/5 hover:bg-white/10 text-white transition-colors"
              >
                <X className="h-5 w-5" />
              </button>
            </Dialog.Close>

            {images.length > 1 && (
              <>
                <button
                  type="button"
                  onClick={prev}
                  aria-label="Previous image"
                  className="absolute left-4 top-1/2 -translate-y-1/2 p-3 rounded-full bg-white/5 hover:bg-white/10 text-white transition-colors"
                >
                  <ChevronLeft className="h-5 w-5" />
                </button>
                <button
                  type="button"
                  onClick={next}
                  aria-label="Next image"
                  className="absolute right-4 top-1/2 -translate-y-1/2 p-3 rounded-full bg-white/5 hover:bg-white/10 text-white transition-colors"
                >
                  <ChevronRight className="h-5 w-5" />
                </button>

                <div className="absolute bottom-6 left-1/2 -translate-x-1/2 px-3 py-1 rounded-full bg-white/5 text-white text-xs font-mono">
                  {(openIndex ?? 0) + 1} / {images.length}
                </div>
              </>
            )}
          </Dialog.Content>
        </Dialog.Portal>
      </Dialog.Root>
    </>
  );
}
