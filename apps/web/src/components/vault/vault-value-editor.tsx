import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Label } from '@/components/ui/label';
import type { UseVaultActionsReturn } from '@/hooks/useVaultData';
import { AlertCircle, Info } from 'lucide-react';

interface VaultValueEditorProps {
  isEditMode: boolean;
  vaultActions: UseVaultActionsReturn;
  originalValue: string;
}

export function VaultValueEditor({
  isEditMode,
  vaultActions,
  originalValue,
}: VaultValueEditorProps) {
  const { error, setError, hasUnsavedChanges } = vaultActions;

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
              value={isEditMode ? vaultActions.editedValue : originalValue}
              onChange={(e) => {
                if (isEditMode) {
                  vaultActions.setEditedValue(e.target.value);
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
            <Alert variant="warning">
              <AlertCircle />
              <AlertTitle>Warning</AlertTitle>
              <AlertDescription>
                This will replace the current encrypted value. This action cannot be undone.
              </AlertDescription>
            </Alert>
          )}

          {!isEditMode && (
            <Alert variant="info">
              <Info />
              <AlertTitle>Info</AlertTitle>
              <AlertDescription>
                This value is decrypted and displayed in plain text. Use the copy button to copy it to clipboard.
              </AlertDescription>
            </Alert>
          )}

          {error && (
            <Alert variant="destructive">
              <AlertCircle />
              <AlertTitle>Error</AlertTitle>
              <AlertDescription>
                {error}
              </AlertDescription>
            </Alert>
          )}

          {hasUnsavedChanges && (
            <Alert variant="warning">
              <AlertCircle />
              <AlertTitle>Unsaved changes</AlertTitle>
              <AlertDescription>
                You have modified the vault value. Save or cancel to continue.
              </AlertDescription>
            </Alert>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
