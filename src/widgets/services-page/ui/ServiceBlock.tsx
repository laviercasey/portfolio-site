'use client';

import { motion } from 'framer-motion';
import { useLocale } from 'next-intl';
import { ArrowRight } from 'lucide-react';
import { Link } from '@/shared/config';
import { Button } from '@/shared/ui';
import ServiceVisual from './ServiceVisual';
import { SERVICE_ICON_MAP } from './service-icons';
import type { Service, ServiceIconKey, ServiceVisualKey } from '@/entities/service';

interface ServiceBlockProps {
  service: Service;
  index: number;
}

export default function ServiceBlock({ service, index }: ServiceBlockProps) {
  const locale = useLocale();
  const isRu = locale === 'ru';
  const flip = index % 2 === 1;

  const reveal = {
    hidden: { opacity: 0, y: 28 },
    visible: { opacity: 1, y: 0, transition: { duration: 0.65, ease: [0.22, 1, 0.36, 1] as const } },
  };

  const Icon = SERVICE_ICON_MAP[service.iconKey as ServiceIconKey] ?? SERVICE_ICON_MAP.bot;

  return (
    <section id={service.slug} className="scroll-mt-28 relative py-20 md:py-28">
      <div
        className="absolute -z-10 rounded-full pointer-events-none hidden md:block"
        style={{
          width: 600, height: 600,
          [flip ? 'left' : 'right']: '-15%',
          top: '20%',
          background: `radial-gradient(circle, ${service.accent}11 0%, transparent 60%)`,
          filter: 'blur(60px)',
        }}
      />

      <motion.div
        variants={reveal}
        initial="hidden"
        whileInView="visible"
        viewport={{ once: true, margin: '-100px' }}
      >
        <div
          className={`grid grid-cols-1 lg:grid-cols-2 gap-10 lg:gap-14 items-start mb-10 ${
            flip ? 'lg:[&>*:first-child]:order-2' : ''
          }`}
        >
          <div>
            <div className="flex items-center gap-4 mb-6">
              <span
                className="font-mono text-xs font-bold px-2.5 py-1 rounded-md"
                style={{ background: `${service.accent}22`, color: service.accent, border: `1px solid ${service.accent}44` }}
              >
                {service.num}
              </span>
              <span className="h-px flex-1 max-w-[80px]" style={{ background: `${service.accent}44` }} />
              <Icon className="h-4 w-4 shrink-0" style={{ color: service.accent }} strokeWidth={2} />
            </div>

            <h2 className="font-heading text-4xl md:text-5xl font-black mb-7 leading-[1.05] tracking-tight">
              {service.title[isRu ? 'ru' : 'en']}
            </h2>

            <p className="text-base md:text-lg text-foreground/75 leading-relaxed mb-7">
              {service.lead[isRu ? 'ru' : 'en']}
            </p>

            <ul className="space-y-2.5">
              {(service.bullets[isRu ? 'ru' : 'en'] ?? []).map((b, i) => (
                <li key={i} className="flex items-start gap-3 text-foreground/75">
                  <span className="mt-2.5 w-1 h-1 rounded-full shrink-0" style={{ background: service.accent }} />
                  <span>{b}</span>
                </li>
              ))}
            </ul>
          </div>

          <div className="space-y-4">
            <ServiceVisual variant={service.visualKey as ServiceVisualKey} accent={service.accent} />

            <div className="grid grid-cols-1 sm:grid-cols-3 gap-2.5">
              <Stat
                label={isRu ? 'Стек' : 'Stack'}
                value={service.stack.split(' · ').slice(0, 3).join(' · ')}
              />
              <Stat
                label={isRu ? 'Сроки' : 'Timeline'}
                value={service.timeline[isRu ? 'ru' : 'en']}
              />
              <Stat
                label={isRu ? 'Цена' : 'Price'}
                value={isRu ? `от ${service.priceRu}` : `from ${service.priceEn}`}
                accent={service.accent}
                bold
              />
            </div>
          </div>
        </div>

        <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-6 pt-6 border-t border-white/5">
          {service.caseProjects.length > 0 ? (
            <div className="flex items-center gap-3 flex-wrap">
              <span className="text-[11px] font-mono uppercase tracking-[0.25em] text-muted-foreground/70 shrink-0">
                {isRu ? 'Кейсы' : 'Cases'}
              </span>
              <div className="flex flex-wrap gap-2">
                {service.caseProjects.map((p) => (
                  <Link
                    key={p.slug}
                    href={`/projects/${p.slug}`}
                    className="group inline-flex items-center gap-1.5 text-sm px-3 py-1.5 rounded-lg border border-white/10 bg-white/[0.03] hover:bg-white/[0.06] transition-all hover:translate-x-0.5"
                    style={{ transitionProperty: 'transform, background-color, border-color' }}
                  >
                    {p.name}
                    <ArrowRight className="h-3 w-3 opacity-50 group-hover:opacity-100 transition-opacity" />
                  </Link>
                ))}
              </div>
            </div>
          ) : (
            <div />
          )}

          <Button asChild size="lg" className="gap-2 group shrink-0">
            <Link href="/contact">
              {isRu ? 'Обсудить' : 'Discuss'}
              <ArrowRight className="h-4 w-4 group-hover:translate-x-0.5 transition-transform" />
            </Link>
          </Button>
        </div>
      </motion.div>
    </section>
  );
}

function Stat({ label, value, accent, bold }: { label: string; value: string; accent?: string; bold?: boolean }) {
  return (
    <div className="glass-card rounded-xl p-3.5">
      <div className="text-[10px] uppercase tracking-[0.2em] text-muted-foreground/60 font-mono mb-1">
        {label}
      </div>
      <div
        className={`text-sm leading-snug ${bold ? 'font-bold' : 'text-foreground/85'}`}
        style={accent ? { color: accent } : undefined}
      >
        {value}
      </div>
    </div>
  );
}
