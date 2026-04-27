import { AdminLogin } from '@/features/admin-login';
import { getAdminSlug } from '@/shared/lib';

export const dynamic = 'force-dynamic';

export default function AdminLoginPage() {
  return <AdminLogin adminSlug={getAdminSlug()} />;
}
