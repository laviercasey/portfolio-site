import { ReactNode } from 'react';

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

export default function ProjectStoryBlock({ label, title, accent = 'approach', children }: Props) {
  return (
    <section className="relative">
      <div className={`font-mono text-xs md:text-sm tracking-[0.3em] uppercase mb-2 md:mb-3 ${ACCENTS[accent]}`}>
        {label}
      </div>
      <h3 className="font-heading text-2xl md:text-3xl font-black mb-4 leading-tight">{title}</h3>
      <div className="text-sm md:text-base lg:text-lg text-muted-foreground leading-relaxed whitespace-pre-line">
        {children}
      </div>
    </section>
  );
}
