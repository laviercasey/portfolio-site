'use client';

import { useState, useMemo } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import type { z } from 'zod';
import { useLocale, useTranslations } from 'next-intl';
import { Button, Input, Textarea, Label, Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/shared/ui';
import { CheckCircle2, Loader2 } from 'lucide-react';
import type { ContactFormField } from '@/entities/content';
import { buildContactSchema } from '../model/schema';
import { submitInquiry } from '../api/submitInquiry';

interface ContactFormProps {
  fields: ContactFormField[];
  submitText?: { ru: string; en: string };
}

export default function ContactForm({ fields, submitText }: ContactFormProps) {
  const locale = useLocale();
  const isRu = locale === 'ru';
  const t = useTranslations('contact');
  const [submitted, setSubmitted] = useState(false);
  const [error, setError] = useState('');

  const schema = useMemo(() => buildContactSchema(fields), [fields]);
  type FormValues = z.infer<typeof schema>;

  const defaults = useMemo(() => {
    const d: FormValues = {};
    for (const f of fields) {
      if (f.type === 'select' && f.options?.length) {
        d[f.id] = f.options[0].value;
      } else {
        d[f.id] = '';
      }
    }
    return d;
  }, [fields]);

  const { register, handleSubmit, setValue, formState: { errors, isSubmitting } } = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: defaults,
  });

  const onSubmit = async (data: FormValues) => {
    setError('');
    try {
      await submitInquiry(data);
      setSubmitted(true);
    } catch {
      setError(t('error'));
    }
  };

  if (submitted) {
    return (
      <div className="flex flex-col items-center justify-center py-12 text-center">
        <CheckCircle2 className="h-12 w-12 text-green-500 mb-4" />
        <p className="text-lg font-medium">
          {t('success')}
        </p>
      </div>
    );
  }

  const rows: ContactFormField[][] = [];
  let currentRow: ContactFormField[] = [];
  for (const f of fields) {
    if (f.gridCol === 2 || f.type === 'textarea') {
      if (currentRow.length > 0) { rows.push(currentRow); currentRow = []; }
      rows.push([f]);
    } else {
      currentRow.push(f);
      if (currentRow.length === 2) { rows.push(currentRow); currentRow = []; }
    }
  }
  if (currentRow.length > 0) rows.push(currentRow);

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
      {rows.map((row, ri) => (
        <div key={ri} className={row.length === 1 && (row[0].gridCol === 2 || row[0].type === 'textarea') ? '' : 'grid grid-cols-1 md:grid-cols-2 gap-5'}>
          {row.map((field) => (
            <div key={field.id} className="space-y-2">
              <Label htmlFor={field.id}>
                {isRu ? field.labelRu : field.labelEn}
                {field.required ? ' *' : ''}
                {!field.required && <span className="text-muted-foreground text-xs ml-1">({t('optional')})</span>}
              </Label>

              {field.type === 'select' && field.options ? (
                <Select
                  defaultValue={field.options[0]?.value}
                  onValueChange={(v) => setValue(field.id, v)}
                >
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {field.options.map((opt) => (
                      <SelectItem key={opt.value} value={opt.value}>
                        {isRu ? opt.labelRu : opt.labelEn}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              ) : field.type === 'textarea' ? (
                <Textarea
                  id={field.id}
                  rows={5}
                  {...register(field.id)}
                  placeholder={isRu ? field.placeholderRu : field.placeholderEn}
                />
              ) : (
                <Input
                  id={field.id}
                  type={field.type === 'email' ? 'email' : 'text'}
                  {...register(field.id)}
                  placeholder={isRu ? field.placeholderRu : field.placeholderEn}
                />
              )}

              {errors[field.id] && (
                <p className="text-xs text-destructive">
                  {field.type === 'email'
                    ? t('validEmail')
                    : field.type === 'textarea'
                    ? t('minChars')
                    : t('fieldRequired')}
                </p>
              )}
            </div>
          ))}
        </div>
      ))}

      {error && <p className="text-sm text-destructive">{error}</p>}

      <Button type="submit" size="lg" disabled={isSubmitting} className="w-full">
        {isSubmitting ? (
          <><Loader2 className="mr-2 h-4 w-4 animate-spin" />{t('sending')}</>
        ) : (
          submitText ? (isRu ? submitText.ru : submitText.en) : t('send')
        )}
      </Button>
    </form>
  );
}
