'use client';

import { useEffect, useState, useRef } from 'react';
import { Button, Input, Textarea, Label, Card, CardContent, CardHeader, CardTitle } from '@/shared/ui';
import { Loader2, Save, Check, Plus, Trash2, ChevronUp, ChevronDown, Upload } from 'lucide-react';
import type { HomepageContent, HomepageVisibility, Achievement, HowIWorkStep } from '@/entities/content';
import { clientApi } from '@/shared/api';
import { AVAILABLE_ICONS, ACHIEVEMENT_TYPES } from '../model/constants';

export default function ContentEditor() {
  const [content, setContent] = useState<HomepageContent | null>(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [saved, setSaved] = useState(false);
  const [error, setError] = useState('');
  const [uploadingField, setUploadingField] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const pendingFieldRef = useRef<string | null>(null);

  useEffect(() => {
    clientApi
      .get<Array<{ section: string; data: unknown }>>('/api/content')
      .then((sections) => {
        const found = sections.find((s) => s.section === 'homepage');
        const d = (found?.data ?? {}) as unknown as HomepageContent;
        if (!d.howIWork) d.howIWork = { steps: [], philosophyTitleRu: '', philosophyTitleEn: '', philosophyTextRu: '', philosophyTextEn: '', payment1Ru: '', payment1En: '', payment1DescRu: '', payment1DescEn: '', payment2Ru: '', payment2En: '', payment2DescRu: '', payment2DescEn: '' };
        if (!d.howIWork.steps) d.howIWork.steps = [];
        if (!d.contactCTA) d.contactCTA = { headingRu: '', headingEn: '', subtitleRu: '', subtitleEn: '', bgWordRu: '', bgWordEn: '' };
        if (!d.marqueeItems) d.marqueeItems = [];
        if (!d.achievements) d.achievements = [];
        if (!d.socialLinks) d.socialLinks = [];
        if (!d.about) d.about = { gifUrl: '', bioEn: '', bioRu: '', stats: { projects: 0, yearsExperience: 0, certificates: 0, githubStars: 0, habrArticles: 0 } };
        if (!d.hero) d.hero = { name: { en: '', ru: '' }, titleAnimated: { en: [], ru: [] }, subtitleEn: '', subtitleRu: '', availableForWork: false };
        if (!d.visibility) d.visibility = {
          showAboutGif: true, showAboutBio: true,
          showStatProjects: true, showStatYears: true, showStatCertificates: true,
          showStatGithubStars: true, showStatHabrArticles: true,
          showProjectStars: true, showProjectStatus: true, showProjectTechStack: true,
          showMarquee: true, showHowIWork: true, showAchievements: true, showContactCTA: true,
        };
        setContent(d);
        setLoading(false);
      })
      .catch((err: unknown) => {
        setError(err instanceof Error ? err.message : 'Failed to load content');
        setLoading(false);
      });
  }, []);

  const handleSave = async () => {
    setSaving(true);
    try {
      await clientApi.put('/api/content?section=homepage', content);
      setSaved(true);
      setTimeout(() => setSaved(false), 2000);
    } catch {
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return <div className="flex justify-center py-20"><Loader2 className="h-6 w-6 animate-spin" /></div>;
  }
  if (error || !content) {
    return <div className="p-8"><p className="text-destructive">Error loading content: {error || 'No data'}</p></div>;
  }

  const set = (path: string, value: unknown) => {
    setContent((prev) => {
      if (!prev) return prev;
      const copy = JSON.parse(JSON.stringify(prev));
      const keys = path.split('.');
      let obj = copy;
      for (let i = 0; i < keys.length - 1; i++) obj = obj[keys[i]];
      obj[keys[keys.length - 1]] = value;
      return copy;
    });
  };

  const get = (path: string): any => {
    let obj: any = content;
    for (const k of path.split('.')) { obj = obj?.[k]; }
    return obj ?? '';
  };

  const triggerUpload = (fieldPath: string) => {
    pendingFieldRef.current = fieldPath;
    fileInputRef.current?.click();
  };

  const handleFileUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    const fieldPath = pendingFieldRef.current;
    if (!file || !fieldPath) return;

    setUploadingField(fieldPath);
    const formData = new FormData();
    formData.append('file', file);
    try {
      const data = await clientApi.post<{ url: string }>('/api/upload', formData);
      set(fieldPath, data.url);
    } catch {
    } finally {
      setUploadingField(null);
      pendingFieldRef.current = null;
      if (fileInputRef.current) fileInputRef.current.value = '';
    }
  };

  const socialLinks = content.socialLinks ?? [];
  const updateSocialLink = (i: number, field: string, value: string) => {
    const links = [...socialLinks];
    links[i] = { ...links[i], [field]: value };
    set('socialLinks', links);
  };
  const addSocialLink = () => set('socialLinks', [...socialLinks, { platform: '', url: '', icon: 'globe' }]);
  const removeSocialLink = (i: number) => set('socialLinks', socialLinks.filter((_: unknown, j: number) => j !== i));
  const moveSocialLink = (i: number, dir: -1 | 1) => {
    const links = [...socialLinks];
    const ni = i + dir;
    if (ni < 0 || ni >= links.length) return;
    [links[i], links[ni]] = [links[ni], links[i]];
    set('socialLinks', links);
  };

  const achievements = content.achievements ?? [];
  const addAchievement = () => set('achievements', [...achievements, {
    id: `ach-${Date.now()}`, title: { en: '', ru: '' }, description: { en: '', ru: '' }, year: new Date().getFullYear().toString(), type: 'achievement',
  }]);
  const removeAchievement = (i: number) => set('achievements', achievements.filter((_: Achievement, j: number) => j !== i));
  const updateAchievement = (i: number, updates: Partial<Achievement>) => {
    const arr = [...achievements];
    arr[i] = { ...arr[i], ...updates };
    set('achievements', arr);
  };
  const moveAchievement = (i: number, dir: -1 | 1) => {
    const arr = [...achievements];
    const ni = i + dir;
    if (ni < 0 || ni >= arr.length) return;
    [arr[i], arr[ni]] = [arr[ni], arr[i]];
    set('achievements', arr);
  };

  const steps = content.howIWork?.steps ?? [];
  const updateStep = (i: number, updates: Partial<HowIWorkStep>) => {
    const arr = [...steps];
    arr[i] = { ...arr[i], ...updates };
    set('howIWork.steps', arr);
  };

  const marqueeItems = content.marqueeItems ?? [];
  const updateMarquee = (i: number, value: string) => {
    const arr = [...marqueeItems];
    arr[i] = value;
    set('marqueeItems', arr);
  };
  const addMarquee = () => set('marqueeItems', [...marqueeItems, '']);
  const removeMarquee = (i: number) => set('marqueeItems', marqueeItems.filter((_: string, j: number) => j !== i));

  return (
    <div className="p-8 max-w-4xl">
      <input ref={fileInputRef} type="file" accept="image/*,video/*" className="hidden" onChange={handleFileUpload} />
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-3xl font-bold">Content Editor</h1>
        <Button onClick={handleSave} disabled={saving}>
          {saving ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : saved ? <Check className="mr-2 h-4 w-4" /> : <Save className="mr-2 h-4 w-4" />}
          {saved ? 'Saved!' : 'Save'}
        </Button>
      </div>

      <div className="space-y-6">

        <Card>
          <CardHeader><CardTitle>Visibility Settings</CardTitle></CardHeader>
          <CardContent>
            {(() => {
              const vis = content.visibility;
              if (!vis) return null;
              const Toggle = ({ field, label }: { field: keyof HomepageVisibility; label: string }) => (
                <label className="flex items-center gap-2.5 cursor-pointer select-none">
                  <input
                    type="checkbox"
                    checked={vis[field]}
                    onChange={(e) => set(`visibility.${field}`, e.target.checked)}
                    className="h-4 w-4 rounded accent-primary"
                  />
                  <span className="text-sm">{label}</span>
                </label>
              );
              return (
                <div className="space-y-5">
                  <div>
                    <p className="text-xs font-semibold text-muted-foreground uppercase tracking-wider mb-2">About Section</p>
                    <div className="grid grid-cols-2 sm:grid-cols-3 gap-y-2.5 gap-x-4">
                      <Toggle field="showAboutGif" label="Photo / GIF" />
                      <Toggle field="showAboutBio" label="Bio text" />
                      <Toggle field="showStatProjects" label="Stat: Projects" />
                      <Toggle field="showStatYears" label="Stat: Years exp." />
                      <Toggle field="showStatCertificates" label="Stat: Certificates" />
                      <Toggle field="showStatGithubStars" label="Stat: GitHub Stars" />
                      <Toggle field="showStatHabrArticles" label="Stat: Habr Articles" />
                    </div>
                  </div>
                  <div className="border-t pt-4">
                    <p className="text-xs font-semibold text-muted-foreground uppercase tracking-wider mb-2">Projects (Carousel)</p>
                    <div className="grid grid-cols-2 sm:grid-cols-3 gap-y-2.5 gap-x-4">
                      <Toggle field="showProjectStars" label="GitHub Stars" />
                      <Toggle field="showProjectStatus" label="Status Badge" />
                      <Toggle field="showProjectTechStack" label="Tech Stack Tags" />
                    </div>
                  </div>
                  <div className="border-t pt-4">
                    <p className="text-xs font-semibold text-muted-foreground uppercase tracking-wider mb-2">Homepage Sections</p>
                    <div className="grid grid-cols-2 sm:grid-cols-3 gap-y-2.5 gap-x-4">
                      <Toggle field="showMarquee" label="Marquee Ticker" />
                      <Toggle field="showHowIWork" label="How I Work" />
                      <Toggle field="showAchievements" label="Achievements" />
                      <Toggle field="showContactCTA" label="Contact CTA" />
                    </div>
                  </div>
                </div>
              );
            })()}
          </CardContent>
        </Card>

        <Card>
          <CardHeader><CardTitle>Hero Section</CardTitle></CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label>Name (EN)</Label>
                <Input value={get('hero.name.en')} onChange={(e) => set('hero.name.en', e.target.value)} />
              </div>
              <div className="space-y-2">
                <Label>Name (RU)</Label>
                <Input value={get('hero.name.ru')} onChange={(e) => set('hero.name.ru', e.target.value)} />
              </div>
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label>Subtitle (RU)</Label>
                <Input value={get('hero.subtitleRu')} onChange={(e) => set('hero.subtitleRu', e.target.value)} />
              </div>
              <div className="space-y-2">
                <Label>Subtitle (EN)</Label>
                <Input value={get('hero.subtitleEn')} onChange={(e) => set('hero.subtitleEn', e.target.value)} />
              </div>
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label>Animated Titles (RU) — comma-separated</Label>
                <Input
                  value={(get('hero.titleAnimated.ru') || []).join(', ')}
                  onChange={(e) => set('hero.titleAnimated.ru', e.target.value.split(',').map((s: string) => s.trim()).filter(Boolean))}
                />
              </div>
              <div className="space-y-2">
                <Label>Animated Titles (EN) — comma-separated</Label>
                <Input
                  value={(get('hero.titleAnimated.en') || []).join(', ')}
                  onChange={(e) => set('hero.titleAnimated.en', e.target.value.split(',').map((s: string) => s.trim()).filter(Boolean))}
                />
              </div>
            </div>
            <div className="flex items-center gap-3">
              <Label>Available for work</Label>
              <input type="checkbox" checked={get('hero.availableForWork') || false} onChange={(e) => set('hero.availableForWork', e.target.checked)} className="h-4 w-4" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader><CardTitle>About Section</CardTitle></CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-2">
              <Label>GIF / Photo URL</Label>
              <div className="flex gap-2">
                <Input value={get('about.gifUrl')} onChange={(e) => set('about.gifUrl', e.target.value)} placeholder="/images/casey.gif" className="flex-1" />
                <Button type="button" variant="outline" size="sm" onClick={() => triggerUpload('about.gifUrl')} disabled={uploadingField === 'about.gifUrl'}>
                  {uploadingField === 'about.gifUrl' ? <Loader2 className="h-4 w-4 animate-spin" /> : <Upload className="h-4 w-4" />}
                </Button>
              </div>
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label>Bio (RU) — short</Label>
                <Textarea rows={4} value={get('about.bioRu')} onChange={(e) => set('about.bioRu', e.target.value)} />
              </div>
              <div className="space-y-2">
                <Label>Bio (EN) — short</Label>
                <Textarea rows={4} value={get('about.bioEn')} onChange={(e) => set('about.bioEn', e.target.value)} />
              </div>
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label>Bio (RU) — full (expandable)</Label>
                <Textarea rows={8} value={get('about.bioFullRu')} onChange={(e) => set('about.bioFullRu', e.target.value)} />
              </div>
              <div className="space-y-2">
                <Label>Bio (EN) — full (expandable)</Label>
                <Textarea rows={8} value={get('about.bioFullEn')} onChange={(e) => set('about.bioFullEn', e.target.value)} />
              </div>
            </div>
            <div className="border-t pt-4">
              <Label className="text-sm font-semibold mb-3 block">Stats</Label>
              <div className="grid grid-cols-2 md:grid-cols-3 gap-3">
                <div className="space-y-1">
                  <Label className="text-xs">Projects</Label>
                  <Input type="number" value={get('about.stats.projects')} onChange={(e) => set('about.stats.projects', Number(e.target.value))} className="h-8" />
                </div>
                <div className="space-y-1">
                  <Label className="text-xs">Years of Experience</Label>
                  <Input type="number" value={get('about.stats.yearsExperience')} onChange={(e) => set('about.stats.yearsExperience', Number(e.target.value))} className="h-8" />
                </div>
                <div className="space-y-1">
                  <Label className="text-xs">Certificates</Label>
                  <Input type="number" value={get('about.stats.certificates')} onChange={(e) => set('about.stats.certificates', Number(e.target.value))} className="h-8" />
                </div>
                <div className="space-y-1">
                  <Label className="text-xs">GitHub Stars (0 = auto from API)</Label>
                  <Input type="number" value={get('about.stats.githubStars')} onChange={(e) => set('about.stats.githubStars', Number(e.target.value))} className="h-8" />
                </div>
                <div className="space-y-1">
                  <Label className="text-xs">Habr Articles</Label>
                  <Input type="number" value={get('about.stats.habrArticles')} onChange={(e) => set('about.stats.habrArticles', Number(e.target.value))} className="h-8" />
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle>Marquee Ticker</CardTitle>
              <Button variant="outline" size="sm" onClick={addMarquee}><Plus className="mr-1.5 h-3.5 w-3.5" />Add</Button>
            </div>
          </CardHeader>
          <CardContent>
            <div className="flex flex-wrap gap-2">
              {marqueeItems.map((item: string, i: number) => (
                <div key={i} className="flex items-center gap-1 bg-muted rounded-md pl-2">
                  <Input
                    value={item}
                    onChange={(e) => updateMarquee(i, e.target.value)}
                    className="h-7 w-32 text-xs border-0 bg-transparent p-0"
                  />
                  <button onClick={() => removeMarquee(i)} className="text-destructive hover:text-destructive/80 p-1">
                    <Trash2 className="h-3 w-3" />
                  </button>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader><CardTitle>How I Work — Steps</CardTitle></CardHeader>
          <CardContent className="space-y-4">
            {steps.map((step: HowIWorkStep, i: number) => (
              <div key={i} className="border rounded-lg p-4 space-y-3 bg-muted/20">
                <span className="text-sm font-bold text-muted-foreground">Step {i + 1}</span>
                <div className="grid grid-cols-2 gap-3">
                  <div className="space-y-1"><Label className="text-xs">Title (RU)</Label><Input value={step.titleRu} onChange={(e) => updateStep(i, { titleRu: e.target.value })} className="h-8 text-sm" /></div>
                  <div className="space-y-1"><Label className="text-xs">Title (EN)</Label><Input value={step.titleEn} onChange={(e) => updateStep(i, { titleEn: e.target.value })} className="h-8 text-sm" /></div>
                </div>
                <div className="grid grid-cols-2 gap-3">
                  <div className="space-y-1"><Label className="text-xs">Description (RU)</Label><Textarea rows={2} value={step.descRu} onChange={(e) => updateStep(i, { descRu: e.target.value })} className="text-sm" /></div>
                  <div className="space-y-1"><Label className="text-xs">Description (EN)</Label><Textarea rows={2} value={step.descEn} onChange={(e) => updateStep(i, { descEn: e.target.value })} className="text-sm" /></div>
                </div>
                <div className="space-y-1">
                  <Label className="text-xs">GIF URL</Label>
                  <div className="flex gap-2">
                    <Input value={step.gifUrl ?? ''} onChange={(e) => updateStep(i, { gifUrl: e.target.value })} className="h-8 text-sm flex-1" placeholder="/images/how-i-work/step1.gif" />
                    <Button type="button" variant="outline" size="sm" className="h-8 px-2" onClick={() => triggerUpload(`howIWork.steps.${i}.gifUrl`)} disabled={uploadingField === `howIWork.steps.${i}.gifUrl`}>
                      {uploadingField === `howIWork.steps.${i}.gifUrl` ? <Loader2 className="h-3 w-3 animate-spin" /> : <Upload className="h-3 w-3" />}
                    </Button>
                  </div>
                  {step.gifUrl && <img src={step.gifUrl} alt="preview" className="mt-1 h-16 rounded object-cover" />}
                </div>
              </div>
            ))}
          </CardContent>
        </Card>

        <Card>
          <CardHeader><CardTitle>How I Work — Philosophy & Payment</CardTitle></CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2"><Label>Philosophy Title (RU)</Label><Input value={get('howIWork.philosophyTitleRu')} onChange={(e) => set('howIWork.philosophyTitleRu', e.target.value)} /></div>
              <div className="space-y-2"><Label>Philosophy Title (EN)</Label><Input value={get('howIWork.philosophyTitleEn')} onChange={(e) => set('howIWork.philosophyTitleEn', e.target.value)} /></div>
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2"><Label>Philosophy Text (RU)</Label><Textarea rows={3} value={get('howIWork.philosophyTextRu')} onChange={(e) => set('howIWork.philosophyTextRu', e.target.value)} /></div>
              <div className="space-y-2"><Label>Philosophy Text (EN)</Label><Textarea rows={3} value={get('howIWork.philosophyTextEn')} onChange={(e) => set('howIWork.philosophyTextEn', e.target.value)} /></div>
            </div>
            <div className="border-t pt-4 grid grid-cols-2 gap-4">
              <div className="space-y-3">
                <Label className="font-semibold">Payment Card 1</Label>
                <div className="space-y-1"><Label className="text-xs">Title (RU)</Label><Input value={get('howIWork.payment1Ru')} onChange={(e) => set('howIWork.payment1Ru', e.target.value)} className="h-8 text-sm" /></div>
                <div className="space-y-1"><Label className="text-xs">Title (EN)</Label><Input value={get('howIWork.payment1En')} onChange={(e) => set('howIWork.payment1En', e.target.value)} className="h-8 text-sm" /></div>
                <div className="space-y-1"><Label className="text-xs">Desc (RU)</Label><Input value={get('howIWork.payment1DescRu')} onChange={(e) => set('howIWork.payment1DescRu', e.target.value)} className="h-8 text-sm" /></div>
                <div className="space-y-1"><Label className="text-xs">Desc (EN)</Label><Input value={get('howIWork.payment1DescEn')} onChange={(e) => set('howIWork.payment1DescEn', e.target.value)} className="h-8 text-sm" /></div>
              </div>
              <div className="space-y-3">
                <Label className="font-semibold">Payment Card 2</Label>
                <div className="space-y-1"><Label className="text-xs">Title (RU)</Label><Input value={get('howIWork.payment2Ru')} onChange={(e) => set('howIWork.payment2Ru', e.target.value)} className="h-8 text-sm" /></div>
                <div className="space-y-1"><Label className="text-xs">Title (EN)</Label><Input value={get('howIWork.payment2En')} onChange={(e) => set('howIWork.payment2En', e.target.value)} className="h-8 text-sm" /></div>
                <div className="space-y-1"><Label className="text-xs">Desc (RU)</Label><Input value={get('howIWork.payment2DescRu')} onChange={(e) => set('howIWork.payment2DescRu', e.target.value)} className="h-8 text-sm" /></div>
                <div className="space-y-1"><Label className="text-xs">Desc (EN)</Label><Input value={get('howIWork.payment2DescEn')} onChange={(e) => set('howIWork.payment2DescEn', e.target.value)} className="h-8 text-sm" /></div>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle>Achievements</CardTitle>
              <Button variant="outline" size="sm" onClick={addAchievement}><Plus className="mr-1.5 h-3.5 w-3.5" />Add</Button>
            </div>
          </CardHeader>
          <CardContent className="space-y-4">
            {achievements.length === 0 && <p className="text-sm text-muted-foreground text-center py-4">No achievements yet.</p>}
            {achievements.map((ach: Achievement, i: number) => (
              <div key={ach.id} className="border rounded-lg p-4 space-y-3 bg-muted/20">
                <div className="flex items-center gap-2">
                  <div className="flex flex-col gap-0.5">
                    <button onClick={() => moveAchievement(i, -1)} disabled={i === 0} className="text-muted-foreground hover:text-foreground disabled:opacity-20 p-0.5"><ChevronUp className="h-3.5 w-3.5" /></button>
                    <button onClick={() => moveAchievement(i, 1)} disabled={i === achievements.length - 1} className="text-muted-foreground hover:text-foreground disabled:opacity-20 p-0.5"><ChevronDown className="h-3.5 w-3.5" /></button>
                  </div>
                  <span className="text-xs font-mono text-muted-foreground">{ach.id}</span>
                  <div className="ml-auto flex items-center gap-2">
                    <select value={ach.type} onChange={(e) => updateAchievement(i, { type: e.target.value as Achievement['type'] })} className="h-7 rounded-md border bg-background px-2 text-xs">
                      {ACHIEVEMENT_TYPES.map((t) => <option key={t} value={t}>{t}</option>)}
                    </select>
                    <Input value={ach.year} onChange={(e) => updateAchievement(i, { year: e.target.value })} className="h-7 w-20 text-xs" placeholder="2024" />
                    <Button variant="ghost" size="sm" onClick={() => removeAchievement(i)} className="text-destructive h-7 w-7 p-0"><Trash2 className="h-3.5 w-3.5" /></Button>
                  </div>
                </div>
                <div className="grid grid-cols-2 gap-3">
                  <div className="space-y-1"><Label className="text-xs">Title (RU)</Label><Input value={ach.title.ru} onChange={(e) => updateAchievement(i, { title: { ...ach.title, ru: e.target.value } })} className="h-8 text-sm" /></div>
                  <div className="space-y-1"><Label className="text-xs">Title (EN)</Label><Input value={ach.title.en} onChange={(e) => updateAchievement(i, { title: { ...ach.title, en: e.target.value } })} className="h-8 text-sm" /></div>
                </div>
                <div className="grid grid-cols-2 gap-3">
                  <div className="space-y-1"><Label className="text-xs">Description (RU)</Label><Textarea rows={2} value={ach.description.ru} onChange={(e) => updateAchievement(i, { description: { ...ach.description, ru: e.target.value } })} className="text-sm" /></div>
                  <div className="space-y-1"><Label className="text-xs">Description (EN)</Label><Textarea rows={2} value={ach.description.en} onChange={(e) => updateAchievement(i, { description: { ...ach.description, en: e.target.value } })} className="text-sm" /></div>
                </div>
                <div className="space-y-1">
                  <Label className="text-xs">URL (optional)</Label>
                  <Input value={ach.url ?? ''} onChange={(e) => updateAchievement(i, { url: e.target.value || undefined })} className="h-8 text-sm" placeholder="https://..." />
                </div>
              </div>
            ))}
          </CardContent>
        </Card>

        <Card>
          <CardHeader><CardTitle>Contact CTA (Bottom of Homepage)</CardTitle></CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2"><Label>Heading (RU)</Label><Input value={get('contactCTA.headingRu')} onChange={(e) => set('contactCTA.headingRu', e.target.value)} /></div>
              <div className="space-y-2"><Label>Heading (EN)</Label><Input value={get('contactCTA.headingEn')} onChange={(e) => set('contactCTA.headingEn', e.target.value)} /></div>
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2"><Label>Subtitle (RU)</Label><Input value={get('contactCTA.subtitleRu')} onChange={(e) => set('contactCTA.subtitleRu', e.target.value)} /></div>
              <div className="space-y-2"><Label>Subtitle (EN)</Label><Input value={get('contactCTA.subtitleEn')} onChange={(e) => set('contactCTA.subtitleEn', e.target.value)} /></div>
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2"><Label>Background Word (RU)</Label><Input value={get('contactCTA.bgWordRu')} onChange={(e) => set('contactCTA.bgWordRu', e.target.value)} /></div>
              <div className="space-y-2"><Label>Background Word (EN)</Label><Input value={get('contactCTA.bgWordEn')} onChange={(e) => set('contactCTA.bgWordEn', e.target.value)} /></div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle>Social Links</CardTitle>
              <Button variant="outline" size="sm" onClick={addSocialLink}><Plus className="mr-1.5 h-3.5 w-3.5" />Add</Button>
            </div>
          </CardHeader>
          <CardContent className="space-y-3">
            {socialLinks.length === 0 && <p className="text-sm text-muted-foreground text-center py-4">No social links yet.</p>}
            {socialLinks.map((link: { platform: string; url: string; icon: string }, i: number) => (
              <div key={i} className="flex items-start gap-2 p-3 rounded-lg border bg-muted/30">
                <div className="flex flex-col gap-0.5 pt-1">
                  <button onClick={() => moveSocialLink(i, -1)} disabled={i === 0} className="text-muted-foreground hover:text-foreground disabled:opacity-20 text-xs">&#9650;</button>
                  <button onClick={() => moveSocialLink(i, 1)} disabled={i === socialLinks.length - 1} className="text-muted-foreground hover:text-foreground disabled:opacity-20 text-xs">&#9660;</button>
                </div>
                <div className="flex-1 grid grid-cols-1 sm:grid-cols-3 gap-2">
                  <div className="space-y-1"><Label className="text-xs">Platform</Label><Input value={link.platform} onChange={(e) => updateSocialLink(i, 'platform', e.target.value)} placeholder="GitHub" className="h-8 text-sm" /></div>
                  <div className="space-y-1"><Label className="text-xs">URL</Label><Input value={link.url} onChange={(e) => updateSocialLink(i, 'url', e.target.value)} placeholder="https://..." className="h-8 text-sm" /></div>
                  <div className="space-y-1">
                    <Label className="text-xs">Icon</Label>
                    <select value={link.icon} onChange={(e) => updateSocialLink(i, 'icon', e.target.value)} className="w-full h-8 rounded-md border bg-background px-2 text-sm">
                      {AVAILABLE_ICONS.map((icon) => <option key={icon} value={icon}>{icon}</option>)}
                    </select>
                  </div>
                </div>
                <Button variant="ghost" size="sm" onClick={() => removeSocialLink(i)} className="text-destructive hover:text-destructive h-8 w-8 p-0 mt-5"><Trash2 className="h-3.5 w-3.5" /></Button>
              </div>
            ))}
          </CardContent>
        </Card>

      </div>
    </div>
  );
}
