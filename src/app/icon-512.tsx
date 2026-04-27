import { ImageResponse } from 'next/og';

export const size = { width: 512, height: 512 };
export const contentType = 'image/png';

export default function Icon512() {
  return new ImageResponse(
    (
      <div
        style={{
          width: '100%',
          height: '100%',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          background: 'radial-gradient(120% 120% at 30% 20%, #3a2418 0%, #1a120c 55%, #0a0706 100%)',
          borderRadius: 112,
          color: 'transparent',
          backgroundClip: 'padding-box',
          fontSize: 240,
          fontWeight: 900,
          fontStyle: 'italic',
          fontFamily: 'Georgia, serif',
          letterSpacing: -14,
        }}
      >
        <span
          style={{
            backgroundImage: 'linear-gradient(135deg, #d4a843 0%, #e07040 55%, #c45a30 100%)',
            backgroundClip: 'text',
            WebkitBackgroundClip: 'text',
            color: 'transparent',
          }}
        >
          CL
        </span>
      </div>
    ),
    { ...size },
  );
}
