import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  basePath: '/portfolio-site',
  assetPrefix: '/portfolio-site',
  output: 'export',
  trailingSlash: true,
  images: {
    unoptimized: true,
  },
};

export default nextConfig;
