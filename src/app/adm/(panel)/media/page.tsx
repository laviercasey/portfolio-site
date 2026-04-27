'use client';

import { useState, useRef } from 'react';
import { Button } from '@/shared/ui';
import { Card, CardContent } from '@/shared/ui';
import { Upload, Loader2, Copy, Check, Image as ImageIcon, Film } from 'lucide-react';

interface UploadedFile {
  url: string;
  filename: string;
  type: 'image' | 'video';
}

export default function MediaPage() {
  const [files, setFiles] = useState<UploadedFile[]>([]);
  const [uploading, setUploading] = useState(false);
  const [copiedUrl, setCopiedUrl] = useState<string | null>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  const handleUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    setUploading(true);
    const formData = new FormData();
    formData.append('file', file);

    try {
      const res = await fetch('/api/upload', { method: 'POST', body: formData });
      if (res.ok) {
        const data = await res.json();
        setFiles((prev) => [{
          url: data.url,
          filename: data.filename,
          type: file.type.startsWith('video/') ? 'video' : 'image',
        }, ...prev]);
      }
    } finally {
      setUploading(false);
      if (inputRef.current) inputRef.current.value = '';
    }
  };

  const copyUrl = (url: string) => {
    navigator.clipboard.writeText(url);
    setCopiedUrl(url);
    setTimeout(() => setCopiedUrl(null), 1500);
  };

  return (
      <div className="p-8">
        <div className="flex items-center justify-between mb-6">
          <h1 className="text-3xl font-bold">Media Library</h1>
          <div>
            <input ref={inputRef} type="file" accept="image/*,video/*" className="hidden" onChange={handleUpload} />
            <Button onClick={() => inputRef.current?.click()} disabled={uploading}>
              {uploading ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : <Upload className="mr-2 h-4 w-4" />}
              Upload file
            </Button>
          </div>
        </div>

        <div className="rounded-xl border-2 border-dashed p-8 text-center mb-8 cursor-pointer hover:bg-muted/30 transition-colors" onClick={() => inputRef.current?.click()}>
          <Upload className="h-8 w-8 text-muted-foreground mx-auto mb-2" />
          <p className="text-muted-foreground">Click to upload images or videos (max 50MB)</p>
          <p className="text-xs text-muted-foreground mt-1">JPG, PNG, WebP, GIF, MP4, WebM</p>
        </div>

        {files.length === 0 ? (
          <p className="text-center text-muted-foreground py-8">No files uploaded yet this session</p>
        ) : (
          <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
            {files.map((file) => (
              <Card key={file.url}>
                <CardContent className="p-3">
                  <div className="rounded-lg bg-muted aspect-video flex items-center justify-center mb-2 overflow-hidden">
                    {file.type === 'image' ? (
                      <img src={file.url} alt={file.filename} className="w-full h-full object-cover" />
                    ) : (
                      <Film className="h-8 w-8 text-muted-foreground" />
                    )}
                  </div>
                  <p className="text-xs text-muted-foreground truncate mb-2">{file.filename}</p>
                  <Button variant="outline" size="sm" className="w-full" onClick={() => copyUrl(file.url)}>
                    {copiedUrl === file.url ? <Check className="mr-1.5 h-3 w-3" /> : <Copy className="mr-1.5 h-3 w-3" />}
                    {copiedUrl === file.url ? 'Copied!' : 'Copy URL'}
                  </Button>
                </CardContent>
              </Card>
            ))}
          </div>
        )}
      </div>
  );
}
