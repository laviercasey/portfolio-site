import { ImageResponse } from 'next/og';
import { siteConfig as siteJson } from '@/shared/config';
import type { SiteConfig } from '@/shared/types';

const site = siteJson as SiteConfig;

export const alt = 'Casey Laviere — Full-Stack Developer & ML Researcher';
export const size = { width: 1200, height: 630 };
export const contentType = 'image/png';
export const revalidate = 86400;

export default async function Image({ params }: { params: Promise<{ locale: string }> }) {
  const { locale } = await params;
  const isRu = locale === 'ru';

  const title = isRu ? (site.nameRu ?? site.name) : site.name;
  const subtitle = isRu ? site.taglineRu : site.taglineEn;
  const skills = ['PHP', 'Laravel', 'Python', 'TypeScript', 'React', 'Next.js', 'PostgreSQL', 'Docker'];

  return new ImageResponse(
    (
      <div
        style={{
          width: '100%',
          height: '100%',
          display: 'flex',
          flexDirection: 'column',
          justifyContent: 'center',
          padding: '80px',
          background: 'linear-gradient(135deg, #0f0f0f 0%, #1a1a2e 50%, #16213e 100%)',
          color: '#ffffff',
        }}
      >
        <div
          style={{
            display: 'flex',
            flexDirection: 'column',
            gap: '16px',
          }}
        >
          <div
            style={{
              fontSize: 72,
              fontWeight: 700,
              letterSpacing: '-2px',
              lineHeight: 1.1,
              background: 'linear-gradient(90deg, #a78bfa, #818cf8, #6366f1)',
              backgroundClip: 'text',
              color: 'transparent',
            }}
          >
            {title}
          </div>

          <div
            style={{
              fontSize: 32,
              fontWeight: 400,
              color: 'rgba(255, 255, 255, 0.7)',
              marginTop: '8px',
            }}
          >
            {subtitle}
          </div>

          <div
            style={{
              display: 'flex',
              flexWrap: 'wrap',
              gap: '12px',
              marginTop: '32px',
            }}
          >
            {skills.map((skill) => (
              <div
                key={skill}
                style={{
                  padding: '8px 20px',
                  borderRadius: '9999px',
                  border: '1px solid rgba(167, 139, 250, 0.3)',
                  backgroundColor: 'rgba(167, 139, 250, 0.1)',
                  fontSize: 18,
                  color: 'rgba(255, 255, 255, 0.8)',
                }}
              >
                {skill}
              </div>
            ))}
          </div>
        </div>

        <div
          style={{
            position: 'absolute',
            bottom: '40px',
            right: '80px',
            fontSize: 20,
            color: 'rgba(255, 255, 255, 0.35)',
          }}
        >
          {site.url.replace('https://', '')}
        </div>
      </div>
    ),
    { ...size },
  );
}
