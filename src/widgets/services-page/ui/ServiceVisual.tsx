'use client';

import { motion } from 'framer-motion';
import type { ServiceVisualKey } from '@/entities/service';

interface ServiceVisualProps {
  variant: ServiceVisualKey;
  accent: string;
}

export default function ServiceVisual({ variant, accent }: ServiceVisualProps) {
  if (variant === 'terminal') return <TerminalMock accent={accent} />;
  if (variant === 'browser') return <BrowserMock accent={accent} />;
  if (variant === 'editor') return <EditorMock accent={accent} />;
  if (variant === 'dataframe') return <DataframeMock accent={accent} />;
  if (variant === 'pipeline') return <PipelineMock accent={accent} />;
  return null;
}

function WindowChrome({ title, accent }: { title: string; accent: string }) {
  return (
    <div className="flex items-center gap-2 px-4 py-3 border-b border-white/8 bg-white/[0.02]">
      <span className="w-2.5 h-2.5 rounded-full bg-red-500/70" />
      <span className="w-2.5 h-2.5 rounded-full bg-yellow-500/70" />
      <span className="w-2.5 h-2.5 rounded-full bg-green-500/70" />
      <span className="ml-3 text-xs font-mono text-muted-foreground/70" style={{ color: `${accent}cc` }}>
        {title}
      </span>
    </div>
  );
}

function Wrapper({ children }: { children: React.ReactNode }) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 24 }}
      whileInView={{ opacity: 1, y: 0 }}
      viewport={{ once: true, margin: '-80px' }}
      transition={{ duration: 0.7, ease: [0.22, 1, 0.36, 1] }}
      className="glass-card overflow-hidden font-mono text-xs md:text-[13px] leading-relaxed"
      style={{ boxShadow: '0 30px 80px -30px rgba(0,0,0,0.5)' }}
    >
      {children}
    </motion.div>
  );
}

function TerminalMock({ accent }: { accent: string }) {
  const lines = [
    { p: '$', t: 'aiogram run bot.py', c: 'text-foreground' },
    { p: '', t: 'Bot @MedReminderBot started ✓', c: 'text-green-400/85' },
    { p: '', t: '', c: '' },
    { p: '>', t: 'Active users: 12,453', c: 'text-foreground/75' },
    { p: '>', t: 'Commands handled today: 3,891', c: 'text-foreground/75' },
    { p: '>', t: 'Avg response: 47ms', c: 'text-foreground/75' },
    { p: '>', t: 'Uptime: 99.8% (30d)', c: 'text-foreground/75' },
    { p: '', t: '', c: '' },
    { p: '$', t: 'tail -f logs/payments.log', c: 'text-foreground' },
    { p: '', t: '[14:23:51] paid: user_8214 → 590₽', c: 'text-cyan-400/80' },
    { p: '', t: '[14:24:18] paid: user_3015 → 1290₽', c: 'text-cyan-400/80' },
    { p: '', t: '[14:24:32] paid: user_9482 → 290₽', c: 'text-cyan-400/80' },
    { p: '_', t: '', c: 'text-primary animate-pulse' },
  ];
  return (
    <Wrapper>
      <WindowChrome title="bot.py — terminal" accent={accent} />
      <div className="p-5 bg-[#0a0a0a]/40 min-h-[320px]">
        {lines.map((l, i) => (
          <div key={i} className="flex gap-3 min-h-[1.6em]">
            <span className="text-muted-foreground/40 select-none w-3">{l.p}</span>
            <span className={l.c}>{l.t}</span>
          </div>
        ))}
      </div>
    </Wrapper>
  );
}

function BrowserMock({ accent }: { accent: string }) {
  return (
    <Wrapper>
      <div className="flex items-center gap-2 px-4 py-3 border-b border-white/8 bg-white/[0.02]">
        <span className="w-2.5 h-2.5 rounded-full bg-red-500/70" />
        <span className="w-2.5 h-2.5 rounded-full bg-yellow-500/70" />
        <span className="w-2.5 h-2.5 rounded-full bg-green-500/70" />
        <div className="ml-3 flex-1 px-3 py-1 rounded-md bg-white/[0.04] text-xs text-muted-foreground/70 truncate">
          🔒 fluxrouter.app/admin/orders
        </div>
      </div>
      <div className="bg-[#0a0a0a]/40 min-h-[340px]">
        <div className="flex">
          <div className="w-32 border-r border-white/5 py-4 px-3 space-y-2.5">
            {['Dashboard', 'Orders', 'Users', 'Products', 'Reports', 'Settings'].map((item, i) => (
              <div
                key={item}
                className={`text-[11px] px-2 py-1.5 rounded-md ${
                  i === 1 ? '' : 'text-muted-foreground/60'
                }`}
                style={i === 1 ? { background: `${accent}22`, color: accent } : {}}
              >
                {item}
              </div>
            ))}
          </div>
          <div className="flex-1 p-5">
            <div className="grid grid-cols-3 gap-3 mb-5">
              {[
                { l: 'Revenue', v: '₽ 2.4M', c: '#22c55e' },
                { l: 'Orders', v: '1,847', c: accent },
                { l: 'AOV', v: '₽ 1,302', c: '#a855f7' },
              ].map((s) => (
                <div key={s.l} className="rounded-lg border border-white/8 bg-white/[0.02] p-3">
                  <div className="text-[10px] uppercase tracking-wider text-muted-foreground/60 mb-1">{s.l}</div>
                  <div className="text-base font-bold" style={{ color: s.c }}>{s.v}</div>
                </div>
              ))}
            </div>
            <div className="space-y-2">
              {[
                ['#4827', 'Анна К.', '590 ₽', 'paid'],
                ['#4826', 'Игорь М.', '1290 ₽', 'paid'],
                ['#4825', 'Лена П.', '290 ₽', 'pending'],
                ['#4824', 'Сергей В.', '4500 ₽', 'paid'],
              ].map((row) => (
                <div key={row[0]} className="grid grid-cols-[60px_1fr_70px_60px] gap-2 py-2 border-b border-white/5 text-[11px]">
                  <span className="text-muted-foreground/70">{row[0]}</span>
                  <span className="text-foreground/85">{row[1]}</span>
                  <span className="text-foreground/85 text-right">{row[2]}</span>
                  <span
                    className="text-center text-[10px] px-1.5 py-0.5 rounded uppercase tracking-wider"
                    style={{
                      background: row[3] === 'paid' ? '#22c55e22' : '#eab30822',
                      color: row[3] === 'paid' ? '#22c55e' : '#eab308',
                    }}
                  >
                    {row[3]}
                  </span>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </Wrapper>
  );
}

function EditorMock({ accent }: { accent: string }) {
  const code: { n: number; s: { c: string; t: string }[] }[] = [
    { n: 1, s: [{ c: 'text-pink-400/85', t: 'export async function' }, { c: 'text-foreground/85', t: ' ' }, { c: 'text-cyan-400/85', t: 'generateMetadata' }, { c: 'text-foreground/65', t: '() {' }] },
    { n: 2, s: [{ c: 'text-foreground/65', t: '  ' }, { c: 'text-pink-400/85', t: 'return' }, { c: 'text-foreground/65', t: ' {' }] },
    { n: 3, s: [{ c: 'text-foreground/65', t: '    title: ' }, { c: 'text-amber-300/85', t: "'Casey Laviere'" }, { c: 'text-foreground/65', t: ',' }] },
    { n: 4, s: [{ c: 'text-foreground/65', t: '    description: ' }, { c: 'text-amber-300/85', t: "'Full-Stack…'" }] },
    { n: 5, s: [{ c: 'text-foreground/65', t: '  };' }] },
    { n: 6, s: [{ c: 'text-foreground/65', t: '}' }] },
    { n: 7, s: [{ c: '', t: '' }] },
    { n: 8, s: [{ c: 'text-pink-400/85', t: 'export default async function' }, { c: 'text-foreground/85', t: ' ' }, { c: 'text-cyan-400/85', t: 'Page' }, { c: 'text-foreground/65', t: '() {' }] },
    { n: 9, s: [{ c: 'text-foreground/65', t: '  ' }, { c: 'text-pink-400/85', t: 'const' }, { c: 'text-foreground/85', t: ' data ' }, { c: 'text-foreground/65', t: '= ' }, { c: 'text-pink-400/85', t: 'await' }, { c: 'text-foreground/85', t: ' ' }, { c: 'text-cyan-400/85', t: 'fetchData' }, { c: 'text-foreground/65', t: '();' }] },
    { n: 10, s: [{ c: 'text-foreground/65', t: '  ' }, { c: 'text-pink-400/85', t: 'return' }, { c: 'text-foreground/65', t: ' <' }, { c: 'text-cyan-400/85', t: 'Hero' }, { c: 'text-foreground/65', t: ' data={data} />;' }] },
    { n: 11, s: [{ c: 'text-foreground/65', t: '}' }] },
  ];
  return (
    <Wrapper>
      <WindowChrome title="app/page.tsx — Next.js" accent={accent} />
      <div className="bg-[#0a0a0a]/40 min-h-[340px] flex">
        <div className="py-4 px-3 border-r border-white/5 bg-white/[0.01]">
          {code.map((row) => (
            <div key={row.n} className="text-muted-foreground/40 text-right select-none min-h-[1.6em]">
              {row.n}
            </div>
          ))}
        </div>
        <div className="flex-1 py-4 px-4">
          {code.map((row) => (
            <div key={row.n} className="min-h-[1.6em] whitespace-pre">
              {row.s.map((seg, i) => (
                <span key={i} className={seg.c}>{seg.t}</span>
              ))}
            </div>
          ))}
        </div>
      </div>
    </Wrapper>
  );
}

function DataframeMock({ accent }: { accent: string }) {
  return (
    <Wrapper>
      <WindowChrome title="heart_disease_benchmark.ipynb" accent={accent} />
      <div className="bg-[#0a0a0a]/40 p-5 min-h-[340px]">
        <div className="text-muted-foreground/60 mb-2">In [4]: df.describe()</div>
        <div className="overflow-x-auto -mx-2 px-2 mb-5">
          <table className="text-[11px] w-full">
            <thead>
              <tr className="text-muted-foreground/65 border-b border-white/8">
                {['', 'age', 'chol', 'thalach', 'oldpeak', 'target'].map((h) => (
                  <th key={h} className="text-right py-1.5 px-2 font-normal">{h}</th>
                ))}
              </tr>
            </thead>
            <tbody>
              {[
                ['mean', '54.4', '246.7', '149.6', '1.04', '0.54'],
                ['std', '9.1', '51.7', '22.9', '1.16', '0.50'],
                ['min', '29.0', '126.0', '71.0', '0.00', '0.00'],
                ['max', '77.0', '564.0', '202.0', '6.20', '1.00'],
              ].map((row) => (
                <tr key={row[0]} className="border-b border-white/5">
                  <td className="py-1.5 px-2 text-right" style={{ color: accent }}>{row[0]}</td>
                  {row.slice(1).map((v, i) => (
                    <td key={i} className="py-1.5 px-2 text-right text-foreground/80">{v}</td>
                  ))}
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        <div className="text-muted-foreground/60 mb-2">In [12]: model.score(X_test, y_test)</div>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-2.5 mb-3">
          {[
            { l: 'Accuracy', v: '0.882', c: '#22c55e' },
            { l: 'F1', v: '0.871', c: accent },
            { l: 'AUC', v: '0.912', c: '#22c55e' },
            { l: 'Recall', v: '0.864', c: accent },
          ].map((m) => (
            <div key={m.l} className="rounded-lg border border-white/8 bg-white/[0.02] p-3">
              <div className="text-[10px] uppercase tracking-wider text-muted-foreground/60">{m.l}</div>
              <div className="text-base font-bold mt-1" style={{ color: m.c }}>{m.v}</div>
            </div>
          ))}
        </div>
        <div className="text-foreground/55 text-[11px]">
          → XGBoost · cross-val 5-fold · features: 13
        </div>
      </div>
    </Wrapper>
  );
}

function PipelineMock({ accent }: { accent: string }) {
  const Node = ({ label, sub, color, x }: { label: string; sub: string; color: string; x: string }) => (
    <div className="absolute" style={{ left: x, transform: 'translateX(-50%)' }}>
      <div
        className="rounded-xl border px-3 py-2 backdrop-blur min-w-[100px] text-center"
        style={{
          background: `${color}1a`,
          borderColor: `${color}55`,
        }}
      >
        <div className="text-[11px] font-bold" style={{ color }}>{label}</div>
        <div className="text-[9px] text-muted-foreground/70 mt-0.5">{sub}</div>
      </div>
    </div>
  );

  return (
    <Wrapper>
      <WindowChrome title="parser.py → pipeline" accent={accent} />
      <div className="bg-[#0a0a0a]/40 min-h-[340px] relative px-4 py-12">
        <div className="relative h-32">
          <svg className="absolute inset-0 w-full h-full" preserveAspectRatio="none" viewBox="0 0 100 30">
            <line x1="10" y1="15" x2="90" y2="15" stroke={accent} strokeWidth="0.4" strokeDasharray="1.5 1" opacity="0.5" />
            <circle cx="35" cy="15" r="0.6" fill={accent}>
              <animate attributeName="cx" values="10;90;10" dur="4s" repeatCount="indefinite" />
            </circle>
          </svg>
          <Node label="SCRAPE" sub="playwright" color="#22c55e" x="15%" />
          <Node label="PARSE" sub="cheerio" color={accent} x="40%" />
          <Node label="VALIDATE" sub="zod" color="#a855f7" x="65%" />
          <Node label="STORE" sub="postgres" color="#5eb3ff" x="90%" />
        </div>

        <div className="mt-10 text-[11px] space-y-1.5 text-foreground/70">
          <div className="flex justify-between">
            <span className="text-muted-foreground/60">→ wildberries.ru/catalog</span>
            <span className="text-green-400/80">2,847 items</span>
          </div>
          <div className="flex justify-between">
            <span className="text-muted-foreground/60">→ ozon.ru/category</span>
            <span className="text-green-400/80">1,392 items</span>
          </div>
          <div className="flex justify-between">
            <span className="text-muted-foreground/60">→ avito.ru/feed</span>
            <span className="text-green-400/80">5,108 items</span>
          </div>
          <div className="flex justify-between border-t border-white/8 pt-2 mt-2">
            <span className="text-foreground/85 font-bold">Total synced today</span>
            <span style={{ color: accent }} className="font-bold">9,347</span>
          </div>
        </div>
      </div>
    </Wrapper>
  );
}
