import { useState, useEffect, useCallback } from 'react';
import { vaultApi } from '@/apis/api';
import type { Vault } from '@lwshen/vault-hub-ts-fetch-client';
import { toast } from 'sonner';

export interface UseVaultDataReturn {
  vault: Vault | null;
  originalValue: string;
  isLoading: boolean;
  error: string | null;
  refetch: () => void;
  updateLocalVault: (updatedVault: Vault) => void;
}

export function useVaultData(vaultId: string): UseVaultDataReturn {
  const [vault, setVault] = useState<Vault | null>(null);
  const [originalValue, setOriginalValue] = useState('');
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const refetch = useCallback(async () => {
    if (!vaultId) return;

    setIsLoading(true);
    setError(null);

    try {
      const fullVault = await vaultApi.getVault(vaultId);
      setVault(fullVault);
      setOriginalValue(fullVault.value || '');
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to fetch vault';
      setError(errorMessage);
    } finally {
      setIsLoading(false);
    }
  }, [vaultId]);

  useEffect(() => {
    refetch();
  }, [refetch]);

  const updateLocalVault = (updatedVault: Vault) => {
    setVault(updatedVault);
    setOriginalValue(updatedVault.value || '');
  };

  return {
    vault,
    originalValue,
    isLoading,
    error,
    refetch,
    updateLocalVault,
  };
}

export interface UseVaultActionsOptions {
  vault: Vault | null;
  originalValue: string;
  onSaveSuccess?: (updatedVault: Vault) => void;
}

export interface UseVaultActionsReturn {
  editedValue: string;
  setEditedValue: (value: string) => void;
  isSaving: boolean;
  error: string | null;
  setError: (error: string | null) => void;
  hasUnsavedChanges: boolean;
  handleSave: () => Promise<boolean>;
  handleCopy: () => Promise<void>;
  resetChanges: () => void;
}

export function useVaultActions({
  vault,
  originalValue,
  onSaveSuccess,
}: UseVaultActionsOptions): UseVaultActionsReturn {
  const [editedValue, setEditedValue] = useState('');
  const [isSaving, setIsSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    // eslint-disable-next-line @eslint-react/hooks-extra/no-direct-set-state-in-use-effect
    setEditedValue(originalValue);
  }, [originalValue]);

  const hasUnsavedChanges = editedValue !== originalValue;

  const handleSave = async (): Promise<boolean> => {
    if (!vault) return false;

    if (!editedValue.trim()) {
      setError('Value is required');
      return false;
    }

    setIsSaving(true);
    setError(null);

    try {
      const updatedVault = await vaultApi.updateVault(vault.uniqueId, {
        value: editedValue.trim(),
      });

      toast.success('Vault value updated successfully');
      onSaveSuccess?.(updatedVault);
      return true;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to update vault value';
      setError(errorMessage);
      return false;
    } finally {
      setIsSaving(false);
    }
  };

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(originalValue);
      toast.success(`Copied value for "${vault?.name}" to clipboard`);
    } catch (err) {
      console.error('Failed to copy to clipboard:', err);
      toast.error('Failed to copy to clipboard');
    }
  };

  const resetChanges = () => {
    setEditedValue(originalValue);
    setError(null);
  };

  return {
    editedValue,
    setEditedValue,
    isSaving,
    error,
    setError,
    hasUnsavedChanges,
    handleSave,
    handleCopy,
    resetChanges,
  };
}
