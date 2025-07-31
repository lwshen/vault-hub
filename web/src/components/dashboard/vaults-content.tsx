import { useState } from 'react';
import { Card } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import {
  Plus,
  MoreVertical,
  Lock,
  Loader2,
  AlertCircle,
} from 'lucide-react';
import { useVaults } from '@/hooks/use-vaults';
import CreateVaultModal from '@/components/modals/create-vault-modal';

const VaultsHeader = ({ onCreateClick }: { onCreateClick: () => void }) => (
  <header className="bg-card border-b border-border p-6 flex-shrink-0">
    <div className="flex items-center justify-between">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Vaults</h1>
        <p className="text-muted-foreground">
          Manage and organize your secret vaults
        </p>
      </div>
      <div className="flex items-center gap-3">
        <Button size="sm" onClick={onCreateClick}>
          <Plus className="h-4 w-4 mr-2" />
          New Vault
        </Button>
      </div>
    </div>
  </header>
);

export default function VaultsContent() {
  const { vaults, isLoading, error, refetch } = useVaults();
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);

  const handleVaultCreated = () => {
    refetch(); // Refresh the vault list after creation
  };

  if (error) {
    return (
      <>
        <VaultsHeader onCreateClick={() => setIsCreateModalOpen(true)} />
        {/* Error State */}
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

        <CreateVaultModal
          open={isCreateModalOpen}
          onOpenChange={setIsCreateModalOpen}
          onVaultCreated={handleVaultCreated}
        />
      </>
    );
  }

  return (
    <>
      <VaultsHeader onCreateClick={() => setIsCreateModalOpen(true)} />
      {/* Main Content */}
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
                    <Button variant="outline" size="sm">
                      View
                    </Button>
                    <Button variant="outline" size="sm">
                      <MoreVertical className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
              </Card>
            ))}
          </div>
        )}
      </main>

      <CreateVaultModal
        open={isCreateModalOpen}
        onOpenChange={setIsCreateModalOpen}
        onVaultCreated={handleVaultCreated}
      />
    </>
  );
}
