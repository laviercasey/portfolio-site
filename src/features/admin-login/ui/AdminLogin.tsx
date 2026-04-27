'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Lock, Loader2, Eye, EyeOff } from 'lucide-react';
import { Button, Input, Label } from '@/shared/ui';
import { adminLogin } from '../api/login';

interface Props {
  adminSlug: string;
}

export default function AdminLogin({ adminSlug }: Props) {
  const router = useRouter();
  const [password, setPassword] = useState('');
  const [showPass, setShowPass] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const adminBase = `/${adminSlug}`;

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    const result = await adminLogin(password);
    setLoading(false);
    if (!result.ok) {
      setError(result.error ?? 'Login failed');
      return;
    }
    router.push(`${adminBase}/dashboard`);
    router.refresh();
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-muted/30">
      <div className="w-full max-w-sm rounded-xl border bg-card p-8 shadow-lg">
        <div className="flex flex-col items-center mb-8">
          <div className="p-3 rounded-full bg-primary/10 mb-4">
            <Lock className="h-6 w-6 text-primary" />
          </div>
          <h1 className="text-2xl font-bold">Admin Panel</h1>
          <p className="text-sm text-muted-foreground mt-1">Enter your password to continue</p>
        </div>

        <form onSubmit={handleLogin} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="password">Password</Label>
            <div className="relative">
              <Input
                id="password"
                type={showPass ? 'text' : 'password'}
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="••••••••"
                autoComplete="current-password"
                className="pr-10"
              />
              <button
                type="button"
                onClick={() => setShowPass(!showPass)}
                className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
              >
                {showPass ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
              </button>
            </div>
          </div>

          {error && (
            <p className="text-sm text-destructive bg-destructive/10 rounded-md px-3 py-2">{error}</p>
          )}

          <Button type="submit" className="w-full" disabled={loading || !password}>
            {loading ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : null}
            Log in
          </Button>
        </form>
      </div>
    </div>
  );
}
