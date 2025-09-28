import { Card } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Loader2, ArrowLeft } from 'lucide-react';
import { useLocation } from 'wouter';
import { useVaultData, useVaultActions } from '@/hooks/useVaultData';
import { useEditMode } from '@/hooks/useEditMode';
import { VaultDetailHeader } from '@/components/vault/vault-detail-header';
import { VaultMetadata } from '@/components/vault/vault-metadata';
import { VaultValueEditor } from '@/components/vault/vault-value-editor';
import DashboardHeader from '@/components/layout/dashboard-header';

interface VaultDetailContentProps {
  vaultId: string;
}

export default function VaultDetailContent({ vaultId }: VaultDetailContentProps) {
  const [, navigate] = useLocation();

  // Custom hooks for clean separation of concerns
  const vaultData = useVaultData(vaultId);
  const editMode = useEditMode();
  const vaultActions = useVaultActions({
    vault: vaultData.vault,
    originalValue: vaultData.originalValue,
    onSaveSuccess: editMode.exitEditMode,
  });

  const goBack = () => {
    // Check if user is in edit mode with unsaved changes
    if (editMode.isEditMode && vaultActions.hasUnsavedChanges) {
      const confirmed = window.confirm(
        'You have unsaved changes. Are you sure you want to leave without saving? Your changes will be lost.',
      );
      if (!confirmed) {
        return; // Don't navigate if user cancels
      }
    }
    navigate('/dashboard/vaults');
  };

  // Loading state
  if (vaultData.isLoading) {
    return (
      <>
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
      </>
    );
  }

  // Error state
  if (vaultData.error && !vaultData.vault) {
    return (
      <>
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
                <p className="text-muted-foreground mb-4">{vaultData.error}</p>
                <Button onClick={vaultData.refetch}>Try Again</Button>
              </div>
            </div>
          </Card>
        </main>
      </>
    );
  }

  if (!vaultData.vault) return null;

  return (
    <>
      <VaultDetailHeader
        vault={vaultData.vault}
        editMode={editMode}
        vaultActions={vaultActions}
        onGoBack={goBack}
      />

      <main className="flex-1 overflow-y-auto p-4 sm:p-6">
        <div className="max-w-4xl mx-auto space-y-6">
          <VaultMetadata vault={vaultData.vault} isEditMode={editMode.isEditMode} />
          <VaultValueEditor isEditMode={editMode.isEditMode} vaultActions={vaultActions} />
        </div>
      </main>
    </>
  );
}
