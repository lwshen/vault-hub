import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Label } from '@/components/ui/label';
import { Edit3, Eye } from 'lucide-react';
import type { Vault } from '@lwshen/vault-hub-ts-fetch-client';

interface VaultMetadataProps {
  vault: Vault;
  isEditMode: boolean;
}

export function VaultMetadata({ vault, isEditMode }: VaultMetadataProps) {
  return (
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
              <Label className="text-xs uppercase tracking-wide text-muted-foreground">
                Category
              </Label>
              <p className="mt-1 text-sm font-medium">{vault.category}</p>
            </div>
          )}
          {vault.description && (
            <div className="sm:col-span-2">
              <Label className="text-xs uppercase tracking-wide text-muted-foreground">
                Description
              </Label>
              <p className="mt-1 text-sm">{vault.description}</p>
            </div>
          )}
          {vault.updatedAt && (
            <div>
              <Label className="text-xs uppercase tracking-wide text-muted-foreground">
                Last Updated
              </Label>
              <p className="mt-1 text-sm">
                {new Date(vault.updatedAt).toLocaleDateString()}
              </p>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
