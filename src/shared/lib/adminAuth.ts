import { cookies } from 'next/headers';
import { NextResponse } from 'next/server';
import { verifyAdminToken } from './jwt';

export async function getAuthToken(): Promise<string | null> {
  const cookieStore = await cookies();
  const value = cookieStore.get('auth-token')?.value?.trim() ?? null;
  if (!value) return null;
  const claims = await verifyAdminToken(value);
  if (!claims) return null;
  return value;
}

export async function requireAdmin(): Promise<{ token: string } | { error: NextResponse }> {
  const token = await getAuthToken();
  if (!token) {
    return { error: NextResponse.json({ error: 'Unauthorized' }, { status: 401 }) };
  }
  return { token };
}
