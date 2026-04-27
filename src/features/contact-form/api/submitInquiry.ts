export async function submitInquiry(payload: Record<string, unknown>): Promise<void> {
  const res = await fetch('/api/inquiries', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  });
  if (!res.ok) {
    throw new Error(`Inquiry submission failed with status ${res.status}`);
  }
}
