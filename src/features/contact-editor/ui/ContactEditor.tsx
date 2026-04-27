'use client';

import { useEffect, useState } from 'react';
import { Button } from '@/shared/ui';
import { Input } from '@/shared/ui';
import { Textarea } from '@/shared/ui';
import { Label } from '@/shared/ui';
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/shared/ui';
import { Loader2, Save, Check, Plus, Trash2, ChevronUp, ChevronDown, GripVertical } from 'lucide-react';
import type { ContactPageConfig, ContactFormField } from '@/entities/content';

const FIELD_TYPES = ['text', 'email', 'textarea', 'select'] as const;
const GRID_OPTIONS = [
  { value: '1', label: 'Половина (1/2)' },
  { value: '2', label: 'Полная ширина' },
];

function emptyField(): ContactFormField {
  return {
    id: `field_${Date.now()}`,
    type: 'text',
    labelRu: '',
    labelEn: '',
    placeholderRu: '',
    placeholderEn: '',
    required: false,
    gridCol: 2,
  };
}

export default function ContactEditor() {
  const [config, setConfig] = useState<ContactPageConfig | null>(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [saved, setSaved] = useState(false);

  useEffect(() => {
    fetch('/api/content')
      .then((r) => r.json())
      .then((sections: Array<{ section: string; data: unknown }>) => {
        const found = sections.find((s) => s.section === 'contact');
        setConfig(found ? (found.data as ContactPageConfig) : null);
        setLoading(false);
      })
      .catch(() => setLoading(false));
  }, []);

  const handleSave = async () => {
    setSaving(true);
    const res = await fetch('/api/content?section=contact', {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(config),
    });
    setSaving(false);
    if (res.ok) { setSaved(true); setTimeout(() => setSaved(false), 2000); }
  };

  if (loading || !config) {
    return <div className="flex justify-center py-20"><Loader2 className="h-6 w-6 animate-spin" /></div>;
  }

  const updateField = (idx: number, patch: Partial<ContactFormField>) => {
    setConfig((prev) => {
      if (!prev) return prev;
      const fields = [...prev.formFields];
      fields[idx] = { ...fields[idx], ...patch };
      return { ...prev, formFields: fields };
    });
  };

  const addField = () => {
    setConfig((prev) => prev ? { ...prev, formFields: [...prev.formFields, emptyField()] } : prev);
  };

  const removeField = (idx: number) => {
    setConfig((prev) => prev ? { ...prev, formFields: prev.formFields.filter((_, i) => i !== idx) } : prev);
  };

  const moveField = (idx: number, dir: -1 | 1) => {
    setConfig((prev) => {
      if (!prev) return prev;
      const fields = [...prev.formFields];
      const target = idx + dir;
      if (target < 0 || target >= fields.length) return prev;
      [fields[idx], fields[target]] = [fields[target], fields[idx]];
      return { ...prev, formFields: fields };
    });
  };

  const updateOption = (fieldIdx: number, optIdx: number, patch: Partial<{ value: string; labelRu: string; labelEn: string }>) => {
    setConfig((prev) => {
      if (!prev) return prev;
      const fields = [...prev.formFields];
      const options = [...(fields[fieldIdx].options ?? [])];
      options[optIdx] = { ...options[optIdx], ...patch };
      fields[fieldIdx] = { ...fields[fieldIdx], options };
      return { ...prev, formFields: fields };
    });
  };

  const addOption = (fieldIdx: number) => {
    setConfig((prev) => {
      if (!prev) return prev;
      const fields = [...prev.formFields];
      const options = [...(fields[fieldIdx].options ?? []), { value: '', labelRu: '', labelEn: '' }];
      fields[fieldIdx] = { ...fields[fieldIdx], options };
      return { ...prev, formFields: fields };
    });
  };

  const removeOption = (fieldIdx: number, optIdx: number) => {
    setConfig((prev) => {
      if (!prev) return prev;
      const fields = [...prev.formFields];
      const options = (fields[fieldIdx].options ?? []).filter((_, i) => i !== optIdx);
      fields[fieldIdx] = { ...fields[fieldIdx], options };
      return { ...prev, formFields: fields };
    });
  };

  return (
    <div className="p-8 max-w-4xl">
      <div className="flex items-center justify-between mb-8">
        <h1 className="text-3xl font-bold">Contact Page Settings</h1>
        <Button onClick={handleSave} disabled={saving}>
          {saving ? <Loader2 className="h-4 w-4 animate-spin mr-2" /> : saved ? <Check className="h-4 w-4 mr-2" /> : <Save className="h-4 w-4 mr-2" />}
          {saved ? 'Saved!' : 'Save'}
        </Button>
      </div>

      <Card className="mb-6">
        <CardHeader><CardTitle>Page Texts</CardTitle></CardHeader>
        <CardContent className="space-y-4">
          <div className="grid grid-cols-2 gap-4">
            <div>
              <Label>Heading (RU)</Label>
              <Input value={config.heading.ru} onChange={(e) => setConfig({ ...config, heading: { ...config.heading, ru: e.target.value } })} />
            </div>
            <div>
              <Label>Heading (EN)</Label>
              <Input value={config.heading.en} onChange={(e) => setConfig({ ...config, heading: { ...config.heading, en: e.target.value } })} />
            </div>
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <Label>Subtitle (RU)</Label>
              <Input value={config.subtitle.ru} onChange={(e) => setConfig({ ...config, subtitle: { ...config.subtitle, ru: e.target.value } })} />
            </div>
            <div>
              <Label>Subtitle (EN)</Label>
              <Input value={config.subtitle.en} onChange={(e) => setConfig({ ...config, subtitle: { ...config.subtitle, en: e.target.value } })} />
            </div>
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <Label>How I Work Text (RU)</Label>
              <Textarea rows={3} value={config.howIWork.ru} onChange={(e) => setConfig({ ...config, howIWork: { ...config.howIWork, ru: e.target.value } })} />
            </div>
            <div>
              <Label>How I Work Text (EN)</Label>
              <Textarea rows={3} value={config.howIWork.en} onChange={(e) => setConfig({ ...config, howIWork: { ...config.howIWork, en: e.target.value } })} />
            </div>
          </div>
        </CardContent>
      </Card>

      <Card className="mb-6">
        <CardHeader><CardTitle>Submit Button</CardTitle></CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <Label>Button Text (RU)</Label>
              <Input value={config.submitTextRu} onChange={(e) => setConfig({ ...config, submitTextRu: e.target.value })} />
            </div>
            <div>
              <Label>Button Text (EN)</Label>
              <Input value={config.submitTextEn} onChange={(e) => setConfig({ ...config, submitTextEn: e.target.value })} />
            </div>
          </div>
        </CardContent>
      </Card>

      <Card className="mb-6">
        <CardHeader className="flex flex-row items-center justify-between">
          <CardTitle>Form Fields</CardTitle>
          <Button size="sm" variant="outline" onClick={addField}>
            <Plus className="h-4 w-4 mr-1" /> Add Field
          </Button>
        </CardHeader>
        <CardContent className="space-y-6">
          {config.formFields.map((field, idx) => (
            <div key={field.id} className="border rounded-lg p-4 space-y-4">
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <GripVertical className="h-4 w-4 text-muted-foreground" />
                  <span className="font-medium text-sm">#{idx + 1} — {field.id}</span>
                </div>
                <div className="flex items-center gap-1">
                  <Button size="icon" variant="ghost" onClick={() => moveField(idx, -1)} disabled={idx === 0}>
                    <ChevronUp className="h-4 w-4" />
                  </Button>
                  <Button size="icon" variant="ghost" onClick={() => moveField(idx, 1)} disabled={idx === config.formFields.length - 1}>
                    <ChevronDown className="h-4 w-4" />
                  </Button>
                  <Button size="icon" variant="ghost" className="text-destructive" onClick={() => removeField(idx)}>
                    <Trash2 className="h-4 w-4" />
                  </Button>
                </div>
              </div>

              <div className="grid grid-cols-4 gap-3">
                <div>
                  <Label className="text-xs">Field ID</Label>
                  <Input value={field.id} onChange={(e) => updateField(idx, { id: e.target.value })} />
                </div>
                <div>
                  <Label className="text-xs">Type</Label>
                  <Select value={field.type} onValueChange={(v) => updateField(idx, { type: v as ContactFormField['type'] })}>
                    <SelectTrigger><SelectValue /></SelectTrigger>
                    <SelectContent>
                      {FIELD_TYPES.map((t) => <SelectItem key={t} value={t}>{t}</SelectItem>)}
                    </SelectContent>
                  </Select>
                </div>
                <div>
                  <Label className="text-xs">Width</Label>
                  <Select value={String(field.gridCol ?? 2)} onValueChange={(v) => updateField(idx, { gridCol: Number(v) as 1 | 2 })}>
                    <SelectTrigger><SelectValue /></SelectTrigger>
                    <SelectContent>
                      {GRID_OPTIONS.map((o) => <SelectItem key={o.value} value={o.value}>{o.label}</SelectItem>)}
                    </SelectContent>
                  </Select>
                </div>
                <div className="flex items-end gap-2 pb-1">
                  <input type="checkbox" checked={field.required} onChange={(e) => updateField(idx, { required: e.target.checked })} className="h-4 w-4 rounded border-gray-300" />
                  <Label className="text-xs">Required</Label>
                </div>
              </div>

              <div className="grid grid-cols-2 gap-3">
                <div>
                  <Label className="text-xs">Label (RU)</Label>
                  <Input value={field.labelRu} onChange={(e) => updateField(idx, { labelRu: e.target.value })} />
                </div>
                <div>
                  <Label className="text-xs">Label (EN)</Label>
                  <Input value={field.labelEn} onChange={(e) => updateField(idx, { labelEn: e.target.value })} />
                </div>
              </div>

              <div className="grid grid-cols-2 gap-3">
                <div>
                  <Label className="text-xs">Placeholder (RU)</Label>
                  <Input value={field.placeholderRu ?? ''} onChange={(e) => updateField(idx, { placeholderRu: e.target.value })} />
                </div>
                <div>
                  <Label className="text-xs">Placeholder (EN)</Label>
                  <Input value={field.placeholderEn ?? ''} onChange={(e) => updateField(idx, { placeholderEn: e.target.value })} />
                </div>
              </div>

              {field.type === 'select' && (
                <div className="pl-4 border-l-2 border-primary/20 space-y-3">
                  <div className="flex items-center justify-between">
                    <Label className="text-xs font-medium">Select Options</Label>
                    <Button size="sm" variant="outline" onClick={() => addOption(idx)}>
                      <Plus className="h-3 w-3 mr-1" /> Option
                    </Button>
                  </div>
                  {(field.options ?? []).map((opt, oi) => (
                    <div key={oi} className="grid grid-cols-[1fr_1fr_1fr_auto] gap-2 items-end">
                      <div>
                        <Label className="text-xs">Value</Label>
                        <Input value={opt.value} onChange={(e) => updateOption(idx, oi, { value: e.target.value })} />
                      </div>
                      <div>
                        <Label className="text-xs">Label RU</Label>
                        <Input value={opt.labelRu} onChange={(e) => updateOption(idx, oi, { labelRu: e.target.value })} />
                      </div>
                      <div>
                        <Label className="text-xs">Label EN</Label>
                        <Input value={opt.labelEn} onChange={(e) => updateOption(idx, oi, { labelEn: e.target.value })} />
                      </div>
                      <Button size="icon" variant="ghost" className="text-destructive" onClick={() => removeOption(idx, oi)}>
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    </div>
                  ))}
                </div>
              )}
            </div>
          ))}
        </CardContent>
      </Card>
    </div>
  );
}
