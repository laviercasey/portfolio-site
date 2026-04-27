export function newLocalId(): string {
  if (typeof crypto !== 'undefined' && 'randomUUID' in crypto) {
    return crypto.randomUUID();
  }
  return `k-${Date.now()}-${Math.random().toString(36).slice(2, 10)}`;
}

export function ensureIds<T extends { id?: string }>(items: T[]): T[] {
  let mutated = false;
  const next = items.map((item) => {
    if (item.id) return item;
    mutated = true;
    return { ...item, id: newLocalId() };
  });
  return mutated ? next : items;
}
