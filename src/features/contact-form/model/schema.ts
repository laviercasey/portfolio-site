import { z } from 'zod';
import type { ContactFormField } from '@/entities/content';

export function buildContactSchema(fields: ContactFormField[]) {
  const shape: Record<string, z.ZodTypeAny> = {};
  for (const f of fields) {
    let s: z.ZodTypeAny;
    if (f.type === 'email') {
      s = f.required
        ? z.string().email().max(320)
        : z.string().email().max(320).or(z.literal('')).optional();
    } else if (f.type === 'textarea') {
      s = f.required
        ? z.string().min(10).max(5000)
        : z.string().max(5000).optional();
    } else {
      s = f.required
        ? z.string().min(1).max(500)
        : z.string().max(500).optional();
    }
    shape[f.id] = s;
  }
  return z.object(shape);
}
