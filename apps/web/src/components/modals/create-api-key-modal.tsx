import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { apiKeyApi } from '@/apis/api';
import { Loader2, X, Copy } from 'lucide-react';

interface CreateApiKeyModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onApiKeyCreated?: () => void;
}

export default function CreateApiKeyModal({ open, onOpenChange, onApiKeyCreated }: CreateApiKeyModalProps) {
  const [name, setName] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [generatedKey, setGeneratedKey] = useState<string | null>(null);

  const validate = (): string | null => {
    if (!name.trim()) return 'Name is required';
    if (name.length > 255) return 'Name must be ≤255 characters';
    return null;
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const validationError = validate();
    if (validationError) {
      setError(validationError);
      return;
    }

    setIsLoading(true);
    setError(null);
    try {
      const response = await apiKeyApi.createAPIKey({ name: name.trim() });
      const plainKey = (response as { key?: string; }).key ?? null;
      setGeneratedKey(plainKey);
      if (!plainKey) {
        setError('API key created but key not returned');
      }
      // Reset input
      setName('');
      onApiKeyCreated?.();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create API key');
    } finally {
      setIsLoading(false);
    }
  };

  const handleClose = () => {
    if (!isLoading) {
      onOpenChange(false);
      setError(null);
      setGeneratedKey(null);
      setName('');
    }
  };

  const copyKey = async () => {
    if (generatedKey) {
      try {
        await navigator.clipboard.writeText(generatedKey);
      } catch {
        // ignore errors
      }
    }
  };

  if (!open) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      <div className="fixed inset-0 bg-black/50" onClick={handleClose} />
      <Card className="relative z-10 w-full max-w-md mx-4">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-xl">Create API Key</CardTitle>
          <Button variant="ghost" size="icon" onClick={handleClose} disabled={isLoading}>
            <X className="h-4 w-4" />
          </Button>
        </CardHeader>
        <CardContent>
          {generatedKey ? (
            <div className="space-y-4">
              <p className="text-sm text-muted-foreground">
                This is your new API key. Please copy and store it securely – you won\'t be able to see it again.
              </p>
              <div className="flex items-center gap-2 bg-muted rounded-md p-2">
                <code className="flex-1 text-sm break-all">{generatedKey}</code>
                <Button variant="outline" size="icon" onClick={copyKey}>
                  <Copy className="h-4 w-4" />
                </Button>
              </div>
              <div className="flex justify-end pt-4">
                <Button onClick={handleClose}>Done</Button>
              </div>
            </div>
          ) : (
            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="name">Name *</Label>
                <Input
                  id="name"
                  placeholder="e.g., CI/CD Pipeline"
                  value={name}
                  onChange={(e: React.ChangeEvent<HTMLInputElement>) => setName(e.target.value)}
                  disabled={isLoading}
                  required
                />
              </div>

              {error && (
                <div className="p-3 text-sm text-red-600 bg-red-50 border border-red-200 rounded-md dark:text-red-400 dark:bg-red-900/20 dark:border-red-800">
                  {error}
                </div>
              )}

              <div className="flex justify-end space-x-2 pt-4">
                <Button type="button" variant="outline" onClick={handleClose} disabled={isLoading}>
                  Cancel
                </Button>
                <Button type="submit" disabled={isLoading}>
                  {isLoading ? (
                    <>
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                      Creating...
                    </>
                  ) : (
                    'Create Key'
                  )}
                </Button>
              </div>
            </form>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
