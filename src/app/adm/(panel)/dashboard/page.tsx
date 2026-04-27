import { Suspense } from 'react';
import { redirect } from 'next/navigation';
import { getAuthToken } from '@/shared/lib/server';
import { getAdminSlug } from '@/shared/lib';
import { projectsService } from '@/entities/project';
import { inquiriesService } from '@/entities/inquiry';
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui';
import { Badge } from '@/shared/ui';
import { FolderOpen, Inbox, Clock, CheckCircle2 } from 'lucide-react';
import Link from 'next/link';
import { AdminAnalytics, AnalyticsSkeleton } from '@/widgets/admin-analytics';
import type { AnalyticsRange } from '@/entities/analytics';

type DashboardSearchParams = { range?: string };

interface DashboardPageProps {
  searchParams?: Promise<DashboardSearchParams>;
}

function parseRange(value: string | undefined): AnalyticsRange | undefined {
  return value === '7d' || value === '30d' ? value : undefined;
}

export default async function DashboardPage({ searchParams }: DashboardPageProps) {
  const token = await getAuthToken();
  const adminSlug = getAdminSlug();
  if (!token) redirect(`/${adminSlug}`);

  const resolvedParams = (await searchParams) ?? {};
  const range = parseRange(resolvedParams.range);

  const [projects, inquiries] = await Promise.all([projectsService.list(), inquiriesService.list(token)]);
  const newInquiries = inquiries.filter((i) => i.status === 'new');
  const recentInquiries = inquiries.slice(0, 5);

  const dashboardPath = `/${adminSlug}/dashboard`;

  return (
      <div className="p-8">
        <h1 className="text-3xl font-bold mb-8">Dashboard</h1>

        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
          <Card>
            <CardContent className="p-5">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-muted-foreground">Projects</p>
                  <p className="text-3xl font-bold">{projects.length}</p>
                </div>
                <FolderOpen className="h-8 w-8 text-primary/30" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-5">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-muted-foreground">Total Inquiries</p>
                  <p className="text-3xl font-bold">{inquiries.length}</p>
                </div>
                <Inbox className="h-8 w-8 text-primary/30" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-5">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-muted-foreground">New Inquiries</p>
                  <p className="text-3xl font-bold text-orange-500">{newInquiries.length}</p>
                </div>
                <Clock className="h-8 w-8 text-orange-300" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-5">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-muted-foreground">Completed</p>
                  <p className="text-3xl font-bold text-green-500">{projects.filter(p => p.status === 'completed').length}</p>
                </div>
                <CheckCircle2 className="h-8 w-8 text-green-300" />
              </div>
            </CardContent>
          </Card>
        </div>

        <div className="mb-8">
          <Suspense fallback={<AnalyticsSkeleton />}>
            <AdminAnalytics token={token} pathname={dashboardPath} range={range} />
          </Suspense>
        </div>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between">
            <CardTitle>Recent Inquiries</CardTitle>
            <Link href={`/${adminSlug}/inquiries`} className="text-sm text-primary hover:underline">
              View all
            </Link>
          </CardHeader>
          <CardContent>
            {recentInquiries.length === 0 ? (
              <p className="text-muted-foreground text-sm py-4 text-center">No inquiries yet</p>
            ) : (
              <div className="space-y-3">
                {recentInquiries.map((inq) => (
                  <div key={inq.id} className="flex items-center justify-between p-3 rounded-lg border">
                    <div>
                      <p className="font-medium text-sm">{inq.name}</p>
                      <p className="text-xs text-muted-foreground">{inq.email} · {inq.type}</p>
                    </div>
                    <div className="flex items-center gap-2">
                      <Badge variant={inq.status === 'new' ? 'warning' : inq.status === 'replied' ? 'success' : 'outline'} className="text-xs">
                        {inq.status}
                      </Badge>
                      <span className="text-xs text-muted-foreground">
                        {new Date(inq.createdAt).toLocaleDateString()}
                      </span>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>
      </div>
  );
}
