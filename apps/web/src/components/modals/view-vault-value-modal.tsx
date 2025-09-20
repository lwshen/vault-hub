import { useState, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Label } from '@/components/ui/label';
import { vaultApi } from '@/apis/api';
import { Copy, Loader2, X } from 'lucide-react';
import type { VaultLite } from '@lwshen/vault-hub-ts-fetch-client';

interface ViewVaultValueModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  vault: VaultLite | null;
}

export default function ViewVaultValueModal({ open, onOpenChange, vault }: ViewVaultValueModalProps) {
  const [value, setValue] = useState('');
  const [isFetchingValue, setIsFetchingValue] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [copySuccess, setCopySuccess] = useState(false);

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
      } else if (!open) {
        // Clear sensitive data when modal is closed
        setValue('');
        setError(null);
        setCopySuccess(false);
      }
    };

    fetchVaultValue();
  }, [vault, open]);

  const handleCopyToClipboard = async () => {
    try {
      await navigator.clipboard.writeText(value);
      setCopySuccess(true);
      setTimeout(() => setCopySuccess(false), 2000);
    } catch (err) {
      console.error('Failed to copy to clipboard:', err);
    }
  };

  const handleClose = () => {
    onOpenChange(false);
    setError(null);
    setCopySuccess(false);
    setValue(''); // Clear sensitive data from memory
  };

  if (!open || !vault) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      <div className="fixed inset-0 bg-black/50" onClick={handleClose} />
      <Card className="relative z-10 w-full max-w-md mx-4">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-xl">View Vault Value</CardTitle>
          <Button variant="ghost" size="icon" onClick={handleClose}>
            <X className="h-4 w-4" />
          </Button>
        </CardHeader>
        <CardContent>
          <div className="mb-4">
            <p className="text-sm text-muted-foreground">
              Viewing encrypted value for: <span className="font-medium">{vault.name}</span>
            </p>
          </div>

          {isFetchingValue ? (
            <div className="flex items-center justify-center py-8">
              <Loader2 className="h-6 w-6 animate-spin text-primary" />
              <span className="ml-2 text-sm text-muted-foreground">Loading vault value...</span>
            </div>
          ) : error ? (
            <div className="p-3 text-sm text-red-600 bg-red-50 border border-red-200 rounded-md dark:text-red-400 dark:bg-red-900/20 dark:border-red-800">
              {error}
            </div>
          ) : (
            <div className="space-y-4">
              <div className="space-y-2">
                <Label>Value</Label>
                <div className="relative">
                  <textarea
                    value={value}
                    readOnly
                    rows={4}
                    className="flex w-full rounded-md border border-input bg-muted px-3 py-2 text-sm shadow-xs resize-none cursor-text select-all"
                  />
                  <div className="absolute top-2 right-2">
                    <Button
                      type="button"
                      variant="ghost"
                      size="sm"
                      onClick={handleCopyToClipboard}
                      disabled={!value}
                      className="h-8 w-8 p-0"
                      title="Copy to clipboard"
                    >
                      <Copy className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
                {copySuccess && (
                  <p className="text-xs text-green-600 dark:text-green-400">
                    Copied to clipboard!
                  </p>
                )}
              </div>

              <div className="p-3 text-sm text-blue-700 bg-blue-50 border border-blue-200 rounded-md dark:text-blue-400 dark:bg-blue-900/20 dark:border-blue-800">
                <strong>Info:</strong> This value is decrypted and displayed in plain text. Use the copy icon to copy it to clipboard.
              </div>
            </div>
          )}

          <div className="flex justify-end pt-4">
            <Button onClick={handleClose}>Close</Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
