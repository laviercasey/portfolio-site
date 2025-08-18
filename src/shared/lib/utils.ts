export function getAssetPath(path: string): string {
  // В режиме static export с basePath нужно добавлять префикс вручную
  const basePath = '/portfolio-site';
  return `${basePath}${path}`;
}