'use client';

import { useEffect, useState } from 'react';
import { Button } from '@/shared/ui';
import { Card, CardContent } from '@/shared/ui';
import { Badge } from '@/shared/ui';
import { Input } from '@/shared/ui';
import { Textarea } from '@/shared/ui';
import { Label } from '@/shared/ui';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/shared/ui';
import type { Project } from '@/entities/project';
import { Plus, Pencil, Trash2, X, Loader2, Check } from 'lucide-react';

const emptyProject = {
  slug: '', category: 'web' as const, status: 'in_development' as const,
  tags: [] as string[], techStack: [] as string[], featured: false, order: 99, stars: 0,
  title: { en: '', ru: '' }, shortDescription: { en: '', ru: '' }, description: { en: '', ru: '' },
  goalDescription: { en: '', ru: '' },
  githubUrl: '', demoUrl: '', siteUrl: '', videoUrl: '', thumbnailUrl: '',
};

export default function ProjectsEditor() {
  const [projects, setProjects] = useState<Project[]>([]);
  const [loading, setLoading] = useState(true);
  const [editing, setEditing] = useState<Partial<Project> | null>(null);
  const [saving, setSaving] = useState(false);
  const [tagInput, setTagInput] = useState('');
  const [techInput, setTechInput] = useState('');

  useEffect(() => {
    fetch('/api/projects').then((r) => r.json()).then((d) => { setProjects(Array.isArray(d) ? d : []); setLoading(false); });
  }, []);

  const save = async () => {
    if (!editing) return;
    setSaving(true);
    const isNew = !editing.id;
    const url = isNew ? '/api/projects' : `/api/projects/${editing.id}`;
    const method = isNew ? 'POST' : 'PUT';
    const res = await fetch(url, { method, headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(editing) });
    if (res.ok) {
      const saved = await res.json();
      setProjects((prev) => isNew ? [...prev, saved] : prev.map((p) => p.id === saved.id ? saved : p));
      setEditing(null);
    }
    setSaving(false);
  };

  const deleteProject = async (id: string) => {
    if (!confirm('Delete this project?')) return;
    const res = await fetch(`/api/projects/${id}`, { method: 'DELETE' });
    if (res.ok) setProjects((prev) => prev.filter((p) => p.id !== id));
  };

  const addTag = (key: string) => {
    if (key && !editing?.tags?.includes(key)) {
      setEditing((prev) => ({ ...prev, tags: [...(prev?.tags ?? []), key] }));
    }
    setTagInput('');
  };

  const addTech = (key: string) => {
    if (key && !editing?.techStack?.includes(key)) {
      setEditing((prev) => ({ ...prev, techStack: [...(prev?.techStack ?? []), key] }));
    }
    setTechInput('');
  };

  if (loading) return <div className="flex justify-center py-20"><Loader2 className="h-6 w-6 animate-spin" /></div>;

  return (
      <div className="p-8">
        <div className="flex items-center justify-between mb-6">
          <h1 className="text-3xl font-bold">Projects</h1>
          <Button onClick={() => setEditing(emptyProject)}>
            <Plus className="mr-2 h-4 w-4" /> Add project
          </Button>
        </div>

        {editing && (
          <Card className="mb-6 border-primary/30">
            <CardContent className="p-6">
              <div className="flex items-center justify-between mb-4">
                <h2 className="text-lg font-semibold">{editing.id ? 'Edit project' : 'New project'}</h2>
                <Button variant="ghost" size="sm" onClick={() => setEditing(null)}><X className="h-4 w-4" /></Button>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label>Title (RU)</Label>
                  <Input value={editing.title?.ru ?? ''} onChange={(e) => setEditing((p) => ({ ...p, title: { ...(p?.title ?? { en: '', ru: '' }), ru: e.target.value } }))} />
                </div>
                <div className="space-y-2">
                  <Label>Title (EN)</Label>
                  <Input value={editing.title?.en ?? ''} onChange={(e) => setEditing((p) => ({ ...p, title: { ...(p?.title ?? { en: '', ru: '' }), en: e.target.value } }))} />
                </div>
                <div className="space-y-2">
                  <Label>Slug</Label>
                  <Input value={editing.slug ?? ''} onChange={(e) => setEditing((p) => ({ ...p, slug: e.target.value }))} placeholder="my-project" />
                </div>
                <div className="space-y-2">
                  <Label>Category</Label>
                  <Select value={editing.category ?? 'web'} onValueChange={(v) => setEditing((p) => ({ ...p, category: v as Project['category'] }))}>
                    <SelectTrigger><SelectValue /></SelectTrigger>
                    <SelectContent>
                      {['web', 'mobile', 'data', 'research', 'other'].map((c) => <SelectItem key={c} value={c}>{c}</SelectItem>)}
                    </SelectContent>
                  </Select>
                </div>
                <div className="space-y-2">
                  <Label>Status</Label>
                  <Select value={editing.status ?? 'in_development'} onValueChange={(v) => setEditing((p) => ({ ...p, status: v as Project['status'] }))}>
                    <SelectTrigger><SelectValue /></SelectTrigger>
                    <SelectContent>
                      <SelectItem value="completed">Completed</SelectItem>
                      <SelectItem value="in_development">In Development</SelectItem>
                      <SelectItem value="open_source">Open Source</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                <div className="space-y-2">
                  <Label>GitHub URL</Label>
                  <Input value={editing.githubUrl ?? ''} onChange={(e) => setEditing((p) => ({ ...p, githubUrl: e.target.value }))} placeholder="https://github.com/..." />
                </div>
                <div className="space-y-2">
                  <Label>Stars (manual, 0 = auto from GitHub API)</Label>
                  <Input type="number" value={editing.stars ?? 0} onChange={(e) => setEditing((p) => ({ ...p, stars: Number(e.target.value) }))} />
                </div>
                <div className="space-y-2">
                  <Label>Order</Label>
                  <Input type="number" value={editing.order ?? 99} onChange={(e) => setEditing((p) => ({ ...p, order: Number(e.target.value) }))} />
                </div>
                <div className="flex items-center gap-2 pt-6">
                  <input type="checkbox" checked={editing.featured ?? false} onChange={(e) => setEditing((p) => ({ ...p, featured: e.target.checked }))} className="h-4 w-4" />
                  <Label>Featured (show on homepage)</Label>
                </div>
                <div className="space-y-2 col-span-full">
                  <Label>Short Description (RU)</Label>
                  <Textarea value={editing.shortDescription?.ru ?? ''} rows={2} onChange={(e) => setEditing((p) => ({ ...p, shortDescription: { ...(p?.shortDescription ?? { en: '', ru: '' }), ru: e.target.value } }))} />
                </div>
                <div className="space-y-2 col-span-full">
                  <Label>Short Description (EN)</Label>
                  <Textarea value={editing.shortDescription?.en ?? ''} rows={2} onChange={(e) => setEditing((p) => ({ ...p, shortDescription: { ...(p?.shortDescription ?? { en: '', ru: '' }), en: e.target.value } }))} />
                </div>
                <div className="space-y-2 col-span-full">
                  <Label>Description (RU)</Label>
                  <Textarea value={editing.description?.ru ?? ''} rows={4} onChange={(e) => setEditing((p) => ({ ...p, description: { ...(p?.description ?? { en: '', ru: '' }), ru: e.target.value } }))} />
                </div>
                <div className="space-y-2 col-span-full">
                  <Label>Description (EN)</Label>
                  <Textarea value={editing.description?.en ?? ''} rows={4} onChange={(e) => setEditing((p) => ({ ...p, description: { ...(p?.description ?? { en: '', ru: '' }), en: e.target.value } }))} />
                </div>
                <div className="space-y-2 col-span-full">
                  <Label>Goal / Problem (RU)</Label>
                  <Textarea value={editing.goalDescription?.ru ?? ''} rows={2} onChange={(e) => setEditing((p) => ({ ...p, goalDescription: { ...(p?.goalDescription ?? { en: '', ru: '' }), ru: e.target.value } }))} placeholder="Какую проблему решает проект..." />
                </div>
                <div className="space-y-2 col-span-full">
                  <Label>Goal / Problem (EN)</Label>
                  <Textarea value={editing.goalDescription?.en ?? ''} rows={2} onChange={(e) => setEditing((p) => ({ ...p, goalDescription: { ...(p?.goalDescription ?? { en: '', ru: '' }), en: e.target.value } }))} placeholder="What problem does the project solve..." />
                </div>
                <div className="space-y-2">
                  <Label>Site URL</Label>
                  <Input value={editing.siteUrl ?? ''} onChange={(e) => setEditing((p) => ({ ...p, siteUrl: e.target.value }))} placeholder="https://..." />
                </div>

                <div className="space-y-2 col-span-full">
                  <Label>Tags</Label>
                  <div className="flex gap-2">
                    <Input value={tagInput} onChange={(e) => setTagInput(e.target.value)} placeholder="Add tag..." onKeyDown={(e) => e.key === 'Enter' && (e.preventDefault(), addTag(tagInput))} />
                    <Button type="button" size="sm" onClick={() => addTag(tagInput)}>Add</Button>
                  </div>
                  <div className="flex flex-wrap gap-1.5">
                    {editing.tags?.map((tag) => (
                      <Badge key={tag} variant="secondary" className="cursor-pointer" onClick={() => setEditing((p) => ({ ...p, tags: p?.tags?.filter((t) => t !== tag) }))}>
                        {tag} ×
                      </Badge>
                    ))}
                  </div>
                </div>

                <div className="space-y-2 col-span-full">
                  <Label>Tech Stack</Label>
                  <div className="flex gap-2">
                    <Input value={techInput} onChange={(e) => setTechInput(e.target.value)} placeholder="Add tech..." onKeyDown={(e) => e.key === 'Enter' && (e.preventDefault(), addTech(techInput))} />
                    <Button type="button" size="sm" onClick={() => addTech(techInput)}>Add</Button>
                  </div>
                  <div className="flex flex-wrap gap-1.5">
                    {editing.techStack?.map((tech) => (
                      <Badge key={tech} variant="outline" className="cursor-pointer" onClick={() => setEditing((p) => ({ ...p, techStack: p?.techStack?.filter((t) => t !== tech) }))}>
                        {tech} ×
                      </Badge>
                    ))}
                  </div>
                </div>

                <div className="space-y-2">
                  <Label>Thumbnail URL</Label>
                  <Input value={editing.thumbnailUrl ?? ''} onChange={(e) => setEditing((p) => ({ ...p, thumbnailUrl: e.target.value }))} />
                </div>
                <div className="space-y-2">
                  <Label>Video URL</Label>
                  <Input value={editing.videoUrl ?? ''} onChange={(e) => setEditing((p) => ({ ...p, videoUrl: e.target.value }))} />
                </div>
              </div>

              <div className="flex gap-3 mt-6">
                <Button onClick={save} disabled={saving}>
                  {saving ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : <Check className="mr-2 h-4 w-4" />}
                  Save
                </Button>
                <Button variant="outline" onClick={() => setEditing(null)}>Cancel</Button>
              </div>
            </CardContent>
          </Card>
        )}

        <div className="space-y-3">
          {projects.map((project) => (
            <Card key={project.id}>
              <CardContent className="p-4 flex items-center justify-between gap-4">
                <div>
                  <p className="font-semibold">{project.title.ru}</p>
                  <p className="text-sm text-muted-foreground">{project.slug} · {project.category}</p>
                </div>
                <div className="flex items-center gap-2">
                  <Badge variant={project.status === 'completed' ? 'success' : project.status === 'open_source' ? 'info' : 'warning'} className="text-xs">
                    {project.status}
                  </Badge>
                  <Button variant="ghost" size="sm" onClick={() => setEditing(project)}>
                    <Pencil className="h-4 w-4" />
                  </Button>
                  <Button variant="ghost" size="sm" className="text-destructive" onClick={() => deleteProject(project.id)}>
                    <Trash2 className="h-4 w-4" />
                  </Button>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
  );
}
