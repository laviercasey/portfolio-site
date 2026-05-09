'use client';

import { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { useLocale } from 'next-intl';
import { Plus } from 'lucide-react';
import type { ServiceFaq } from '@/entities/service';

interface ServicesFaqProps {
  faqs: ServiceFaq[];
}

export default function ServicesFaq({ faqs }: ServicesFaqProps) {
  const locale = useLocale();
  const isRu = locale === 'ru';
  const [openIdx, setOpenIdx] = useState<number | null>(0);

  if (faqs.length === 0) return null;

  return (
    <section className="py-24 relative">
      <div className="container max-w-4xl">
        <header className="mb-16">
          <div className="flex items-center gap-3 mb-5 text-xs font-mono uppercase tracking-[0.3em] text-muted-foreground">
            <span className="w-8 h-px bg-primary/60" />
            <span className="text-primary">FAQ</span>
          </div>
          <h2 className="font-heading text-4xl md:text-6xl font-black leading-[1.05] mb-5">
            {isRu ? (
              <>
                Частые <span className="italic gradient-text">вопросы</span>
              </>
            ) : (
              <>
                Common <span className="italic gradient-text">questions</span>
              </>
            )}
          </h2>
        </header>

        <div className="space-y-3">
          {faqs.map((f, i) => {
            const open = openIdx === i;
            return (
              <motion.div
                key={f.id}
                initial={{ opacity: 0, y: 12 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true, margin: '-30px' }}
                transition={{ duration: 0.4, delay: i * 0.04 }}
                className={`glass-card overflow-hidden transition-colors ${
                  open ? 'border-primary/40' : 'hover:border-white/20'
                }`}
              >
                <button
                  onClick={() => setOpenIdx(open ? null : i)}
                  className="w-full flex items-center justify-between gap-4 p-5 md:p-6 text-left"
                  aria-expanded={open}
                >
                  <span className="font-heading font-bold text-lg md:text-xl text-foreground/95 leading-snug">
                    {f.question[isRu ? 'ru' : 'en']}
                  </span>
                  <span
                    className={`shrink-0 w-9 h-9 rounded-full flex items-center justify-center transition-all duration-300 ${
                      open ? 'rotate-45 bg-primary/15 text-primary border-primary/40' : 'bg-white/[0.04] text-muted-foreground border-white/10'
                    } border`}
                  >
                    <Plus className="h-4 w-4" strokeWidth={2} />
                  </span>
                </button>
                <AnimatePresence initial={false}>
                  {open && (
                    <motion.div
                      initial={{ height: 0, opacity: 0 }}
                      animate={{ height: 'auto', opacity: 1 }}
                      exit={{ height: 0, opacity: 0 }}
                      transition={{ duration: 0.35, ease: [0.22, 1, 0.36, 1] }}
                      className="overflow-hidden"
                    >
                      <div className="px-5 md:px-6 pb-6 text-foreground/70 leading-relaxed text-[15px] md:text-base border-t border-white/5 pt-5">
                        {f.answer[isRu ? 'ru' : 'en']}
                      </div>
                    </motion.div>
                  )}
                </AnimatePresence>
              </motion.div>
            );
          })}
        </div>
      </div>
    </section>
  );
}
