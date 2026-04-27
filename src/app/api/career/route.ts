import { NextResponse } from 'next/server';
import { sanitizeBackendError } from '@/shared/api';

const BACKEND = process.env.API_INTERNAL_URL || 'http://localhost:8080';

export async function GET() {
  const res = await fetch(`${BACKEND}/api/career`, {
    next: { revalidate: 60 },
  });
  if (!res.ok) return sanitizeBackendError(res.status);
  const data = await res.json();
  return NextResponse.json(data);
}
