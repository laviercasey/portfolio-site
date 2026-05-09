'use client';

import { useState } from 'react';
import { useLocale } from 'next-intl';
import { Terminal, Copy, Check, Shield, User, Briefcase, Settings } from 'lucide-react';
import type { DemoCredential } from '@/entities/project';

interface Props {
  credentials: DemoCredential[];
}

const ROLE_ICONS = {
  admin: Shield,
  manager: Briefcase,
  user: User,
  custom: Settings,
} as const;

const ROLE_LABELS = {
  admin: { en: 'Admin', ru: 'Админ' },
  manager: { en: 'Manager', ru: 'Менеджер' },
  user: { en: 'User', ru: 'Пользователь' },
  custom: { en: 'Custom', ru: 'Другая роль' },
} as const;

export default function ProjectCredentials({ credentials }: Props) {
  const locale = useLocale();
  const isRu = locale === 'ru';
  const [copiedKey, setCopiedKey] = useState<string | null>(null);

  if (!credentials || credentials.length === 0) return null;

  const copy = async (value: string, key: string) => {
    try {
      await navigator.clipboard.writeText(value);
      setCopiedKey(key);
      setTimeout(() => setCopiedKey((k) => (k === key ? null : k)), 1200);
    } catch {
    }
  };

  return (
    <section className="glass-card rounded-xl overflow-hidden">
      <header className="flex items-center gap-2 px-5 py-3 md:py-3.5 border-b border-border/60 bg-black/20 font-mono text-xs md:text-sm">
        <Terminal className="h-3.5 w-3.5 md:h-4 md:w-4 text-emerald-400" />
        <span className="text-muted-foreground">
          {isRu ? 'Тестовые доступы — скопируй и попробуй' : 'Demo access — copy and try'}
        </span>
      </header>

      <div className="divide-y divide-border/40">
        {credentials.map((cred, i) => {
          const Icon = ROLE_ICONS[cred.role] ?? Settings;
          const roleLabel = cred.label && (isRu ? cred.label.ru : cred.label.en)
            ? (isRu ? cred.label.ru : cred.label.en)
            : (isRu ? ROLE_LABELS[cred.role].ru : ROLE_LABELS[cred.role].en);
          const note = cred.note && (isRu ? cred.note.ru : cred.note.en);
          const keyPrefix = `${i}-${cred.role}`;

          return (
            <div key={keyPrefix} className="p-5 md:p-6 font-mono text-sm md:text-base">
              <div className="flex items-center gap-2 mb-3 md:mb-4">
                <Icon className="h-4 w-4 md:h-5 md:w-5 text-primary" />
                <span className="text-xs md:text-sm uppercase tracking-[0.18em] text-primary">
                  {roleLabel}
                </span>
              </div>

              <div className="space-y-1.5 md:space-y-2">
                <CredRow
                  label="login"
                  value={cred.login}
                  copied={copiedKey === `${keyPrefix}-login`}
                  onCopy={() => copy(cred.login, `${keyPrefix}-login`)}
                />
                <CredRow
                  label="password"
                  value={cred.password}
                  copied={copiedKey === `${keyPrefix}-password`}
                  onCopy={() => copy(cred.password, `${keyPrefix}-password`)}
                  masked
                />
              </div>

              {note && (
                <p className="mt-3 text-xs md:text-sm text-muted-foreground font-sans leading-relaxed">
                  {note}
                </p>
              )}
            </div>
          );
        })}
      </div>

      <footer className="px-5 py-2 md:py-2.5 border-t border-border/60 bg-black/20 text-[10px] md:text-xs font-mono text-muted-foreground/70">
        {isRu ? 'Данные могут обнуляться при перезапусках демо' : 'Demo data may reset on redeploys'}
      </footer>
    </section>
  );
}

function CredRow({
  label,
  value,
  copied,
  onCopy,
  masked,
}: {
  label: string;
  value: string;
  copied: boolean;
  onCopy: () => void;
  masked?: boolean;
}) {
  const [reveal, setReveal] = useState(false);
  const displayed = masked && !reveal ? '•'.repeat(Math.max(6, value.length)) : value;

  return (
    <div className="flex items-center gap-3 text-emerald-300/90">
      <span className="w-20 md:w-24 shrink-0 text-muted-foreground/60">{label}</span>
      <span className="flex-1 truncate select-all">{displayed}</span>
      {masked && (
        <button
          type="button"
          onClick={() => setReveal((v) => !v)}
          className="text-[10px] md:text-xs uppercase tracking-widest text-muted-foreground hover:text-foreground transition-colors"
        >
          {reveal ? 'hide' : 'show'}
        </button>
      )}
      <button
        type="button"
        onClick={onCopy}
        className="p-1 md:p-1.5 rounded text-muted-foreground hover:text-primary transition-colors"
        aria-label={`Copy ${label}`}
      >
        {copied ? <Check className="h-3.5 w-3.5 md:h-4 md:w-4 text-emerald-400" /> : <Copy className="h-3.5 w-3.5 md:h-4 md:w-4" />}
      </button>
    </div>
  );
}
