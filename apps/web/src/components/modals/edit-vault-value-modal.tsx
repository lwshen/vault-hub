import { useState, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Label } from '@/components/ui/label';
import { vaultApi } from '@/apis/api';
import { Loader2, X } from 'lucide-react';
import type { VaultLite } from '@lwshen/vault-hub-ts-fetch-client';

interface EditVaultValueModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  vault: VaultLite | null;
  onVaultUpdated?: () => void;
}

export default function EditVaultValueModal({ open, onOpenChange, vault, onVaultUpdated }: EditVaultValueModalProps) {
  const [value, setValue] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [isFetchingValue, setIsFetchingValue] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchVaultValue = async () => {
      if (vault && open) {
        setIsFetchingValue(true);
        setError(null);
        try {
          const fullVault = await vaultApi.getVault(vault.uniqueId);
          setValue(fullVault.value || '');
        } catch (err) {
          setError(err instanceof Error ? err.message : 'Failed to fetch vault value');
        } finally {
          setIsFetchingValue(false);
        }
      }
    };

    fetchVaultValue();
  }, [vault, open]);

  const validateForm = (): string | null => {
    if (!value.trim()) return 'Value is required';
    return null;
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    if (!vault) return;

    const validationError = validateForm();
    if (validationError) {
      setError(validationError);
      return;
    }

    setIsLoading(true);
    setError(null);

    try {
      await vaultApi.updateVault(vault.uniqueId, {
        value: value.trim(),
      });

      onVaultUpdated?.();
      onOpenChange(false);
      setValue(''); // Clear form after success
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to update vault value');
    } finally {
      setIsLoading(false);
    }
  };

  const handleClose = () => {
    if (!isLoading) {
      onOpenChange(false);
      setValue('');
      setError(null);
    }
  };

  if (!open || !vault) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      <div className="fixed inset-0 bg-black/50" onClick={handleClose} />
      <Card className="relative z-10 w-full max-w-md mx-4">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-xl">Edit Vault Value</CardTitle>
          <Button variant="ghost" size="icon" onClick={handleClose} disabled={isLoading}>
            <X className="h-4 w-4" />
          </Button>
        </CardHeader>
        <CardContent>
          <div className="mb-4">
            <p className="text-sm text-muted-foreground">
              Updating the encrypted value for: <span className="font-medium">{vault.name}</span>
            </p>
          </div>

          {isFetchingValue ? (
            <div className="flex items-center justify-center py-8">
              <Loader2 className="h-6 w-6 animate-spin text-primary" />
              <span className="ml-2 text-sm text-muted-foreground">Loading current value...</span>
            </div>
          ) : (
            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="value">Current Value *</Label>
                <textarea
                  id="value"
                  placeholder="Enter the new secret value to be encrypted and stored"
                  value={value}
                  onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) => {
                    setValue(e.target.value);
                    if (error) setError(null);
                  }}
                  disabled={isLoading}
                  required
                  rows={4}
                  className="flex w-full rounded-md border border-input bg-transparent px-3 py-2 text-sm shadow-xs placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
                />
              </div>

              <div className="p-3 text-sm text-amber-700 bg-amber-50 border border-amber-200 rounded-md dark:text-amber-400 dark:bg-amber-900/20 dark:border-amber-800">
                <strong>Warning:</strong> This will replace the current encrypted value. This action cannot be undone.
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
                      Updating...
                    </>
                  ) : (
                    'Update Value'
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
