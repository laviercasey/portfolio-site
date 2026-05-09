'use client';

import { useEffect, useState } from 'react';
import { Loader2, Save, Plus, Trash2, ChevronDown, ChevronUp, X } from 'lucide-react';
import { Button, Input, Textarea, Label, Card, CardContent, CardHeader, CardTitle } from '@/shared/ui';
import type {
  Service,
  ServiceFaq,
  ServiceProcessStep,
  ServiceIconKey,
  ServiceVisualKey,
  ServiceCaseProject,
  CreateServiceInput,
  CreateServiceFaqInput,
  CreateServiceProcessStepInput,
} from '@/entities/service';

type Tab = 'services' | 'faqs' | 'process';

const ICON_KEYS: ServiceIconKey[] = ['bot', 'code', 'layers', 'database', 'workflow'];
const VISUAL_KEYS: ServiceVisualKey[] = ['terminal', 'browser', 'editor', 'dataframe', 'pipeline'];

async function api<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(path, {
    ...init,
    headers: { 'Content-Type': 'application/json', ...(init?.headers ?? {}) },
  });
  if (!res.ok) {
    const text = await res.text().catch(() => '');
    throw new Error(text || `${res.status} ${res.statusText}`);
  }
  if (res.status === 204) return undefined as T;
  return res.json() as Promise<T>;
}

export default function ServicesEditor() {
  const [tab, setTab] = useState<Tab>('services');
  const [services, setServices] = useState<Service[]>([]);
  const [faqs, setFaqs] = useState<ServiceFaq[]>([]);
  const [steps, setSteps] = useState<ServiceProcessStep[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const reload = async () => {
    setLoading(true);
    setError(null);
    try {
      const [s, f, p] = await Promise.all([
        api<Service[]>('/api/services/list').catch(() => []),
        api<ServiceFaq[]>('/api/services/faqs').catch(() => []),
        api<ServiceProcessStep[]>('/api/services/process').catch(() => []),
      ]);
      setServices(s ?? []);
      setFaqs(f ?? []);
      setSteps(p ?? []);
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to load');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { reload(); }, []);

  if (loading) {
    return <div className="flex items-center justify-center py-20"><Loader2 className="h-8 w-8 animate-spin" /></div>;
  }

  return (
    <div className="space-y-6">
      <div className="flex gap-2 border-b border-white/10">
        {(['services', 'faqs', 'process'] as Tab[]).map((t) => (
          <button
            key={t}
            onClick={() => setTab(t)}
            className={`px-4 py-2 text-sm font-medium border-b-2 transition-colors ${
              tab === t ? 'border-primary text-foreground' : 'border-transparent text-muted-foreground hover:text-foreground'
            }`}
          >
            {t === 'services' ? `Услуги (${services.length})` : t === 'faqs' ? `FAQ (${faqs.length})` : `Этапы (${steps.length})`}
          </button>
        ))}
      </div>

      {error && <div className="text-sm text-red-400 px-3 py-2 bg-red-500/10 border border-red-500/30 rounded-lg">{error}</div>}

      {tab === 'services' && <ServicesTab items={services} onChange={reload} />}
      {tab === 'faqs' && <FaqsTab items={faqs} onChange={reload} />}
      {tab === 'process' && <ProcessTab items={steps} onChange={reload} />}
    </div>
  );
}

function ServicesTab({ items, onChange }: { items: Service[]; onChange: () => void }) {
  const [openId, setOpenId] = useState<string | null>(null);
  const [creating, setCreating] = useState(false);

  return (
    <div className="space-y-3">
      <Button onClick={() => setCreating(true)} className="gap-2"><Plus className="h-4 w-4" />Добавить услугу</Button>

      {creating && (
        <ServiceForm
          onCancel={() => setCreating(false)}
          onSave={async (input) => {
            await api('/api/services', { method: 'POST', body: JSON.stringify(input) });
            setCreating(false);
            onChange();
          }}
        />
      )}

      {items.map((s) => (
        <Card key={s.id}>
          <CardHeader className="cursor-pointer" onClick={() => setOpenId(openId === s.id ? null : s.id)}>
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                <span className="font-mono text-xs px-2 py-0.5 rounded" style={{ background: `${s.accent}22`, color: s.accent }}>{s.num}</span>
                <CardTitle className="text-base">{s.title.ru || s.title.en || s.slug}</CardTitle>
                <span className="text-xs text-muted-foreground">от {s.priceRu}</span>
              </div>
              {openId === s.id ? <ChevronUp className="h-4 w-4" /> : <ChevronDown className="h-4 w-4" />}
            </div>
          </CardHeader>
          {openId === s.id && (
            <CardContent>
              <ServiceForm
                initial={s}
                onCancel={() => setOpenId(null)}
                onSave={async (input) => {
                  await api(`/api/services/${s.id}`, { method: 'PUT', body: JSON.stringify(input) });
                  setOpenId(null);
                  onChange();
                }}
                onDelete={async () => {
                  if (!confirm(`Удалить услугу "${s.title.ru}"?`)) return;
                  await api(`/api/services/${s.id}`, { method: 'DELETE' });
                  setOpenId(null);
                  onChange();
                }}
              />
            </CardContent>
          )}
        </Card>
      ))}
    </div>
  );
}

function ServiceForm({
  initial,
  onSave,
  onCancel,
  onDelete,
}: {
  initial?: Service;
  onSave: (input: CreateServiceInput) => Promise<void>;
  onCancel: () => void;
  onDelete?: () => Promise<void>;
}) {
  const [form, setForm] = useState<CreateServiceInput>(() => ({
    slug: initial?.slug ?? '',
    num: initial?.num ?? '',
    iconKey: initial?.iconKey ?? 'bot',
    visualKey: initial?.visualKey ?? 'terminal',
    accent: initial?.accent ?? '#5eb3ff',
    title: initial?.title ?? { ru: '', en: '' },
    lead: initial?.lead ?? { ru: '', en: '' },
    bullets: initial?.bullets ?? { ru: [], en: [] },
    stack: initial?.stack ?? '',
    timeline: initial?.timeline ?? { ru: '', en: '' },
    priceRu: initial?.priceRu ?? '',
    priceEn: initial?.priceEn ?? '',
    caseProjects: initial?.caseProjects ?? [],
    order: initial?.order ?? 0,
  }));
  const [saving, setSaving] = useState(false);
  const [err, setErr] = useState<string | null>(null);

  const submit = async () => {
    setSaving(true);
    setErr(null);
    try {
      await onSave(form);
    } catch (e) {
      setErr(e instanceof Error ? e.message : 'Save failed');
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="space-y-4">
      <div className="grid grid-cols-2 gap-3">
        <Field label="Slug"><Input value={form.slug} onChange={(e) => setForm({ ...form, slug: e.target.value })} /></Field>
        <Field label="Номер (01)"><Input value={form.num} onChange={(e) => setForm({ ...form, num: e.target.value })} /></Field>
        <Field label="Иконка">
          <select className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm" value={form.iconKey} onChange={(e) => setForm({ ...form, iconKey: e.target.value as ServiceIconKey })}>
            {ICON_KEYS.map((k) => <option key={k} value={k}>{k}</option>)}
          </select>
        </Field>
        <Field label="Визуал">
          <select className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm" value={form.visualKey} onChange={(e) => setForm({ ...form, visualKey: e.target.value as ServiceVisualKey })}>
            {VISUAL_KEYS.map((k) => <option key={k} value={k}>{k}</option>)}
          </select>
        </Field>
        <Field label="Акцент-цвет (#hex)"><Input value={form.accent} onChange={(e) => setForm({ ...form, accent: e.target.value })} /></Field>
        <Field label="Сортировка"><Input type="number" value={form.order} onChange={(e) => setForm({ ...form, order: parseInt(e.target.value) || 0 })} /></Field>
      </div>

      <div className="grid grid-cols-2 gap-3">
        <Field label="Заголовок RU"><Input value={form.title.ru} onChange={(e) => setForm({ ...form, title: { ...form.title, ru: e.target.value } })} /></Field>
        <Field label="Title EN"><Input value={form.title.en} onChange={(e) => setForm({ ...form, title: { ...form.title, en: e.target.value } })} /></Field>
      </div>

      <div className="grid grid-cols-2 gap-3">
        <Field label="Лид RU"><Textarea rows={4} value={form.lead.ru} onChange={(e) => setForm({ ...form, lead: { ...form.lead, ru: e.target.value } })} /></Field>
        <Field label="Lead EN"><Textarea rows={4} value={form.lead.en} onChange={(e) => setForm({ ...form, lead: { ...form.lead, en: e.target.value } })} /></Field>
      </div>

      <BulletsField label="Пункты RU" items={form.bullets.ru} onChange={(v) => setForm({ ...form, bullets: { ...form.bullets, ru: v } })} />
      <BulletsField label="Bullets EN" items={form.bullets.en} onChange={(v) => setForm({ ...form, bullets: { ...form.bullets, en: v } })} />

      <Field label="Стек (через ' · ')"><Input value={form.stack} onChange={(e) => setForm({ ...form, stack: e.target.value })} /></Field>

      <div className="grid grid-cols-2 gap-3">
        <Field label="Сроки RU"><Input value={form.timeline.ru} onChange={(e) => setForm({ ...form, timeline: { ...form.timeline, ru: e.target.value } })} /></Field>
        <Field label="Timeline EN"><Input value={form.timeline.en} onChange={(e) => setForm({ ...form, timeline: { ...form.timeline, en: e.target.value } })} /></Field>
        <Field label="Цена RU"><Input value={form.priceRu} onChange={(e) => setForm({ ...form, priceRu: e.target.value })} /></Field>
        <Field label="Price EN"><Input value={form.priceEn} onChange={(e) => setForm({ ...form, priceEn: e.target.value })} /></Field>
      </div>

      <CaseProjectsField items={form.caseProjects} onChange={(v) => setForm({ ...form, caseProjects: v })} />

      {err && <div className="text-sm text-red-400">{err}</div>}

      <div className="flex items-center justify-between pt-2">
        <div className="flex gap-2">
          <Button onClick={submit} disabled={saving} className="gap-2">
            {saving ? <Loader2 className="h-4 w-4 animate-spin" /> : <Save className="h-4 w-4" />}
            Сохранить
          </Button>
          <Button variant="outline" onClick={onCancel}>Отмена</Button>
        </div>
        {onDelete && (
          <Button variant="outline" onClick={onDelete} className="text-red-400 border-red-400/30 hover:bg-red-500/10 gap-2">
            <Trash2 className="h-4 w-4" />Удалить
          </Button>
        )}
      </div>
    </div>
  );
}

function BulletsField({ label, items, onChange }: { label: string; items: string[]; onChange: (v: string[]) => void }) {
  return (
    <div>
      <Label className="mb-2 block">{label}</Label>
      <div className="space-y-2">
        {items.map((b, i) => (
          <div key={i} className="flex gap-2">
            <Input value={b} onChange={(e) => {
              const next = [...items];
              next[i] = e.target.value;
              onChange(next);
            }} />
            <Button variant="outline" size="icon" onClick={() => onChange(items.filter((_, idx) => idx !== i))}>
              <X className="h-4 w-4" />
            </Button>
          </div>
        ))}
        <Button variant="outline" size="sm" onClick={() => onChange([...items, ''])} className="gap-2">
          <Plus className="h-3.5 w-3.5" />Добавить
        </Button>
      </div>
    </div>
  );
}

function CaseProjectsField({ items, onChange }: { items: ServiceCaseProject[]; onChange: (v: ServiceCaseProject[]) => void }) {
  return (
    <div>
      <Label className="mb-2 block">Кейсы (slug + name)</Label>
      <div className="space-y-2">
        {items.map((c, i) => (
          <div key={i} className="flex gap-2">
            <Input placeholder="slug" value={c.slug} onChange={(e) => {
              const next = [...items];
              next[i] = { ...next[i], slug: e.target.value };
              onChange(next);
            }} />
            <Input placeholder="Name" value={c.name} onChange={(e) => {
              const next = [...items];
              next[i] = { ...next[i], name: e.target.value };
              onChange(next);
            }} />
            <Button variant="outline" size="icon" onClick={() => onChange(items.filter((_, idx) => idx !== i))}>
              <X className="h-4 w-4" />
            </Button>
          </div>
        ))}
        <Button variant="outline" size="sm" onClick={() => onChange([...items, { slug: '', name: '' }])} className="gap-2">
          <Plus className="h-3.5 w-3.5" />Добавить кейс
        </Button>
      </div>
    </div>
  );
}

function FaqsTab({ items, onChange }: { items: ServiceFaq[]; onChange: () => void }) {
  const [openId, setOpenId] = useState<string | null>(null);
  const [creating, setCreating] = useState(false);

  return (
    <div className="space-y-3">
      <Button onClick={() => setCreating(true)} className="gap-2"><Plus className="h-4 w-4" />Добавить FAQ</Button>

      {creating && (
        <FaqForm
          onCancel={() => setCreating(false)}
          onSave={async (input) => {
            await api('/api/services/faqs', { method: 'POST', body: JSON.stringify(input) });
            setCreating(false);
            onChange();
          }}
        />
      )}

      {items.map((f) => (
        <Card key={f.id}>
          <CardHeader className="cursor-pointer" onClick={() => setOpenId(openId === f.id ? null : f.id)}>
            <div className="flex items-center justify-between">
              <CardTitle className="text-base">{f.question.ru || f.question.en}</CardTitle>
              {openId === f.id ? <ChevronUp className="h-4 w-4" /> : <ChevronDown className="h-4 w-4" />}
            </div>
          </CardHeader>
          {openId === f.id && (
            <CardContent>
              <FaqForm
                initial={f}
                onCancel={() => setOpenId(null)}
                onSave={async (input) => {
                  await api(`/api/services/faqs/${f.id}`, { method: 'PUT', body: JSON.stringify(input) });
                  setOpenId(null);
                  onChange();
                }}
                onDelete={async () => {
                  if (!confirm('Удалить вопрос?')) return;
                  await api(`/api/services/faqs/${f.id}`, { method: 'DELETE' });
                  setOpenId(null);
                  onChange();
                }}
              />
            </CardContent>
          )}
        </Card>
      ))}
    </div>
  );
}

function FaqForm({
  initial,
  onSave,
  onCancel,
  onDelete,
}: {
  initial?: ServiceFaq;
  onSave: (input: CreateServiceFaqInput) => Promise<void>;
  onCancel: () => void;
  onDelete?: () => Promise<void>;
}) {
  const [form, setForm] = useState<CreateServiceFaqInput>(() => ({
    question: initial?.question ?? { ru: '', en: '' },
    answer: initial?.answer ?? { ru: '', en: '' },
    order: initial?.order ?? 0,
  }));
  const [saving, setSaving] = useState(false);
  const [err, setErr] = useState<string | null>(null);

  const submit = async () => {
    setSaving(true);
    setErr(null);
    try {
      await onSave(form);
    } catch (e) {
      setErr(e instanceof Error ? e.message : 'Save failed');
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="space-y-3">
      <Field label="Сортировка"><Input type="number" value={form.order} onChange={(e) => setForm({ ...form, order: parseInt(e.target.value) || 0 })} /></Field>
      <div className="grid grid-cols-2 gap-3">
        <Field label="Вопрос RU"><Input value={form.question.ru} onChange={(e) => setForm({ ...form, question: { ...form.question, ru: e.target.value } })} /></Field>
        <Field label="Question EN"><Input value={form.question.en} onChange={(e) => setForm({ ...form, question: { ...form.question, en: e.target.value } })} /></Field>
      </div>
      <div className="grid grid-cols-2 gap-3">
        <Field label="Ответ RU"><Textarea rows={4} value={form.answer.ru} onChange={(e) => setForm({ ...form, answer: { ...form.answer, ru: e.target.value } })} /></Field>
        <Field label="Answer EN"><Textarea rows={4} value={form.answer.en} onChange={(e) => setForm({ ...form, answer: { ...form.answer, en: e.target.value } })} /></Field>
      </div>
      {err && <div className="text-sm text-red-400">{err}</div>}
      <div className="flex items-center justify-between pt-2">
        <div className="flex gap-2">
          <Button onClick={submit} disabled={saving} className="gap-2">
            {saving ? <Loader2 className="h-4 w-4 animate-spin" /> : <Save className="h-4 w-4" />}
            Сохранить
          </Button>
          <Button variant="outline" onClick={onCancel}>Отмена</Button>
        </div>
        {onDelete && (
          <Button variant="outline" onClick={onDelete} className="text-red-400 border-red-400/30 hover:bg-red-500/10 gap-2">
            <Trash2 className="h-4 w-4" />Удалить
          </Button>
        )}
      </div>
    </div>
  );
}

function ProcessTab({ items, onChange }: { items: ServiceProcessStep[]; onChange: () => void }) {
  const [openId, setOpenId] = useState<string | null>(null);
  const [creating, setCreating] = useState(false);

  return (
    <div className="space-y-3">
      <Button onClick={() => setCreating(true)} className="gap-2"><Plus className="h-4 w-4" />Добавить шаг</Button>

      {creating && (
        <ProcessForm
          onCancel={() => setCreating(false)}
          onSave={async (input) => {
            await api('/api/services/process', { method: 'POST', body: JSON.stringify(input) });
            setCreating(false);
            onChange();
          }}
        />
      )}

      {items.map((p) => (
        <Card key={p.id}>
          <CardHeader className="cursor-pointer" onClick={() => setOpenId(openId === p.id ? null : p.id)}>
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                <span className="font-mono text-xs px-2 py-0.5 rounded bg-primary/15 text-primary">{p.num}</span>
                <CardTitle className="text-base">{p.title.ru || p.title.en}</CardTitle>
              </div>
              {openId === p.id ? <ChevronUp className="h-4 w-4" /> : <ChevronDown className="h-4 w-4" />}
            </div>
          </CardHeader>
          {openId === p.id && (
            <CardContent>
              <ProcessForm
                initial={p}
                onCancel={() => setOpenId(null)}
                onSave={async (input) => {
                  await api(`/api/services/process/${p.id}`, { method: 'PUT', body: JSON.stringify(input) });
                  setOpenId(null);
                  onChange();
                }}
                onDelete={async () => {
                  if (!confirm('Удалить шаг?')) return;
                  await api(`/api/services/process/${p.id}`, { method: 'DELETE' });
                  setOpenId(null);
                  onChange();
                }}
              />
            </CardContent>
          )}
        </Card>
      ))}
    </div>
  );
}

function ProcessForm({
  initial,
  onSave,
  onCancel,
  onDelete,
}: {
  initial?: ServiceProcessStep;
  onSave: (input: CreateServiceProcessStepInput) => Promise<void>;
  onCancel: () => void;
  onDelete?: () => Promise<void>;
}) {
  const [form, setForm] = useState<CreateServiceProcessStepInput>(() => ({
    num: initial?.num ?? '',
    title: initial?.title ?? { ru: '', en: '' },
    description: initial?.description ?? { ru: '', en: '' },
    order: initial?.order ?? 0,
  }));
  const [saving, setSaving] = useState(false);
  const [err, setErr] = useState<string | null>(null);

  const submit = async () => {
    setSaving(true);
    setErr(null);
    try {
      await onSave(form);
    } catch (e) {
      setErr(e instanceof Error ? e.message : 'Save failed');
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="space-y-3">
      <div className="grid grid-cols-2 gap-3">
        <Field label="Номер (01)"><Input value={form.num} onChange={(e) => setForm({ ...form, num: e.target.value })} /></Field>
        <Field label="Сортировка"><Input type="number" value={form.order} onChange={(e) => setForm({ ...form, order: parseInt(e.target.value) || 0 })} /></Field>
      </div>
      <div className="grid grid-cols-2 gap-3">
        <Field label="Заголовок RU"><Input value={form.title.ru} onChange={(e) => setForm({ ...form, title: { ...form.title, ru: e.target.value } })} /></Field>
        <Field label="Title EN"><Input value={form.title.en} onChange={(e) => setForm({ ...form, title: { ...form.title, en: e.target.value } })} /></Field>
      </div>
      <div className="grid grid-cols-2 gap-3">
        <Field label="Описание RU"><Textarea rows={3} value={form.description.ru} onChange={(e) => setForm({ ...form, description: { ...form.description, ru: e.target.value } })} /></Field>
        <Field label="Description EN"><Textarea rows={3} value={form.description.en} onChange={(e) => setForm({ ...form, description: { ...form.description, en: e.target.value } })} /></Field>
      </div>
      {err && <div className="text-sm text-red-400">{err}</div>}
      <div className="flex items-center justify-between pt-2">
        <div className="flex gap-2">
          <Button onClick={submit} disabled={saving} className="gap-2">
            {saving ? <Loader2 className="h-4 w-4 animate-spin" /> : <Save className="h-4 w-4" />}
            Сохранить
          </Button>
          <Button variant="outline" onClick={onCancel}>Отмена</Button>
        </div>
        {onDelete && (
          <Button variant="outline" onClick={onDelete} className="text-red-400 border-red-400/30 hover:bg-red-500/10 gap-2">
            <Trash2 className="h-4 w-4" />Удалить
          </Button>
        )}
      </div>
    </div>
  );
}

function Field({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <div>
      <Label className="mb-1.5 block">{label}</Label>
      {children}
    </div>
  );
}
