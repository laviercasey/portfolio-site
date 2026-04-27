import { redirect } from 'next/navigation';
import { getAuthToken } from '@/shared/lib/server';
import { getAdminSlug } from '@/shared/lib';
import { AdminLayoutClient } from '@/widgets/admin-layout';

export default async function AdminPanelLayout({ children }: { children: React.ReactNode }) {
  const token = await getAuthToken();
  const slug = getAdminSlug();

  if (!token) {
    redirect(`/${slug}`);
  }

  return <AdminLayoutClient adminSlug={slug}>{children}</AdminLayoutClient>;
}
