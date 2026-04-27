'use client';

import { Input } from './input';
import { Textarea } from './textarea';
import { Label } from './label';

interface I18nValue {
  en: string;
  ru: string;
}

interface Props {
  label: string;
  value: I18nValue | undefined;
  onChange: (next: I18nValue) => void;
  placeholder?: string;
  multiline?: boolean;
  rows?: number;
}

export default function I18nField({ label, value, onChange, placeholder, multiline, rows = 2 }: Props) {
  const safe: I18nValue = { en: value?.en ?? '', ru: value?.ru ?? '' };
  const Component = multiline ? Textarea : Input;

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
      <div className="space-y-1">
        <Label className="text-xs">{label} (RU)</Label>
        <Component
          value={safe.ru}
          rows={multiline ? rows : undefined}
          onChange={(e) => onChange({ ...safe, ru: e.target.value })}
          placeholder={placeholder}
        />
      </div>
      <div className="space-y-1">
        <Label className="text-xs">{label} (EN)</Label>
        <Component
          value={safe.en}
          rows={multiline ? rows : undefined}
          onChange={(e) => onChange({ ...safe, en: e.target.value })}
          placeholder={placeholder}
        />
      </div>
    </div>
  );
}
