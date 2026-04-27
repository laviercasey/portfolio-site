'use client';

import { useEffect, useState } from 'react';
import { Badge } from '@/shared/ui';
import { Button } from '@/shared/ui';
import { Card, CardContent } from '@/shared/ui';
import type { Inquiry, InquiryStatus } from '@/entities/inquiry';
import { Loader2, ChevronDown, ChevronUp } from 'lucide-react';

const statusVariants: Record<InquiryStatus, 'warning' | 'info' | 'success' | 'outline'> = {
  new: 'warning', read: 'info', replied: 'success', archived: 'outline',
};

export default function InquiriesPage() {
  const [inquiries, setInquiries] = useState<Inquiry[]>([]);
  const [loading, setLoading] = useState(true);
  const [expanded, setExpanded] = useState<string | null>(null);
  const [filter, setFilter] = useState<InquiryStatus | 'all'>('all');

  useEffect(() => {
    fetch('/api/inquiries')
      .then((r) => r.json())
      .then((data) => { setInquiries(Array.isArray(data) ? data : []); setLoading(false); })
      .catch(() => setLoading(false));
  }, []);

  const updateStatus = async (id: string, status: InquiryStatus) => {
    const res = await fetch(`/api/inquiries/${id}`, {
      method: 'PATCH',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ status }),
    });
    if (res.ok) {
      setInquiries((prev) => prev.map((i) => i.id === id ? { ...i, status } : i));
    }
  };

  const filtered = filter === 'all' ? inquiries : inquiries.filter((i) => i.status === filter);

  return (
      <div className="p-8">
        <h1 className="text-3xl font-bold mb-6">Inquiries</h1>

        <div className="flex flex-wrap gap-2 mb-6">
          {(['all', 'new', 'read', 'replied', 'archived'] as const).map((s) => (
            <Button key={s} variant={filter === s ? 'default' : 'outline'} size="sm" onClick={() => setFilter(s)}>
              {s === 'all' ? 'All' : s.charAt(0).toUpperCase() + s.slice(1)}
              {s !== 'all' && (
                <span className="ml-1.5 text-xs opacity-70">
                  ({inquiries.filter((i) => i.status === s).length})
                </span>
              )}
            </Button>
          ))}
        </div>

        {loading ? (
          <div className="flex justify-center py-12"><Loader2 className="h-6 w-6 animate-spin" /></div>
        ) : filtered.length === 0 ? (
          <p className="text-center py-12 text-muted-foreground">No inquiries found</p>
        ) : (
          <div className="space-y-3">
            {filtered.map((inq) => (
              <Card key={inq.id}>
                <CardContent className="p-5">
                  <div className="flex flex-wrap items-start justify-between gap-3">
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2 mb-1">
                        <p className="font-semibold">{inq.name}</p>
                        <Badge variant={statusVariants[inq.status]} className="text-xs">{inq.status}</Badge>
                        <Badge variant="outline" className="text-xs">{inq.type}</Badge>
                      </div>
                      <p className="text-sm text-muted-foreground">{inq.email}{inq.company ? ` · ${inq.company}` : ''}</p>
                      {inq.telegram && <p className="text-xs text-muted-foreground mt-0.5">Telegram: {inq.telegram}</p>}
                      {inq.budget && <p className="text-xs text-muted-foreground mt-0.5">Budget: {inq.budget}</p>}
                    </div>
                    <div className="flex items-center gap-2">
                      <span className="text-xs text-muted-foreground">{new Date(inq.createdAt).toLocaleDateString()}</span>
                      <Button variant="ghost" size="sm" onClick={() => setExpanded(expanded === inq.id ? null : inq.id)}>
                        {expanded === inq.id ? <ChevronUp className="h-4 w-4" /> : <ChevronDown className="h-4 w-4" />}
                      </Button>
                    </div>
                  </div>

                  {expanded === inq.id && (
                    <div className="mt-4 pt-4 border-t space-y-4">
                      <div className="bg-muted/50 rounded-lg p-4">
                        <p className="text-sm whitespace-pre-wrap">{inq.message}</p>
                      </div>

                      <div>
                        <p className="text-xs font-medium text-muted-foreground mb-2">Update status:</p>
                        <div className="flex flex-wrap gap-2">
                          {(['new', 'read', 'replied', 'archived'] as InquiryStatus[]).map((s) => (
                            <Button
                              key={s}
                              size="sm"
                              variant={inq.status === s ? 'default' : 'outline'}
                              onClick={() => updateStatus(inq.id, s)}
                            >
                              {s}
                            </Button>
                          ))}
                        </div>
                      </div>
                    </div>
                  )}
                </CardContent>
              </Card>
            ))}
          </div>
        )}
      </div>
  );
}
