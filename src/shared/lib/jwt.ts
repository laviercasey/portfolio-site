import { jwtVerify } from 'jose';

export interface AdminClaims {
  sub: string;
  role: string;
  exp: number;
  iat: number;
}

function getSecret(): Uint8Array {
  const secret = process.env.JWT_SECRET;
  if (!secret || secret.length < 32) {
    throw new Error('JWT_SECRET is not configured or is too short (min 32 chars)');
  }
  return new TextEncoder().encode(secret);
}

export async function verifyAdminToken(token: string): Promise<AdminClaims | null> {
  if (!token || token.split('.').length !== 3) return null;
  try {
    const { payload } = await jwtVerify(token, getSecret(), {
      algorithms: ['HS256'],
    });
    if (payload.role !== 'admin') return null;
    return payload as unknown as AdminClaims;
  } catch {
    return null;
  }
}
