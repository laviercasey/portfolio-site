import createMiddleware from 'next-intl/middleware';
import { NextRequest, NextResponse } from 'next/server';
import { routing } from '@/shared/config';
import { verifyAdminToken } from '@/shared/lib/server';
import { getAdminSlug } from '@/shared/lib';

const handleI18nRouting = createMiddleware(routing);

export default async function middleware(request: NextRequest) {
  const pathname = request.nextUrl.pathname;
  const adminSlug = getAdminSlug();

  if (adminSlug) {
    if (pathname === `/${adminSlug}` || pathname.startsWith(`/${adminSlug}/`)) {
      const rest = pathname.slice(adminSlug.length + 1);

      if (rest && rest !== '/') {
        const token = request.cookies.get('auth-token')?.value?.trim();
        const claims = token ? await verifyAdminToken(token) : null;
        if (!claims) {
          return NextResponse.redirect(new URL(`/${adminSlug}`, request.url));
        }
      }

      const internalPath = `/adm${rest || ''}`;
      return NextResponse.rewrite(new URL(internalPath, request.url));
    }
  }

  if (pathname === '/adm' || pathname.startsWith('/adm/')) {
    return new NextResponse(null, { status: 404 });
  }

  if (pathname.startsWith('/api/')) {
    return NextResponse.next();
  }

  return handleI18nRouting(request);
}

export const config = {
  matcher: [
    '/((?!_next|_vercel|.*\\..*).*)',
  ],
};
