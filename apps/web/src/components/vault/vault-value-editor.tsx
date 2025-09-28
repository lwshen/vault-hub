import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Label } from '@/components/ui/label';
import { StatusAlert } from '@/components/ui/status-alert';
import type { UseVaultActionsReturn } from '@/hooks/useVaultData';

interface VaultValueEditorProps {
  isEditMode: boolean;
  vaultActions: UseVaultActionsReturn;
}

export function VaultValueEditor({ isEditMode, vaultActions }: VaultValueEditorProps) {
  const {
    editedValue,
    setEditedValue,
    error,
    setError,
    hasUnsavedChanges,
  } = vaultActions;

  return (
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
              value={isEditMode ? editedValue : vaultActions.editedValue}
              onChange={(e) => {
                if (isEditMode) {
                  setEditedValue(e.target.value);
                  if (error) setError(null);
                }
              }}
              placeholder={isEditMode ? 'Enter the secret value to be encrypted and stored' : ''}
              readOnly={!isEditMode}
              rows={6}
              className={`
                flex w-full rounded-md border border-input px-3 py-2 text-sm shadow-xs resize-none
                ${isEditMode
      ? 'bg-transparent placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring'
      : 'bg-muted cursor-text select-all'
    }
              `}
              style={{
                minHeight: isEditMode ? '200px' : '150px',
              }}
            />
          </div>

          {/* Alert Messages */}
          {isEditMode && (
            <StatusAlert variant="warning" title="Warning">
              This will replace the current encrypted value. This action cannot be undone.
            </StatusAlert>
          )}

          {!isEditMode && (
            <StatusAlert variant="info" title="Info">
              This value is decrypted and displayed in plain text. Use the copy button to copy it to clipboard.
            </StatusAlert>
          )}

          {error && (
            <StatusAlert variant="error">
              {error}
            </StatusAlert>
          )}

          {hasUnsavedChanges && (
            <StatusAlert variant="warning" title="Unsaved changes">
              You have modified the vault value. Save or cancel to continue.
            </StatusAlert>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
