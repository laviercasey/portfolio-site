const createNextIntlPlugin = require('next-intl/plugin');
const withNextIntl = createNextIntlPlugin('./src/shared/config/i18n/request.ts');

const nextConfig = {
  output: 'standalone',
  trailingSlash: false,
  experimental: {
    inlineCss: true,
  },
  images: {
    formats: ['image/avif', 'image/webp'],
    deviceSizes: [640, 750, 828, 1080, 1200, 1920],
    imageSizes: [16, 32, 48, 64, 96, 128, 256, 384],
    minimumCacheTTL: 2592000,
    remotePatterns: [
      { protocol: 'https', hostname: '**.vercel-storage.com' },
      { protocol: 'https', hostname: 'images.unsplash.com' },
    ],
  },
  async headers() {
    const umamiSrc = process.env.NEXT_PUBLIC_UMAMI_SRC;
    let umamiOrigin = '';
    if (umamiSrc) {
      try {
        umamiOrigin = new URL(umamiSrc).origin;
      } catch {}
    }

    const csp = [
      "default-src 'self'",
      `script-src 'self' 'unsafe-inline'${umamiOrigin ? ' ' + umamiOrigin : ''}`,
      "style-src 'self' 'unsafe-inline' https://fonts.googleapis.com",
      "img-src 'self' data: blob: https:",
      "font-src 'self' https://fonts.gstatic.com",
      `connect-src 'self'${umamiOrigin ? ' ' + umamiOrigin : ''}`,
      "object-src 'none'",
      "base-uri 'self'",
      "frame-ancestors 'none'",
      "form-action 'self'",
    ].join('; ');

    return [
      {
        source: '/(.*)',
        headers: [
          { key: 'X-Frame-Options', value: 'DENY' },
          { key: 'X-Content-Type-Options', value: 'nosniff' },
          { key: 'Referrer-Policy', value: 'strict-origin-when-cross-origin' },
          { key: 'Permissions-Policy', value: 'camera=(), microphone=(), geolocation=()' },
          { key: 'Strict-Transport-Security', value: 'max-age=63072000; includeSubDomains; preload' },
          { key: 'Content-Security-Policy', value: csp },
        ],
      },
    ];
  },
  async rewrites() {
    return [
      {
        source: '/:key([a-f0-9]{8,64}).txt',
        destination: '/api/indexnow-verification/:key',
      },
    ];
  },
};

module.exports = withNextIntl(nextConfig);
