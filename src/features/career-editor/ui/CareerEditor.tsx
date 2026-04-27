'use client';

import { useEffect, useState } from 'react';
import { Button, Input, Textarea, Label, Card, CardContent, CardHeader, CardTitle } from '@/shared/ui';
import { Loader2, Save, Check, Plus, Trash2, ChevronUp, ChevronDown } from 'lucide-react';
import type { CareerContent, Education, WorkExperience, Certificate, Publication } from '@/entities/career';
import { toYearMonth } from '@/shared/lib';
import { saveCareerDiff } from '../lib/saveDiff';

function moveItem<T>(arr: T[], i: number, dir: -1 | 1): T[] {
  const copy = [...arr];
  const ni = i + dir;
  if (ni < 0 || ni >= copy.length) return copy;
  [copy[i], copy[ni]] = [copy[ni], copy[i]];
  return copy;
}

async function fetchCareer(): Promise<CareerContent | null> {
  const res = await fetch('/api/career');
  if (!res.ok) return null;
  const raw = (await res.json()) as Partial<CareerContent> | null;
  return {
    education: (raw?.education ?? []).map((e) => ({
      ...e,
      startYear: e.startYear?.toString() ?? '',
      endYear: e.endYear?.toString() ?? '',
    })),
    workHistory: (raw?.workHistory ?? []).map((w) => ({
      ...w,
      startDate: toYearMonth(w.startDate),
      endDate: w.endDate ? toYearMonth(w.endDate) : undefined,
    })),
    certificates: (raw?.certificates ?? []).map((c) => ({
      ...c,
      date: toYearMonth(c.date),
    })),
    publications: raw?.publications ?? [],
  };
}

export default function CareerEditor() {
  const [data, setData] = useState<CareerContent | null>(null);
  const [original, setOriginal] = useState<CareerContent | null>(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [saved, setSaved] = useState(false);
  const [saveError, setSaveError] = useState<string | null>(null);

  const loadCareer = async () => {
    try {
      const normalized = await fetchCareer();
      if (!normalized) {
        setLoading(false);
        return;
      }
      setData(normalized);
      setOriginal(JSON.parse(JSON.stringify(normalized)));
    } catch {
      setSaveError('Failed to load career data');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadCareer();
  }, []);

  const handleSave = async () => {
    if (!data || !original) return;
    setSaving(true);
    setSaveError(null);
    try {
      await saveCareerDiff(original, data);
      setSaved(true);
      setTimeout(() => setSaved(false), 2000);
      await loadCareer();
    } catch (e) {
      setSaveError(e instanceof Error ? e.message : 'Save failed');
    } finally {
      setSaving(false);
    }
  };

  if (loading || !data) {
    return <div className="flex justify-center py-20"><Loader2 className="h-6 w-6 animate-spin" /></div>;
  }

  const addEducation = () => setData((p) => p ? { ...p, education: [...p.education, {
    id: `edu-${Date.now()}`, institution: { en: '', ru: '' }, degree: { en: '', ru: '' }, field: { en: '', ru: '' }, startYear: '', endYear: '',
  }] } : p);
  const updateEdu = (i: number, updates: Partial<Education>) => setData((p) => {
    if (!p) return p;
    const arr = [...p.education]; arr[i] = { ...arr[i], ...updates }; return { ...p, education: arr };
  });

  const addWork = () => setData((p) => p ? { ...p, workHistory: [...p.workHistory, {
    id: `work-${Date.now()}`, company: { en: '', ru: '' }, position: { en: '', ru: '' }, startDate: '', current: false, description: { en: '', ru: '' }, technologies: [],
  }] } : p);
  const updateWork = (i: number, updates: Partial<WorkExperience>) => setData((p) => {
    if (!p) return p;
    const arr = [...p.workHistory]; arr[i] = { ...arr[i], ...updates }; return { ...p, workHistory: arr };
  });

  const addCert = () => setData((p) => p ? { ...p, certificates: [...p.certificates, {
    id: `cert-${Date.now()}`, title: { en: '', ru: '' }, issuer: { en: '', ru: '' }, date: '',
  }] } : p);
  const updateCert = (i: number, updates: Partial<Certificate>) => setData((p) => {
    if (!p) return p;
    const arr = [...p.certificates]; arr[i] = { ...arr[i], ...updates }; return { ...p, certificates: arr };
  });

  const addPub = () => setData((p) => p ? { ...p, publications: [...p.publications, {
    id: `pub-${Date.now()}`, title: { en: '', ru: '' }, year: new Date().getFullYear().toString(),
  }] } : p);
  const updatePub = (i: number, updates: Partial<Publication>) => setData((p) => {
    if (!p) return p;
    const arr = [...p.publications]; arr[i] = { ...arr[i], ...updates }; return { ...p, publications: arr };
  });

  return (
    <div className="p-8 max-w-4xl">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-3xl font-bold">Career & Education</h1>
        <div className="flex items-center gap-3">
          {saveError && <span className="text-sm text-destructive">{saveError}</span>}
          <Button onClick={handleSave} disabled={saving}>
            {saving ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : saved ? <Check className="mr-2 h-4 w-4" /> : <Save className="mr-2 h-4 w-4" />}
            {saved ? 'Saved!' : 'Save'}
          </Button>
        </div>
      </div>

      <div className="space-y-6">
        <Card>
          <CardHeader><div className="flex items-center justify-between"><CardTitle>Education</CardTitle><Button variant="outline" size="sm" onClick={addEducation}><Plus className="mr-1.5 h-3.5 w-3.5" />Add</Button></div></CardHeader>
          <CardContent className="space-y-4">
            {data.education.map((edu, i) => (
              <div key={edu.id} className="border rounded-lg p-4 space-y-3 bg-muted/20">
                <div className="flex items-center gap-2">
                  <div className="flex flex-col gap-0.5">
                    <button onClick={() => setData((p) => p ? { ...p, education: moveItem(p.education, i, -1) } : p)} disabled={i === 0} className="disabled:opacity-20 p-0.5"><ChevronUp className="h-3.5 w-3.5" /></button>
                    <button onClick={() => setData((p) => p ? { ...p, education: moveItem(p.education, i, 1) } : p)} disabled={i === data.education.length - 1} className="disabled:opacity-20 p-0.5"><ChevronDown className="h-3.5 w-3.5" /></button>
                  </div>
                  <span className="text-xs font-mono text-muted-foreground">{edu.id}</span>
                  <Button variant="ghost" size="sm" onClick={() => setData((p) => p ? { ...p, education: p.education.filter((_, j) => j !== i) } : p)} className="ml-auto text-destructive h-7 w-7 p-0"><Trash2 className="h-3.5 w-3.5" /></Button>
                </div>
                <div className="grid grid-cols-2 gap-3">
                  <div className="space-y-1"><Label className="text-xs">Institution (RU)</Label><Input value={edu.institution.ru} onChange={(e) => updateEdu(i, { institution: { ...edu.institution, ru: e.target.value } })} className="h-8 text-sm" /></div>
                  <div className="space-y-1"><Label className="text-xs">Institution (EN)</Label><Input value={edu.institution.en} onChange={(e) => updateEdu(i, { institution: { ...edu.institution, en: e.target.value } })} className="h-8 text-sm" /></div>
                </div>
                <div className="grid grid-cols-2 gap-3">
                  <div className="space-y-1"><Label className="text-xs">Degree (RU)</Label><Input value={edu.degree.ru} onChange={(e) => updateEdu(i, { degree: { ...edu.degree, ru: e.target.value } })} className="h-8 text-sm" /></div>
                  <div className="space-y-1"><Label className="text-xs">Degree (EN)</Label><Input value={edu.degree.en} onChange={(e) => updateEdu(i, { degree: { ...edu.degree, en: e.target.value } })} className="h-8 text-sm" /></div>
                </div>
                <div className="grid grid-cols-2 gap-3">
                  <div className="space-y-1"><Label className="text-xs">Field (RU)</Label><Input value={edu.field.ru} onChange={(e) => updateEdu(i, { field: { ...edu.field, ru: e.target.value } })} className="h-8 text-sm" /></div>
                  <div className="space-y-1"><Label className="text-xs">Field (EN)</Label><Input value={edu.field.en} onChange={(e) => updateEdu(i, { field: { ...edu.field, en: e.target.value } })} className="h-8 text-sm" /></div>
                </div>
                <div className="grid grid-cols-3 gap-3">
                  <div className="space-y-1"><Label className="text-xs">Start Year</Label><Input value={edu.startYear} onChange={(e) => updateEdu(i, { startYear: e.target.value })} className="h-8 text-sm" /></div>
                  <div className="space-y-1"><Label className="text-xs">End Year</Label><Input value={edu.endYear} onChange={(e) => updateEdu(i, { endYear: e.target.value })} className="h-8 text-sm" /></div>
                  <div className="space-y-1"><Label className="text-xs">Logo URL</Label><Input value={edu.logoUrl ?? ''} onChange={(e) => updateEdu(i, { logoUrl: e.target.value || undefined })} className="h-8 text-sm" /></div>
                </div>
                <div className="grid grid-cols-2 gap-3">
                  <div className="space-y-1"><Label className="text-xs">Description (RU)</Label><Textarea rows={2} value={edu.description?.ru ?? ''} onChange={(e) => updateEdu(i, { description: { en: edu.description?.en ?? '', ru: e.target.value } })} className="text-sm" /></div>
                  <div className="space-y-1"><Label className="text-xs">Description (EN)</Label><Textarea rows={2} value={edu.description?.en ?? ''} onChange={(e) => updateEdu(i, { description: { ru: edu.description?.ru ?? '', en: e.target.value } })} className="text-sm" /></div>
                </div>
                <div className="space-y-1">
                  <Label className="text-xs">Related project slugs (comma-separated — matches project.slug)</Label>
                  <Input
                    value={(edu.relatedProjectSlugs ?? []).join(', ')}
                    onChange={(e) => updateEdu(i, { relatedProjectSlugs: e.target.value.split(',').map(s => s.trim()).filter(Boolean) })}
                    className="h-8 text-sm font-mono"
                    placeholder="heart-disease-ml-benchmark, genreneuro"
                  />
                </div>
              </div>
            ))}
          </CardContent>
        </Card>

        <Card>
          <CardHeader><div className="flex items-center justify-between"><CardTitle>Work History</CardTitle><Button variant="outline" size="sm" onClick={addWork}><Plus className="mr-1.5 h-3.5 w-3.5" />Add</Button></div></CardHeader>
          <CardContent className="space-y-4">
            {data.workHistory.map((w, i) => (
              <div key={w.id} className="border rounded-lg p-4 space-y-3 bg-muted/20">
                <div className="flex items-center gap-2">
                  <div className="flex flex-col gap-0.5">
                    <button onClick={() => setData((p) => p ? { ...p, workHistory: moveItem(p.workHistory, i, -1) } : p)} disabled={i === 0} className="disabled:opacity-20 p-0.5"><ChevronUp className="h-3.5 w-3.5" /></button>
                    <button onClick={() => setData((p) => p ? { ...p, workHistory: moveItem(p.workHistory, i, 1) } : p)} disabled={i === data.workHistory.length - 1} className="disabled:opacity-20 p-0.5"><ChevronDown className="h-3.5 w-3.5" /></button>
                  </div>
                  <span className="text-xs font-mono text-muted-foreground">{w.id}</span>
                  <label className="ml-auto flex items-center gap-1.5 text-xs"><input type="checkbox" checked={w.current} onChange={(e) => updateWork(i, { current: e.target.checked, endDate: e.target.checked ? undefined : w.endDate })} className="h-3.5 w-3.5" />Current</label>
                  <Button variant="ghost" size="sm" onClick={() => setData((p) => p ? { ...p, workHistory: p.workHistory.filter((_, j) => j !== i) } : p)} className="text-destructive h-7 w-7 p-0"><Trash2 className="h-3.5 w-3.5" /></Button>
                </div>
                <div className="grid grid-cols-2 gap-3">
                  <div className="space-y-1"><Label className="text-xs">Company (RU)</Label><Input value={w.company.ru} onChange={(e) => updateWork(i, { company: { ...w.company, ru: e.target.value } })} className="h-8 text-sm" /></div>
                  <div className="space-y-1"><Label className="text-xs">Company (EN)</Label><Input value={w.company.en} onChange={(e) => updateWork(i, { company: { ...w.company, en: e.target.value } })} className="h-8 text-sm" /></div>
                </div>
                <div className="grid grid-cols-2 gap-3">
                  <div className="space-y-1"><Label className="text-xs">Position (RU)</Label><Input value={w.position.ru} onChange={(e) => updateWork(i, { position: { ...w.position, ru: e.target.value } })} className="h-8 text-sm" /></div>
                  <div className="space-y-1"><Label className="text-xs">Position (EN)</Label><Input value={w.position.en} onChange={(e) => updateWork(i, { position: { ...w.position, en: e.target.value } })} className="h-8 text-sm" /></div>
                </div>
                <div className="grid grid-cols-3 gap-3">
                  <div className="space-y-1"><Label className="text-xs">Start Date</Label><Input value={w.startDate} onChange={(e) => updateWork(i, { startDate: e.target.value })} className="h-8 text-sm" placeholder="2023-09" /></div>
                  <div className="space-y-1"><Label className="text-xs">End Date</Label><Input value={w.endDate ?? ''} onChange={(e) => updateWork(i, { endDate: e.target.value || undefined })} disabled={w.current} className="h-8 text-sm" /></div>
                  <div className="space-y-1"><Label className="text-xs">Logo URL</Label><Input value={w.logoUrl ?? ''} onChange={(e) => updateWork(i, { logoUrl: e.target.value || undefined })} className="h-8 text-sm" /></div>
                </div>
                <div className="grid grid-cols-2 gap-3">
                  <div className="space-y-1"><Label className="text-xs">Description (RU)</Label><Textarea rows={2} value={w.description.ru} onChange={(e) => updateWork(i, { description: { ...w.description, ru: e.target.value } })} className="text-sm" /></div>
                  <div className="space-y-1"><Label className="text-xs">Description (EN)</Label><Textarea rows={2} value={w.description.en} onChange={(e) => updateWork(i, { description: { ...w.description, en: e.target.value } })} className="text-sm" /></div>
                </div>
                <div className="space-y-1"><Label className="text-xs">Technologies (comma-separated)</Label><Input value={(w.technologies ?? []).join(', ')} onChange={(e) => updateWork(i, { technologies: e.target.value.split(',').map(s => s.trim()).filter(Boolean) })} className="h-8 text-sm" /></div>

                <div className="border-t pt-3 space-y-2">
                  <p className="text-xs font-semibold text-primary">Long-form description (expand/collapse on public page)</p>
                  <div className="grid grid-cols-2 gap-3">
                    <div className="space-y-1"><Label className="text-xs">Full description (RU)</Label><Textarea rows={4} value={w.fullDescription?.ru ?? ''} onChange={(e) => updateWork(i, { fullDescription: { en: w.fullDescription?.en ?? '', ru: e.target.value } })} className="text-sm" placeholder="Расширенный рассказ о позиции..." /></div>
                    <div className="space-y-1"><Label className="text-xs">Full description (EN)</Label><Textarea rows={4} value={w.fullDescription?.en ?? ''} onChange={(e) => updateWork(i, { fullDescription: { ru: w.fullDescription?.ru ?? '', en: e.target.value } })} className="text-sm" placeholder="Extended story about the role..." /></div>
                  </div>
                </div>

                <div className="border-t pt-3 space-y-2">
                  <div className="flex items-center justify-between">
                    <p className="text-xs font-semibold text-primary">Achievements (resume bullets)</p>
                    <Button
                      type="button"
                      variant="outline"
                      size="sm"
                      onClick={() => updateWork(i, { achievements: [...(w.achievements ?? []), { en: '', ru: '' }] })}
                    >
                      <Plus className="mr-1 h-3.5 w-3.5" />Add bullet
                    </Button>
                  </div>
                  {(w.achievements ?? []).map((ach, k) => (
                    <div key={k} className="border rounded p-2 space-y-2 bg-background/40">
                      <div className="flex items-center gap-2">
                        <span className="text-[10px] font-mono text-muted-foreground">#{k + 1}</span>
                        <Button
                          type="button"
                          variant="ghost"
                          size="sm"
                          onClick={() => updateWork(i, { achievements: (w.achievements ?? []).filter((_, j) => j !== k) })}
                          className="ml-auto text-destructive h-6 w-6 p-0"
                        >
                          <Trash2 className="h-3 w-3" />
                        </Button>
                      </div>
                      <div className="grid grid-cols-2 gap-2">
                        <Textarea
                          rows={2}
                          value={ach.ru}
                          onChange={(e) => {
                            const next = [...(w.achievements ?? [])];
                            next[k] = { ...ach, ru: e.target.value };
                            updateWork(i, { achievements: next });
                          }}
                          className="text-xs"
                          placeholder="RU — bullet из резюме"
                        />
                        <Textarea
                          rows={2}
                          value={ach.en}
                          onChange={(e) => {
                            const next = [...(w.achievements ?? [])];
                            next[k] = { ...ach, en: e.target.value };
                            updateWork(i, { achievements: next });
                          }}
                          className="text-xs"
                          placeholder="EN — resume bullet"
                        />
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            ))}
          </CardContent>
        </Card>

        <Card>
          <CardHeader><div className="flex items-center justify-between"><CardTitle>Certificates</CardTitle><Button variant="outline" size="sm" onClick={addCert}><Plus className="mr-1.5 h-3.5 w-3.5" />Add</Button></div></CardHeader>
          <CardContent className="space-y-4">
            {data.certificates.map((c, i) => (
              <div key={c.id} className="border rounded-lg p-4 space-y-3 bg-muted/20">
                <div className="flex items-center gap-2">
                  <span className="text-xs font-mono text-muted-foreground">{c.id}</span>
                  <Button variant="ghost" size="sm" onClick={() => setData((p) => p ? { ...p, certificates: p.certificates.filter((_, j) => j !== i) } : p)} className="ml-auto text-destructive h-7 w-7 p-0"><Trash2 className="h-3.5 w-3.5" /></Button>
                </div>
                <div className="grid grid-cols-2 gap-3">
                  <div className="space-y-1"><Label className="text-xs">Title (RU)</Label><Input value={c.title.ru} onChange={(e) => updateCert(i, { title: { ...c.title, ru: e.target.value } })} className="h-8 text-sm" /></div>
                  <div className="space-y-1"><Label className="text-xs">Title (EN)</Label><Input value={c.title.en} onChange={(e) => updateCert(i, { title: { ...c.title, en: e.target.value } })} className="h-8 text-sm" /></div>
                </div>
                <div className="grid grid-cols-2 gap-3">
                  <div className="space-y-1"><Label className="text-xs">Issuer (RU)</Label><Input value={c.issuer.ru} onChange={(e) => updateCert(i, { issuer: { ...c.issuer, ru: e.target.value } })} className="h-8 text-sm" /></div>
                  <div className="space-y-1"><Label className="text-xs">Issuer (EN)</Label><Input value={c.issuer.en} onChange={(e) => updateCert(i, { issuer: { ...c.issuer, en: e.target.value } })} className="h-8 text-sm" /></div>
                </div>
                <div className="grid grid-cols-3 gap-3">
                  <div className="space-y-1"><Label className="text-xs">Date</Label><Input value={c.date} onChange={(e) => updateCert(i, { date: e.target.value })} className="h-8 text-sm" placeholder="2023-05" /></div>
                  <div className="space-y-1"><Label className="text-xs">Credential ID</Label><Input value={c.credentialId ?? ''} onChange={(e) => updateCert(i, { credentialId: e.target.value || undefined })} className="h-8 text-sm" /></div>
                  <div className="space-y-1"><Label className="text-xs">URL</Label><Input value={c.url ?? ''} onChange={(e) => updateCert(i, { url: e.target.value || undefined })} className="h-8 text-sm" /></div>
                </div>
              </div>
            ))}
          </CardContent>
        </Card>

        <Card>
          <CardHeader><div className="flex items-center justify-between"><CardTitle>Publications</CardTitle><Button variant="outline" size="sm" onClick={addPub}><Plus className="mr-1.5 h-3.5 w-3.5" />Add</Button></div></CardHeader>
          <CardContent className="space-y-4">
            {data.publications.map((p, i) => (
              <div key={p.id} className="border rounded-lg p-4 space-y-3 bg-muted/20">
                <div className="flex items-center gap-2">
                  <span className="text-xs font-mono text-muted-foreground">{p.id}</span>
                  <Input value={p.year} onChange={(e) => updatePub(i, { year: e.target.value })} className="h-7 w-20 text-xs ml-auto" placeholder="2024" />
                  <Button variant="ghost" size="sm" onClick={() => setData((prev) => prev ? { ...prev, publications: prev.publications.filter((_, j) => j !== i) } : prev)} className="text-destructive h-7 w-7 p-0"><Trash2 className="h-3.5 w-3.5" /></Button>
                </div>
                <div className="grid grid-cols-2 gap-3">
                  <div className="space-y-1"><Label className="text-xs">Title (RU)</Label><Input value={p.title.ru} onChange={(e) => updatePub(i, { title: { ...p.title, ru: e.target.value } })} className="h-8 text-sm" /></div>
                  <div className="space-y-1"><Label className="text-xs">Title (EN)</Label><Input value={p.title.en} onChange={(e) => updatePub(i, { title: { ...p.title, en: e.target.value } })} className="h-8 text-sm" /></div>
                </div>
                <div className="grid grid-cols-2 gap-3">
                  <div className="space-y-1"><Label className="text-xs">Journal (RU)</Label><Input value={p.journal?.ru ?? ''} onChange={(e) => updatePub(i, { journal: { en: p.journal?.en ?? '', ru: e.target.value } })} className="h-8 text-sm" /></div>
                  <div className="space-y-1"><Label className="text-xs">Journal (EN)</Label><Input value={p.journal?.en ?? ''} onChange={(e) => updatePub(i, { journal: { ru: p.journal?.ru ?? '', en: e.target.value } })} className="h-8 text-sm" /></div>
                </div>
                <div className="grid grid-cols-2 gap-3">
                  <div className="space-y-1"><Label className="text-xs">DOI</Label><Input value={p.doi ?? ''} onChange={(e) => updatePub(i, { doi: e.target.value || undefined })} className="h-8 text-sm" /></div>
                  <div className="space-y-1"><Label className="text-xs">URL</Label><Input value={p.url ?? ''} onChange={(e) => updatePub(i, { url: e.target.value || undefined })} className="h-8 text-sm" /></div>
                </div>
                <div className="grid grid-cols-2 gap-3">
                  <div className="space-y-1"><Label className="text-xs">Abstract (RU)</Label><Textarea rows={2} value={p.abstract?.ru ?? ''} onChange={(e) => updatePub(i, { abstract: { en: p.abstract?.en ?? '', ru: e.target.value } })} className="text-sm" /></div>
                  <div className="space-y-1"><Label className="text-xs">Abstract (EN)</Label><Textarea rows={2} value={p.abstract?.en ?? ''} onChange={(e) => updatePub(i, { abstract: { ru: p.abstract?.ru ?? '', en: e.target.value } })} className="text-sm" /></div>
                </div>
              </div>
            ))}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
