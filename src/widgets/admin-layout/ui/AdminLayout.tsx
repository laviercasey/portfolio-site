'use client';

import Link from 'next/link';
import { usePathname, useRouter } from 'next/navigation';
import {
  LayoutDashboard, FolderOpen, FileText, Inbox, Image as ImageIcon,
  LogOut, ChevronRight, Code2, MessageSquare, GraduationCap, Wrench,
} from 'lucide-react';
import { Button } from '@/shared/ui';
import { cn } from '@/shared/lib';

const navItems = [
  { label: 'Dashboard', icon: LayoutDashboard, path: 'dashboard' },
  { label: 'Projects', icon: FolderOpen, path: 'projects' },
  { label: 'Services', icon: Wrench, path: 'services' },
  { label: 'Content', icon: FileText, path: 'content' },
  { label: 'Career', icon: GraduationCap, path: 'career' },
  { label: 'Contact', icon: MessageSquare, path: 'contact' },
  { label: 'Inquiries', icon: Inbox, path: 'inquiries' },
  { label: 'Media', icon: ImageIcon, path: 'media' },
];

interface AdminLayoutProps {
  children: React.ReactNode;
  adminSlug: string;
}

export default function AdminLayoutClient({ children, adminSlug }: AdminLayoutProps) {
  const pathname = usePathname();
  const router = useRouter();

  const adminBase = `/${adminSlug}`;

  const parts = pathname.split('/').filter(Boolean);
  const subPath = parts.slice(1).join('/');

  const handleLogout = async () => {
    await fetch('/api/auth/logout', { method: 'POST' });
    router.push(adminBase);
    router.refresh();
  };

  return (
    <div className="min-h-screen flex bg-muted/20">
      <aside className="w-60 border-r bg-card flex flex-col">
        <div className="p-4 border-b">
          <div className="flex items-center gap-2 font-bold text-primary">
            <Code2 className="h-5 w-5" />
            <span>Admin Panel</span>
          </div>
          <p className="text-xs text-muted-foreground mt-0.5">Casey Laviere Portfolio</p>
        </div>

        <nav className="flex-1 p-3 space-y-1">
          {navItems.map(({ label, icon: Icon, path }) => {
            const href = `${adminBase}/${path}`;
            const active = subPath === path || subPath.startsWith(`${path}/`);
            return (
              <Link
                key={path}
                href={href}
                className={cn(
                  'flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors',
                  active
                    ? 'bg-primary text-primary-foreground'
                    : 'text-muted-foreground hover:bg-muted hover:text-foreground'
                )}
              >
                <Icon className="h-4 w-4" />
                {label}
                {active && <ChevronRight className="ml-auto h-3 w-3" />}
              </Link>
            );
          })}
        </nav>

        <div className="p-3 border-t">
          <Button variant="ghost" size="sm" className="w-full justify-start text-muted-foreground" onClick={handleLogout}>
            <LogOut className="mr-2 h-4 w-4" />
            Log out
          </Button>
        </div>
      </aside>

      <main className="flex-1 overflow-auto">
        {children}
      </main>
    </div>
  );
}
