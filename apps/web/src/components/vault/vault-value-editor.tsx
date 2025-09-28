import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Label } from '@/components/ui/label';
import type { UseVaultActionsReturn } from '@/hooks/useVaultData';
import { AlertCircle, Info } from 'lucide-react';
import { useMemo } from 'react';

interface VaultValueEditorProps {
  isEditMode: boolean;
  vaultActions: UseVaultActionsReturn;
  originalValue: string;
}

// Utility functions for content analysis
function analyzeContent(content: string) {
  const lines = content.split('\n');
  const characterCount = content.length;
  const lineCount = lines.length;

  // Determine if content is "long" based on multiple criteria
  const isLongByCharacters = characterCount > 500;
  const isLongByLines = lineCount > 10;
  const hasVeryLongLines = lines.some(line => line.length > 100);

  const isLong = isLongByCharacters || isLongByLines || hasVeryLongLines;

  return {
    characterCount,
    lineCount,
    isLong,
    isLongByCharacters,
    isLongByLines,
    hasVeryLongLines,
  };
}

function calculateTextareaProps(content: string, isEditMode: boolean) {
  const analysis = analyzeContent(content);

  // Base configuration for short content
  let rows = 6;
  let minHeight = isEditMode ? '200px' : '150px';
  let maxHeight = '600px';

  if (analysis.isLong) {
    // For long content, provide significantly more space
    const baseRows = Math.min(Math.max(analysis.lineCount + 2, 12), 25);
    rows = baseRows;

    // Responsive height based on content and mode
    if (analysis.isLongByCharacters && analysis.characterCount > 2000) {
      // Very long content gets even more space
      minHeight = isEditMode ? '500px' : '400px';
      maxHeight = '85vh';
    } else {
      // Regular long content
      minHeight = isEditMode ? '400px' : '300px';
      maxHeight = '75vh';
    }
  }

  return {
    rows,
    minHeight,
    maxHeight,
    analysis,
  };
}

export function VaultValueEditor({
  isEditMode,
  vaultActions,
  originalValue,
}: VaultValueEditorProps) {
  const { error, setError, hasUnsavedChanges } = vaultActions;

  // Calculate dynamic textarea properties based on content
  const currentValue = isEditMode ? vaultActions.editedValue : originalValue;
  const textareaProps = useMemo(() => calculateTextareaProps(currentValue, isEditMode), [currentValue, isEditMode]);

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
            <div className="flex items-center justify-between">
              <Label htmlFor="vault-value">
                {isEditMode ? 'Encrypted Value *' : 'Decrypted Value'}
              </Label>
              {textareaProps.analysis.isLong && (
                <span className="text-xs text-muted-foreground">
                  {textareaProps.analysis.characterCount} chars, {textareaProps.analysis.lineCount} lines
                </span>
              )}
            </div>
            <textarea
              id="vault-value"
              value={currentValue}
              onChange={(e) => {
                if (isEditMode) {
                  vaultActions.setEditedValue(e.target.value);
                  if (error) setError(null);
                }
              }}
              placeholder={isEditMode ? 'Enter the secret value to be encrypted and stored' : ''}
              readOnly={!isEditMode}
              rows={textareaProps.rows}
              className={`
                flex w-full rounded-md border border-input px-3 py-2 text-sm shadow-xs transition-all duration-200
                ${textareaProps.analysis.isLong ? 'resize-y' : 'resize-none'}
                ${isEditMode
      ? 'bg-transparent placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring'
      : 'bg-muted cursor-text select-all'
    }
                ${textareaProps.analysis.isLong ? 'font-mono text-xs sm:text-sm' : ''}
              `}
              style={{
                minHeight: textareaProps.minHeight,
                maxHeight: textareaProps.maxHeight,
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
