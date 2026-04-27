import { NextRequest, NextResponse } from 'next/server';
import { getAuthToken } from '@/shared/lib/server';
import { sanitizeBackendError } from '@/shared/api';

const BACKEND = process.env.API_INTERNAL_URL || 'http://localhost:8080';

const ALLOWED_MIME_TYPES = [
  'image/jpeg', 'image/png', 'image/webp', 'image/gif',
  'video/mp4', 'video/webm',
];
const MAX_FILE_SIZE = 50 * 1024 * 1024;

export async function POST(request: NextRequest) {
  const token = await getAuthToken();
  if (!token) return NextResponse.json({ error: 'Unauthorized' }, { status: 401 });

  const formData = await request.formData().catch(() => null);
  if (!formData) return NextResponse.json({ error: 'Invalid form data' }, { status: 400 });

  const file = formData.get('file');
  if (!(file instanceof File)) {
    return NextResponse.json({ error: 'No file provided' }, { status: 400 });
  }

  if (!ALLOWED_MIME_TYPES.includes(file.type)) {
    return NextResponse.json({ error: 'File type not allowed' }, { status: 400 });
  }

  if (file.size > MAX_FILE_SIZE) {
    return NextResponse.json({ error: 'File too large' }, { status: 400 });
  }

  const cleanFormData = new FormData();
  cleanFormData.append('file', file);

  const res = await fetch(`${BACKEND}/api/upload`, {
    method: 'POST',
    headers: { Authorization: `Bearer ${token}` },
    body: cleanFormData,
  });

  if (!res.ok) return sanitizeBackendError(res.status);
  const data = await res.json();
  return NextResponse.json(data);
}
