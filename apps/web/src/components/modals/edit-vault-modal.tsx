import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { vaultApi } from '@/apis/api';
import { Loader2, X } from 'lucide-react';
import type { VaultLite } from '@lwshen/vault-hub-ts-fetch-client';

interface EditVaultModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  vault: VaultLite | null;
  onVaultUpdated?: () => void;
}

interface FormData {
  name: string;
  description: string;
  category: string;
}

export default function EditVaultModal({ open, onOpenChange, vault, onVaultUpdated }: EditVaultModalProps) {
  const [formData, setFormData] = useState<FormData>({
    name: vault?.name ?? '',
    description: vault?.description ?? '',
    category: vault?.category ?? '',
  });
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Use vault data directly when available, otherwise use form state
  const displayData = vault && open ? {
    name: vault.name || '',
    description: vault.description || '',
    category: vault.category || '',
  } : formData;

  const handleInputChange = (field: keyof FormData, value: string) => {
    setFormData(prev => ({ ...prev, [field]: value }));
    if (error) setError(null);
  };

  const validateForm = (data: FormData): string | null => {
    const validations = [
      { check: !data.name.trim(), message: 'Name is required' },
      { check: data.name.length > 100, message: 'Name must be ≤100 characters' },
      { check: data.description?.length > 500, message: 'Description must be ≤500 characters' },
    ];

    const failed = validations.find(v => v.check);
    return failed?.message || null;
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    if (!vault) return;

    const validationError = validateForm(displayData);
    if (validationError) {
      setError(validationError);
      return;
    }

    setIsLoading(true);
    setError(null);

    try {
      await vaultApi.updateVault(vault.uniqueId, {
        name: displayData.name.trim(),
        description: displayData.description.trim() || undefined,
        category: displayData.category.trim() || undefined,
      });

      onVaultUpdated?.();
      onOpenChange(false);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to update vault');
    } finally {
      setIsLoading(false);
    }
  };

  const handleClose = () => {
    if (!isLoading) {
      onOpenChange(false);
      setError(null);
    }
  };

  if (!open || !vault) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      <div className="fixed inset-0 bg-black/50" onClick={handleClose} />
      <Card className="relative z-10 w-full max-w-lg mx-4">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-xl">Edit Vault</CardTitle>
          <Button variant="ghost" size="icon" onClick={handleClose} disabled={isLoading}>
            <X className="h-4 w-4" />
          </Button>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="name">Name *</Label>
              <Input
                id="name"
                placeholder="e.g., Production API Keys"
                value={displayData.name}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => handleInputChange('name', e.target.value)}
                disabled={isLoading}
                required
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="category">Category</Label>
              <Input
                id="category"
                placeholder="e.g., API Keys, Database, Certificates"
                value={displayData.category}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => handleInputChange('category', e.target.value)}
                disabled={isLoading}
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="description">Description</Label>
              <textarea
                id="description"
                placeholder="Optional description for this vault"
                value={displayData.description}
                onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) => handleInputChange('description', e.target.value)}
                disabled={isLoading}
                rows={2}
                className="flex w-full rounded-md border border-input bg-transparent px-3 py-2 text-sm shadow-xs placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
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
                    Saving...
                  </>
                ) : (
                  'Save'
                )}
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
