import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { vaultApi } from '@/apis/api';
import { Loader2, X } from 'lucide-react';

interface CreateVaultModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onVaultCreated?: () => void;
}

interface FormData {
  name: string;
  value: string;
  description: string;
  category: string;
}

export default function CreateVaultModal({ open, onOpenChange, onVaultCreated }: CreateVaultModalProps) {
  const [formData, setFormData] = useState<FormData>({
    name: '',
    value: '',
    description: '',
    category: '',
  });
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleInputChange = (field: keyof FormData, value: string) => {
    setFormData(prev => ({ ...prev, [field]: value }));
    if (error) setError(null); // Clear error when user starts typing
  };

  const validateForm = (data: FormData): string | null => {
    const validations = [
      { check: !data.name.trim(), message: 'Name is required' },
      { check: !data.value.trim(), message: 'Value is required' },
      { check: data.name.length > 100, message: 'Name must be ≤100 characters' },
      { check: data.description?.length > 500, message: 'Description must be ≤500 characters' },
    ];

    const failed = validations.find(v => v.check);
    return failed?.message || null;
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    const validationError = validateForm(formData);
    if (validationError) {
      setError(validationError);
      return;
    }

    setIsLoading(true);
    setError(null);

    try {
      await vaultApi.createVault({
        name: formData.name.trim(),
        value: formData.value.trim(),
        description: formData.description.trim() || undefined,
        category: formData.category.trim() || undefined,
      });

      // Reset form
      setFormData({
        name: '',
        value: '',
        description: '',
        category: '',
      });

      // Close modal and notify parent
      onOpenChange(false);
      onVaultCreated?.();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create vault');
    } finally {
      setIsLoading(false);
    }
  };

  const handleClose = () => {
    if (!isLoading) {
      onOpenChange(false);
      setError(null);
    }
  };

  if (!open) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      {/* Backdrop */}
      <div
        className="fixed inset-0 bg-black/50"
        onClick={handleClose}
      />

      {/* Modal Content */}
      <Card className="relative z-10 w-full max-w-md mx-4">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-xl">Create New Vault</CardTitle>
          <Button
            variant="ghost"
            size="icon"
            onClick={handleClose}
            disabled={isLoading}
          >
            <X className="h-4 w-4" />
          </Button>
        </CardHeader>

        <CardContent>
          <p className="text-muted-foreground text-sm mb-6">
            Create a new vault to securely store your secrets and sensitive data.
          </p>

          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="name">Name *</Label>
              <Input
                id="name"
                placeholder="e.g., Production API Keys"
                value={formData.name}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => handleInputChange('name', e.target.value)}
                disabled={isLoading}
                required
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="value">Value *</Label>
              <textarea
                id="value"
                placeholder="Enter the secret value to be encrypted and stored"
                value={formData.value}
                onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) => handleInputChange('value', e.target.value)}
                disabled={isLoading}
                required
                rows={3}
                className="flex w-full rounded-md border border-input bg-transparent px-3 py-2 text-sm shadow-xs placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="category">Category</Label>
              <Input
                id="category"
                placeholder="e.g., API Keys, Database, Certificates"
                value={formData.category}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => handleInputChange('category', e.target.value)}
                disabled={isLoading}
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="description">Description</Label>
              <textarea
                id="description"
                placeholder="Optional description for this vault"
                value={formData.description}
                onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) => handleInputChange('description', e.target.value)}
                disabled={isLoading}
                rows={2}
                className="flex w-full rounded-md border border-input bg-transparent px-3 py-2 text-sm shadow-xs placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
              />
            </div>

            {error && (
              <div className="p-3 text-sm text-red-600 bg-red-50 border border-red-200 rounded-md dark:text-red-400 dark:bg-red-900/20 dark:border-red-800">
                {error}
              </div>
            )}

            <div className="flex justify-end space-x-2 pt-4">
              <Button
                type="button"
                variant="outline"
                onClick={handleClose}
                disabled={isLoading}
              >
                Cancel
              </Button>
              <Button type="submit" disabled={isLoading}>
                {isLoading ? (
                  <>
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    Creating...
                  </>
                ) : (
                  'Create Vault'
                )}
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
