import { ImageResponse } from 'next/og';
import { projectsService } from '@/entities/project';
import { siteConfig as siteJson } from '@/shared/config';
import type { SiteConfig } from '@/shared/types';

export const alt = 'Project case study';
export const size = { width: 1200, height: 630 };
export const contentType = 'image/png';
export const revalidate = 86400;

const site = siteJson as SiteConfig;

export default async function Image({
  params,
}: {
  params: Promise<{ locale: string; slug: string }>;
}) {
  const { locale, slug } = await params;
  const isRu = locale === 'ru';

  let project;
  try {
    project = await projectsService.getBySlug(slug);
  } catch {
    project = null;
  }

  const title = project
    ? (isRu ? project.title.ru : project.title.en)
    : site.name;
  const subtitle = project
    ? (isRu ? project.shortDescription.ru : project.shortDescription.en)
    : (isRu ? site.taglineRu : site.taglineEn);
  const tags = project ? project.techStack.slice(0, 6) : [];

  return new ImageResponse(
    (
      <div
        style={{
          width: '100%',
          height: '100%',
          display: 'flex',
          flexDirection: 'column',
          justifyContent: 'space-between',
          padding: '72px',
          background: 'linear-gradient(135deg, #0f0f0f 0%, #1a1a2e 50%, #16213e 100%)',
          color: '#ffffff',
        }}
      >
        <div style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
          <div
            style={{
              fontSize: 24,
              letterSpacing: '2px',
              textTransform: 'uppercase',
              color: 'rgba(167, 139, 250, 0.9)',
            }}
          >
            {isRu ? 'Проект' : 'Case Study'}
          </div>
          <div
            style={{
              fontSize: 72,
              fontWeight: 700,
              lineHeight: 1.05,
              letterSpacing: '-2px',
              background: 'linear-gradient(90deg, #a78bfa, #818cf8, #6366f1)',
              backgroundClip: 'text',
              color: 'transparent',
            }}
          >
            {title}
          </div>
          <div
            style={{
              fontSize: 28,
              lineHeight: 1.3,
              color: 'rgba(255, 255, 255, 0.75)',
              maxWidth: '920px',
            }}
          >
            {subtitle}
          </div>
          {tags.length > 0 && (
            <div style={{ display: 'flex', flexWrap: 'wrap', gap: '12px', marginTop: '12px' }}>
              {tags.map((t) => (
                <div
                  key={t}
                  style={{
                    padding: '6px 18px',
                    borderRadius: '9999px',
                    border: '1px solid rgba(167, 139, 250, 0.3)',
                    backgroundColor: 'rgba(167, 139, 250, 0.12)',
                    fontSize: 18,
                    color: 'rgba(255, 255, 255, 0.85)',
                  }}
                >
                  {t}
                </div>
              ))}
            </div>
          )}
        </div>

        <div
          style={{
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'flex-end',
            fontSize: 22,
            color: 'rgba(255, 255, 255, 0.55)',
          }}
        >
          <span>{site.name}</span>
          <span>{site.url.replace('https://', '')}</span>
        </div>
      </div>
    ),
    { ...size },
  );
}
