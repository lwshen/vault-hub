import DashboardHeader from '@/components/layout/dashboard-header';
import { Button } from '@/components/ui/button';
import {
  ArrowLeft,
  Copy,
  Edit3,
  Loader2,
  Save,
  X,
} from 'lucide-react';
import type { Vault } from '@lwshen/vault-hub-ts-fetch-client';
import type { UseEditModeReturn } from '@/hooks/useEditMode';
import type { UseVaultActionsReturn } from '@/hooks/useVaultData';

interface VaultDetailHeaderProps {
  vault: Vault;
  editMode: UseEditModeReturn;
  vaultActions: UseVaultActionsReturn;
  onGoBack: () => void;
}

export function VaultDetailHeader({
  vault,
  editMode,
  vaultActions,
  onGoBack,
}: VaultDetailHeaderProps) {
  const { isEditMode, enterEditMode, exitEditMode } = editMode;
  const { handleSave, handleCopy, resetChanges, isSaving } = vaultActions;

  const handleSaveAndExit = async () => {
    await handleSave();
    exitEditMode();
  };

  const handleCancel = () => {
    resetChanges();
    exitEditMode();
  };

  return (
    <>
      {/* Desktop Header */}
      <DashboardHeader
        title={vault.name}
        description={`${isEditMode ? 'Editing' : 'Viewing'} vault value`}
        actions={
          <div className="flex items-center gap-2">
            <Button variant="outline" size="sm" onClick={onGoBack}>
              <ArrowLeft className="h-4 w-4 mr-0 sm:mr-2" />
              <span className="hidden sm:inline">Back to Vaults</span>
            </Button>

            {!isEditMode && (
              <>
                <Button variant="outline" size="sm" onClick={handleCopy}>
                  <Copy className="h-4 w-4 mr-0 sm:mr-2" />
                  <span className="hidden sm:inline">Copy</span>
                </Button>
                <Button variant="default" size="sm" onClick={enterEditMode}>
                  <Edit3 className="h-4 w-4 mr-0 sm:mr-2" />
                  <span className="hidden sm:inline">Edit</span>
                </Button>
              </>
            )}

            {isEditMode && (
              <>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={handleCancel}
                  disabled={isSaving}
                >
                  <X className="h-4 w-4 mr-0 sm:mr-2" />
                  <span className="hidden sm:inline">Cancel</span>
                </Button>
                <Button
                  variant="default"
                  size="sm"
                  onClick={handleSaveAndExit}
                  disabled={isSaving}
                >
                  {isSaving ? (
                    <Loader2 className="h-4 w-4 mr-0 sm:mr-2 animate-spin" />
                  ) : (
                    <Save className="h-4 w-4 mr-0 sm:mr-2" />
                  )}
                  <span className="hidden sm:inline">
                    {isSaving ? 'Saving...' : 'Save'}
                  </span>
                </Button>
              </>
            )}
          </div>
        }
      />

      {/* Mobile: Sticky action bar at bottom for better thumb access */}
      <div className="sm:hidden border-t border-border bg-background p-4">
        <div className="flex gap-2">
          {!isEditMode ? (
            <>
              <Button
                variant="outline"
                size="lg"
                onClick={handleCopy}
                className="flex-1"
              >
                <Copy className="h-4 w-4 mr-2" />
                Copy Value
              </Button>
              <Button
                variant="default"
                size="lg"
                onClick={enterEditMode}
                className="flex-1"
              >
                <Edit3 className="h-4 w-4 mr-2" />
                Edit Value
              </Button>
            </>
          ) : (
            <>
              <Button
                variant="outline"
                size="lg"
                onClick={handleCancel}
                disabled={isSaving}
                className="flex-1"
              >
                <X className="h-4 w-4 mr-2" />
                Cancel
              </Button>
              <Button
                variant="default"
                size="lg"
                onClick={handleSaveAndExit}
                disabled={isSaving}
                className="flex-1"
              >
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
    </>
  );
}
