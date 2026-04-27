import type { I18nString } from '@/shared/types';

const EMPTY: I18nString = { en: '', ru: '' };

export function normalizeI18n(value: unknown): I18nString {
  if (!value || typeof value !== 'object') return { ...EMPTY };
  const v = value as Partial<I18nString>;
  return {
    en: typeof v.en === 'string' ? v.en : '',
    ru: typeof v.ru === 'string' ? v.ru : '',
  };
}

export function normalizeI18nOptional(value: unknown): I18nString | undefined {
  if (value === null || value === undefined) return undefined;
  return normalizeI18n(value);
}
