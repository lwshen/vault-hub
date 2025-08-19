import DashboardHeader from '@/components/layout/dashboard-header';
import CreateVaultModal from '@/components/modals/create-vault-modal';
import EditVaultModal from '@/components/modals/edit-vault-modal';
import EditVaultValueModal from '@/components/modals/edit-vault-value-modal';
import ViewVaultValueModal from '@/components/modals/view-vault-value-modal';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { useVaults } from '@/hooks/use-vaults';
import { vaultApi } from '@/apis/api';
import type { VaultLite } from '@lwshen/vault-hub-ts-fetch-client';
import {
  AlertCircle,
  Edit,
  Eye,
  Key,
  Loader2,
  Lock,
  MoreVertical,
  Plus,
  Trash2,
} from 'lucide-react';
import { useState } from 'react';


export default function VaultsContent() {
  const { vaults, isLoading, error, refetch } = useVaults();
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [isEditValueModalOpen, setIsEditValueModalOpen] = useState(false);
  const [isViewValueModalOpen, setIsViewValueModalOpen] = useState(false);
  const [deleteConfirmOpen, setDeleteConfirmOpen] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);
  const [selectedVault, setSelectedVault] = useState<VaultLite | null>(null);

  const handleVaultCreated = () => {
    refetch(); // Refresh the vault list after creation
  };

  const handleEditVault = (vault: VaultLite) => {
    setSelectedVault(vault);
    setIsEditModalOpen(true);
  };

  const handleEditVaultValue = (vault: VaultLite) => {
    setSelectedVault(vault);
    setIsEditValueModalOpen(true);
  };

  const handleViewVaultValue = (vault: VaultLite) => {
    setSelectedVault(vault);
    setIsViewValueModalOpen(true);
  };

  const handleDeleteVault = (vault: VaultLite) => {
    setSelectedVault(vault);
    setDeleteConfirmOpen(true);
  };

  const confirmDeleteVault = async () => {
    if (!selectedVault) return;

    setIsDeleting(true);
    try {
      await vaultApi.deleteVault(selectedVault.uniqueId);
      setDeleteConfirmOpen(false);
      refetch(); // Refresh the vault list after deletion
    } catch (err) {
      console.error('Failed to delete vault:', err);
      // TODO: Show error toast or notification
    } finally {
      setIsDeleting(false);
    }
  };

  const handleVaultUpdated = () => {
    refetch(); // Refresh the vault list after update
  };

  const renderContent = () => {
    if (error) {
      return (
        <main className="flex-1 overflow-y-auto p-6">
          <Card className="p-6">
            <div className="flex items-center justify-center min-h-[200px] flex-col gap-4">
              <AlertCircle className="h-12 w-12 text-red-500" />
              <div className="text-center">
                <h3 className="text-lg font-semibold">Failed to load vaults</h3>
                <p className="text-muted-foreground mb-4">{error}</p>
                <Button onClick={refetch}>Try Again</Button>
              </div>
            </div>
          </Card>
        </main>
      );
    }

    return (
      <main className="flex-1 overflow-y-auto p-6">
        {isLoading ? (
          <Card className="p-6">
            <div className="flex items-center justify-center min-h-[200px] flex-col gap-4">
              <Loader2 className="h-8 w-8 animate-spin text-primary" />
              <p className="text-muted-foreground">Loading vaults...</p>
            </div>
          </Card>
        ) : vaults.length === 0 ? (
          <Card className="p-6">
            <div className="flex items-center justify-center min-h-[200px] flex-col gap-4">
              <Lock className="h-12 w-12 text-muted-foreground" />
              <div className="text-center">
                <h3 className="text-lg font-semibold">No vaults found</h3>
                <p className="text-muted-foreground mb-4">Create your first vault to get started</p>
                <Button onClick={() => setIsCreateModalOpen(true)}>
                  <Plus className="h-4 w-4 mr-2" />
                  New Vault
                </Button>
              </div>
            </div>
          </Card>
        ) : (
          <div className="grid gap-4">
            {vaults.map((vault) => (
              <Card key={vault.uniqueId} className="p-6">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-4">
                    {/* Show lock icon - could be based on category or other logic */}
                    <Lock className="h-5 w-5 text-blue-500" />
                    <div>
                      <h3 className="text-lg font-semibold">{vault.name}</h3>
                      <div className="flex items-center gap-4 text-sm text-muted-foreground">
                        {vault.category && (
                          <span className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-blue-50 text-blue-700 dark:bg-blue-900/20 dark:text-blue-300">
                            {vault.category}
                          </span>
                        )}
                        {vault.description && <span>{vault.description}</span>}
                        {vault.updatedAt && (
                          <span>Last Updated {new Date(vault.updatedAt).toLocaleDateString()}</span>
                        )}
                      </div>
                    </div>
                  </div>
                  <div className="flex items-center gap-2">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => handleViewVaultValue(vault)}
                      title="View Value"
                    >
                      <Eye className="h-4 w-4" />
                    </Button>
                    <DropdownMenu>
                      <DropdownMenuTrigger asChild>
                        <Button variant="outline" size="sm">
                          <MoreVertical className="h-4 w-4" />
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="end">
                        <DropdownMenuItem onClick={() => handleEditVaultValue(vault)}>
                          <Key className="h-4 w-4" />
                          Edit Value
                        </DropdownMenuItem>
                        <DropdownMenuItem onClick={() => handleEditVault(vault)}>
                          <Edit className="h-4 w-4" />
                          Edit Properties
                        </DropdownMenuItem>
                        <DropdownMenuSeparator />
                        <DropdownMenuItem
                          onClick={() => handleDeleteVault(vault)}
                          variant="destructive"
                        >
                          <Trash2 className="h-4 w-4" />
                          Delete Vault
                        </DropdownMenuItem>
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </div>
                </div>
              </Card>
            ))}
          </div>
        )}
      </main>
    );
  };

  return (
    <>
      <DashboardHeader
        title="Vaults"
        description="Manage and organize your secret vaults"
        actions={
          <Button size="sm" onClick={() => setIsCreateModalOpen(true)}>
            <Plus className="h-4 w-4 mr-2" />
            New Vault
          </Button>
        }
      />
      {renderContent()}

      <CreateVaultModal
        open={isCreateModalOpen}
        onOpenChange={setIsCreateModalOpen}
        onVaultCreated={handleVaultCreated}
      />

      <EditVaultModal
        open={isEditModalOpen}
        onOpenChange={setIsEditModalOpen}
        vault={selectedVault}
        onVaultUpdated={handleVaultUpdated}
      />

      <EditVaultValueModal
        open={isEditValueModalOpen}
        onOpenChange={setIsEditValueModalOpen}
        vault={selectedVault}
        onVaultUpdated={handleVaultUpdated}
      />

      <ViewVaultValueModal
        open={isViewValueModalOpen}
        onOpenChange={setIsViewValueModalOpen}
        vault={selectedVault}
      />

      {/* Delete Confirmation Modal */}
      {deleteConfirmOpen && selectedVault && (
        <div className="fixed inset-0 z-50 flex items-center justify-center">
          <div className="fixed inset-0 bg-black/50" onClick={() => setDeleteConfirmOpen(false)} />
          <Card className="relative z-10 w-full max-w-md mx-4">
            <CardHeader>
              <CardTitle className="text-xl">Delete Vault</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-muted-foreground mb-4">
                Are you sure you want to delete the vault <span className="font-medium">{selectedVault.name}</span>?
                This action cannot be undone.
              </p>
              <div className="flex justify-end space-x-2">
                <Button
                  variant="outline"
                  onClick={() => setDeleteConfirmOpen(false)}
                  disabled={isDeleting}
                >
                  Cancel
                </Button>
                <Button
                  variant="destructive"
                  onClick={confirmDeleteVault}
                  disabled={isDeleting}
                >
                  {isDeleting ? (
                    <>
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                      Deleting...
                    </>
                  ) : (
                    'Delete'
                  )}
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>
      )}
    </>
  );
}
