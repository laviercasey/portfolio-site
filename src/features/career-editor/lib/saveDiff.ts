import type {
  CareerContent,
  Education,
  WorkExperience,
  Certificate,
  Publication,
} from '@/entities/career';

const UUID_REGEX = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;

type CareerType = 'education' | 'work' | 'certificate' | 'publication';

async function callApi(method: string, path: string, body?: unknown): Promise<void> {
  const res = await fetch(path, {
    method,
    headers: body ? { 'Content-Type': 'application/json' } : undefined,
    body: body ? JSON.stringify(body) : undefined,
  });
  if (!res.ok && res.status !== 204) {
    const err = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error(`${method} ${path}: ${err.error ?? res.statusText}`);
  }
}

function toInt(value: string): number | null {
  if (!value) return null;
  const n = parseInt(value, 10);
  return Number.isFinite(n) ? n : null;
}

function educationPayload(item: Education, sortOrder: number) {
  const startYear = toInt(item.startYear) ?? 0;
  const endYear = toInt(item.endYear);
  return {
    institution: item.institution,
    degree: item.degree,
    field: item.field,
    startYear,
    endYear,
    description: item.description ?? { en: '', ru: '' },
    logoUrl: item.logoUrl ?? '',
    sortOrder,
    relatedProjectSlugs: item.relatedProjectSlugs ?? [],
  };
}

function workPayload(item: WorkExperience, sortOrder: number) {
  return {
    company: item.company,
    position: item.position,
    startDate: item.startDate,
    endDate: item.current ? null : (item.endDate ?? null),
    current: item.current,
    description: item.description,
    technologies: item.technologies ?? [],
    logoUrl: item.logoUrl ?? '',
    sortOrder,
    achievements: item.achievements ?? [],
    fullDescription: item.fullDescription ?? { en: '', ru: '' },
  };
}

function certificatePayload(item: Certificate, sortOrder: number) {
  return {
    title: item.title,
    issuer: item.issuer,
    date: item.date,
    url: item.url ?? '',
    credentialId: item.credentialId ?? '',
    sortOrder,
  };
}

function publicationPayload(item: Publication, sortOrder: number) {
  return {
    title: item.title,
    journal: item.journal ?? { en: '', ru: '' },
    year: toInt(item.year) ?? 0,
    doi: item.doi ?? '',
    url: item.url ?? '',
    abstract: item.abstract ?? { en: '', ru: '' },
    sortOrder,
  };
}

async function diffSection<T extends { id: string }>(
  type: CareerType,
  original: T[],
  current: T[],
  buildPayload: (item: T, sortOrder: number) => unknown,
): Promise<void> {
  const originalIds = new Set(original.map((i) => i.id));
  const currentIds = new Set(current.map((i) => i.id));

  for (const id of originalIds) {
    if (!currentIds.has(id)) {
      await callApi('DELETE', `/api/career/${type}/${id}`);
    }
  }

  for (let i = 0; i < current.length; i++) {
    const item = current[i];
    const payload = buildPayload(item, i);
    if (!UUID_REGEX.test(item.id)) {
      await callApi('POST', `/api/career/${type}`, payload);
    } else {
      const orig = original.find((o) => o.id === item.id);
      if (!orig || JSON.stringify(orig) !== JSON.stringify(item)) {
        await callApi('PUT', `/api/career/${type}/${item.id}`, payload);
      }
    }
  }
}

export async function saveCareerDiff(
  original: CareerContent,
  current: CareerContent,
): Promise<void> {
  await diffSection('education', original.education, current.education, educationPayload);
  await diffSection('work', original.workHistory, current.workHistory, workPayload);
  await diffSection('certificate', original.certificates, current.certificates, certificatePayload);
  await diffSection('publication', original.publications, current.publications, publicationPayload);
}
