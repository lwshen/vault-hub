import DashboardHeader from '@/components/layout/dashboard-header';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Label } from '@/components/ui/label';
import { vaultApi } from '@/apis/api';
import type { Vault } from '@lwshen/vault-hub-ts-fetch-client';
import {
  ArrowLeft,
  Copy,
  Edit3,
  Eye,
  Loader2,
  Save,
  X,
} from 'lucide-react';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { useLocation } from 'wouter';

interface VaultDetailProps {
  vaultId: string;
}

type ViewMode = 'view' | 'edit';

export default function VaultDetail({ vaultId }: VaultDetailProps) {
  // Get mode from URL search params
  const [location, navigate] = useLocation();
  const searchParams = new URLSearchParams(location.split('?')[1] || '');
  const mode = (searchParams.get('mode') as ViewMode) || 'view';

  // State management
  const [vault, setVault] = useState<Vault | null>(null);
  const [originalValue, setOriginalValue] = useState('');
  const [editedValue, setEditedValue] = useState('');
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Fetch vault data
  useEffect(() => {
    const fetchVault = async () => {
      if (!vaultId) return;

      setIsLoading(true);
      setError(null);

      try {
        const fullVault = await vaultApi.getVault(vaultId);
        setVault(fullVault);
        setOriginalValue(fullVault.value || '');
        setEditedValue(fullVault.value || '');
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to fetch vault');
      } finally {
        setIsLoading(false);
      }
    };

    fetchVault();
  }, [vaultId]);

  // Navigation helpers
  const goBack = () => {
    navigate('/dashboard/vaults');
  };

  const setMode = (newMode: ViewMode) => {
    const params = new URLSearchParams(location.split('?')[1] || '');
    if (newMode === 'view') {
      params.delete('mode');
    } else {
      params.set('mode', newMode);
    }
    const queryString = params.toString();
    navigate(`/dashboard/vaults/${vaultId}${queryString ? `?${queryString}` : ''}`);
  };

  // Actions
  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(originalValue);
      toast.success(`Copied value for "${vault?.name}" to clipboard`);
    } catch (err) {
      console.error('Failed to copy to clipboard:', err);
      toast.error('Failed to copy to clipboard');
    }
  };

  const handleSave = async () => {
    if (!vault) return;

    if (!editedValue.trim()) {
      setError('Value is required');
      return;
    }

    setIsSaving(true);
    setError(null);

    try {
      await vaultApi.updateVault(vault.uniqueId, {
        value: editedValue.trim(),
      });

      setOriginalValue(editedValue.trim());
      toast.success('Vault value updated successfully');
      setMode('view');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to update vault value');
    } finally {
      setIsSaving(false);
    }
  };

  const handleCancel = () => {
    setEditedValue(originalValue);
    setError(null);
    setMode('view');
  };

  const hasUnsavedChanges = mode === 'edit' && editedValue !== originalValue;

  // Loading state
  if (isLoading) {
    return (
      <div className="flex flex-col min-h-screen">
        <DashboardHeader
          title="Loading..."
          description="Fetching vault details"
          actions={
            <Button variant="outline" size="sm" onClick={goBack}>
              <ArrowLeft className="h-4 w-4 mr-2" />
              <span className="hidden sm:inline">Back to Vaults</span>
              <span className="sm:hidden">Back</span>
            </Button>
          }
        />
        <main className="flex-1 overflow-y-auto p-4 sm:p-6">
          <Card className="p-6">
            <div className="flex items-center justify-center min-h-[200px] flex-col gap-4">
              <Loader2 className="h-8 w-8 animate-spin text-primary" />
              <p className="text-muted-foreground">Loading vault details...</p>
            </div>
          </Card>
        </main>
      </div>
    );
  }

  // Error state
  if (error && !vault) {
    return (
      <div className="flex flex-col min-h-screen">
        <DashboardHeader
          title="Error"
          description="Failed to load vault"
          actions={
            <Button variant="outline" size="sm" onClick={goBack}>
              <ArrowLeft className="h-4 w-4 mr-2" />
              <span className="hidden sm:inline">Back to Vaults</span>
              <span className="sm:hidden">Back</span>
            </Button>
          }
        />
        <main className="flex-1 overflow-y-auto p-4 sm:p-6">
          <Card className="p-6">
            <div className="flex items-center justify-center min-h-[200px] flex-col gap-4">
              <div className="text-center">
                <h3 className="text-lg font-semibold text-red-600">Failed to load vault</h3>
                <p className="text-muted-foreground mb-4">{error}</p>
                <Button onClick={() => window.location.reload()}>Try Again</Button>
              </div>
            </div>
          </Card>
        </main>
      </div>
    );
  }

  if (!vault) return null;

  const isEditMode = mode === 'edit';

  return (
    <div className="flex flex-col min-h-screen">
      {/* Header with responsive navigation */}
      <DashboardHeader
        title={vault.name}
        description={`${isEditMode ? 'Editing' : 'Viewing'} vault value`}
        actions={
          <div className="flex items-center gap-2">
            {/* Mobile: Stack buttons vertically in dropdown or show essential ones */}
            <div className="flex items-center gap-2">
              <Button variant="outline" size="sm" onClick={goBack}>
                <ArrowLeft className="h-4 w-4 mr-0 sm:mr-2" />
                <span className="hidden sm:inline">Back to Vaults</span>
              </Button>

              {!isEditMode && (
                <>
                  <Button variant="outline" size="sm" onClick={handleCopy}>
                    <Copy className="h-4 w-4 mr-0 sm:mr-2" />
                    <span className="hidden sm:inline">Copy</span>
                  </Button>
                  <Button variant="default" size="sm" onClick={() => setMode('edit')}>
                    <Edit3 className="h-4 w-4 mr-0 sm:mr-2" />
                    <span className="hidden sm:inline">Edit</span>
                  </Button>
                </>
              )}

              {isEditMode && (
                <>
                  <Button variant="outline" size="sm" onClick={handleCancel} disabled={isSaving}>
                    <X className="h-4 w-4 mr-0 sm:mr-2" />
                    <span className="hidden sm:inline">Cancel</span>
                  </Button>
                  <Button variant="default" size="sm" onClick={handleSave} disabled={isSaving}>
                    {isSaving ? (
                      <Loader2 className="h-4 w-4 mr-0 sm:mr-2 animate-spin" />
                    ) : (
                      <Save className="h-4 w-4 mr-0 sm:mr-2" />
                    )}
                    <span className="hidden sm:inline">{isSaving ? 'Saving...' : 'Save'}</span>
                  </Button>
                </>
              )}
            </div>
          </div>
        }
      />

      {/* Main content area with responsive layout */}
      <main className="flex-1 overflow-y-auto p-4 sm:p-6">
        <div className="max-w-4xl mx-auto space-y-6">
          {/* Vault metadata card */}
          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="flex items-center gap-2 text-lg sm:text-xl">
                {isEditMode ? <Edit3 className="h-5 w-5" /> : <Eye className="h-5 w-5" />}
                {vault.name}
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
                {vault.category && (
                  <div>
                    <Label className="text-xs uppercase tracking-wide text-muted-foreground">Category</Label>
                    <p className="mt-1 text-sm font-medium">{vault.category}</p>
                  </div>
                )}
                {vault.description && (
                  <div className="sm:col-span-2">
                    <Label className="text-xs uppercase tracking-wide text-muted-foreground">Description</Label>
                    <p className="mt-1 text-sm">{vault.description}</p>
                  </div>
                )}
                {vault.updatedAt && (
                  <div>
                    <Label className="text-xs uppercase tracking-wide text-muted-foreground">Last Updated</Label>
                    <p className="mt-1 text-sm">{new Date(vault.updatedAt).toLocaleDateString()}</p>
                  </div>
                )}
              </div>
            </CardContent>
          </Card>

          {/* Vault value card with responsive textarea */}
          <Card>
            <CardHeader>
              <CardTitle className="text-lg">
                {isEditMode ? 'Edit Value' : 'Vault Value'}
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="vault-value">
                    {isEditMode ? 'Encrypted Value *' : 'Decrypted Value'}
                  </Label>
                  <textarea
                    id="vault-value"
                    value={isEditMode ? editedValue : originalValue}
                    onChange={(e) => {
                      if (isEditMode) {
                        setEditedValue(e.target.value);
                        if (error) setError(null);
                      }
                    }}
                    placeholder={isEditMode ? 'Enter the secret value to be encrypted and stored' : ''}
                    readOnly={!isEditMode}
                    rows={6} // Mobile-friendly height
                    className={`
                      flex w-full rounded-md border border-input px-3 py-2 text-sm shadow-xs resize-none
                      sm:rows-8 lg:rows-12
                      ${isEditMode
      ? 'bg-transparent placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring'
      : 'bg-muted cursor-text select-all'
    }
                      ${!isEditMode && 'cursor-text select-all'}
                    `}
                    style={{
                      minHeight: isEditMode ? '200px' : '150px', // Responsive min-height
                    }}
                  />
                </div>

                {/* Warning messages */}
                {isEditMode && (
                  <div className="p-3 text-sm text-amber-700 bg-amber-50 border border-amber-200 rounded-md dark:text-amber-400 dark:bg-amber-900/20 dark:border-amber-800">
                    <strong>Warning:</strong> This will replace the current encrypted value. This action cannot be undone.
                  </div>
                )}

                {!isEditMode && (
                  <div className="p-3 text-sm text-blue-700 bg-blue-50 border border-blue-200 rounded-md dark:text-blue-400 dark:bg-blue-900/20 dark:border-blue-800">
                    <strong>Info:</strong> This value is decrypted and displayed in plain text. Use the copy button to copy it to clipboard.
                  </div>
                )}

                {/* Error display */}
                {error && (
                  <div className="p-3 text-sm text-red-600 bg-red-50 border border-red-200 rounded-md dark:text-red-400 dark:bg-red-900/20 dark:border-red-800">
                    {error}
                  </div>
                )}

                {/* Unsaved changes indicator */}
                {hasUnsavedChanges && (
                  <div className="p-3 text-sm text-orange-700 bg-orange-50 border border-orange-200 rounded-md dark:text-orange-400 dark:bg-orange-900/20 dark:border-orange-800">
                    <strong>Unsaved changes:</strong> You have modified the vault value. Save or cancel to continue.
                  </div>
                )}
              </div>
            </CardContent>
          </Card>
        </div>
      </main>

      {/* Mobile: Sticky action bar at bottom for better thumb access */}
      <div className="sm:hidden border-t border-border bg-background p-4">
        <div className="flex gap-2">
          {!isEditMode ? (
            <>
              <Button variant="outline" size="lg" onClick={handleCopy} className="flex-1">
                <Copy className="h-4 w-4 mr-2" />
                Copy Value
              </Button>
              <Button variant="default" size="lg" onClick={() => setMode('edit')} className="flex-1">
                <Edit3 className="h-4 w-4 mr-2" />
                Edit Value
              </Button>
            </>
          ) : (
            <>
              <Button variant="outline" size="lg" onClick={handleCancel} disabled={isSaving} className="flex-1">
                <X className="h-4 w-4 mr-2" />
                Cancel
              </Button>
              <Button variant="default" size="lg" onClick={handleSave} disabled={isSaving} className="flex-1">
                {isSaving ? (
                  <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                ) : (
                  <Save className="h-4 w-4 mr-2" />
                )}
                {isSaving ? 'Saving...' : 'Save Changes'}
              </Button>
            </>
          )}
        </div>
      </div>
    </div>
  );
}
