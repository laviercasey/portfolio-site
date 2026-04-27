import { NextRequest } from 'next/server';

export async function GET(_request: NextRequest, { params }: { params: Promise<{ key: string }> }) {
  const expected = process.env.INDEXNOW_KEY?.trim();
  if (!expected) {
    return new Response('Not found', { status: 404 });
  }
  const { key } = await params;
  if (key !== expected) {
    return new Response('Not found', { status: 404 });
  }
  return new Response(expected, {
    status: 200,
    headers: { 'Content-Type': 'text/plain', 'Cache-Control': 'public, max-age=86400' },
  });
}
