import { Card } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog';
import { Loader2, ArrowLeft } from 'lucide-react';
import { useLocation } from 'wouter';
import { useState, useEffect } from 'react';
import { useVaultData, useVaultActions } from '@/hooks/useVaultData';
import { useEditMode, type ViewMode } from '@/hooks/useEditMode';
import { VaultDetailHeader } from '@/components/vault/vault-detail-header';
import { VaultMetadata } from '@/components/vault/vault-metadata';
import { VaultValueEditor } from '@/components/vault/vault-value-editor';
import DashboardHeader from '@/components/layout/dashboard-header';

// Helper function to parse URL query parameters
function getQueryParam(param: string): string | null {
  const urlParams = new URLSearchParams(window.location.search);
  return urlParams.get(param);
}

// Helper function to get initial mode from URL
function getInitialModeFromUrl(): ViewMode {
  const modeParam = getQueryParam('mode');
  return modeParam === 'edit' ? 'edit' : 'view';
}

interface VaultDetailContentProps {
  vaultId: string;
}

export default function VaultDetailContent({ vaultId }: VaultDetailContentProps) {
  const [, navigate] = useLocation();
  const [showUnsavedChangesDialog, setShowUnsavedChangesDialog] = useState(false);

  // Custom hooks for clean separation of concerns
  const vaultData = useVaultData(vaultId);
  const editMode = useEditMode(getInitialModeFromUrl());
  const vaultActions = useVaultActions({
    vault: vaultData.vault,
    originalValue: vaultData.originalValue,
    onSaveSuccess: editMode.exitEditMode,
  });

  const goBack = () => {
    // Check if user is in edit mode with unsaved changes
    if (editMode.isEditMode && vaultActions.hasUnsavedChanges) {
      setShowUnsavedChangesDialog(true);
      return; // Show dialog instead of navigating
    }
    navigate('/dashboard/vaults');
  };

  const handleConfirmLeave = () => {
    setShowUnsavedChangesDialog(false);
    navigate('/dashboard/vaults');
  };

  const handleCancelLeave = () => {
    setShowUnsavedChangesDialog(false);
  };

  // Sync URL when mode changes for proper browser history
  useEffect(() => {
    const currentUrl = new URL(window.location.href);
    const currentMode = currentUrl.searchParams.get('mode');
    const newMode = editMode.mode;

    // Only update URL if mode has actually changed
    if (currentMode !== newMode) {
      if (newMode === 'edit') {
        currentUrl.searchParams.set('mode', 'edit');
      } else {
        currentUrl.searchParams.delete('mode');
      }

      // Update URL without causing a page reload
      window.history.replaceState({}, '', currentUrl.toString());
    }
  }, [editMode.mode]);

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
        <div className="space-y-6">
          <VaultMetadata vault={vaultData.vault} isEditMode={editMode.isEditMode} />
          <VaultValueEditor isEditMode={editMode.isEditMode} vaultActions={vaultActions} />
        </div>
      </main>

      {/* Unsaved Changes Confirmation Dialog */}
      <AlertDialog open={showUnsavedChangesDialog} onOpenChange={setShowUnsavedChangesDialog}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Unsaved Changes</AlertDialogTitle>
            <AlertDialogDescription>
              You have unsaved changes that will be lost if you leave this page. Are you sure you want to continue?
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel onClick={handleCancelLeave}>
              Stay and Continue Editing
            </AlertDialogCancel>
            <AlertDialogAction onClick={handleConfirmLeave}>
              Leave Without Saving
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}
