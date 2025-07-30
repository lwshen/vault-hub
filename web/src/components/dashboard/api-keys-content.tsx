import { useState } from 'react';
import { Card } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Key, Plus, Loader2, AlertCircle, Pencil, Trash2 } from 'lucide-react';
import { useApiKeys } from '@/hooks/use-api-keys';
import type { APIKey } from '@lwshen/vault-hub-ts-fetch-client';
import { CreateApiKeyModal } from '@/components/modals/create-api-key-modal';
import { EditApiKeyModal } from '@/components/modals/edit-api-key-modal';

export default function ApiKeysContent() {
  const { apiKeys, isLoading, error, refetch } = useApiKeys();
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [selectedApiKey, setSelectedApiKey] = useState<APIKey | null>(null);

  const handleKeyCreated = () => {
    refetch();
  };

  if (error) {
    return (
      <>
        {/* Top Header */}
        <header className="bg-card border-b border-border p-6 flex-shrink-0">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold tracking-tight">API Keys</h1>
              <p className="text-muted-foreground">Manage and create API keys for programmatic access</p>
            </div>
            <div className="flex items-center gap-3">
              <Button size="sm" onClick={() => setIsCreateModalOpen(true)}>
                <Plus className="h-4 w-4 mr-2" />
                New Key
              </Button>
            </div>
          </div>
        </header>

        {/* Error State */}
        <main className="flex-1 overflow-y-auto p-6">
          <Card className="p-6">
            <div className="flex items-center justify-center min-h-[200px] flex-col gap-4">
              <AlertCircle className="h-12 w-12 text-red-500" />
              <div className="text-center">
                <h3 className="text-lg font-semibold">Failed to load API keys</h3>
                <p className="text-muted-foreground mb-4">{error}</p>
                <Button onClick={refetch}>Try Again</Button>
              </div>
            </div>
          </Card>
        </main>

        <CreateApiKeyModal
          open={isCreateModalOpen}
          onOpenChange={setIsCreateModalOpen}
          onApiKeyCreated={handleKeyCreated}
        />
      </>
    );
  }

  return (
    <>
      {/* Top Header */}
      <header className="bg-card border-b border-border p-6 flex-shrink-0">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold tracking-tight">API Keys</h1>
            <p className="text-muted-foreground">Manage and create API keys for programmatic access</p>
          </div>
          <div className="flex items-center gap-3">
            <Button size="sm" onClick={() => setIsCreateModalOpen(true)}>
              <Plus className="h-4 w-4 mr-2" />
              New Key
            </Button>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="flex-1 overflow-y-auto p-6">
        {isLoading ? (
          <Card className="p-6 flex items-center justify-center">
            <Loader2 className="h-8 w-8 animate-spin text-primary" />
          </Card>
        ) : apiKeys.length === 0 ? (
          <Card className="p-6 flex items-center justify-center flex-col gap-4">
            <Key className="h-12 w-12 text-muted-foreground" />
            <div className="text-center">
              <h3 className="text-lg font-semibold">No API keys</h3>
              <p className="text-muted-foreground mb-4">Create your first key to get started</p>
              <Button size="sm" onClick={() => setIsCreateModalOpen(true)}>
                <Plus className="h-4 w-4 mr-2" />
                New Key
              </Button>
            </div>
          </Card>
        ) : (
          <div className="grid gap-4">
            {apiKeys.map((key) => (
              <Card key={key.id} className="p-6 flex items-center justify-between">
                <div>
                  <p className="font-medium flex items-center gap-2">
                    <Key className="h-4 w-4 text-blue-500" /> {key.name}
                    {key.isActive === false && (
                      <span className="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-red-50 text-red-700 dark:bg-red-900/20 dark:text-red-300">
                        Disabled
                      </span>
                    )}
                  </p>
                  <div className="flex items-center gap-4 text-sm text-muted-foreground flex-wrap">
                    <span>Created {new Date(key.createdAt as any).toLocaleDateString()}</span>
                    {key.expiresAt && <span>Expires {new Date(key.expiresAt as any).toLocaleDateString()}</span>}
                    {key.lastUsedAt && <span>Last Used {new Date(key.lastUsedAt as any).toLocaleDateString()}</span>}
                  </div>
                </div>
                <div className="flex items-center gap-2">
                  <Button
                    variant="outline"
                    size="icon"
                    onClick={() => {
                      setSelectedApiKey(key);
                      setIsEditModalOpen(true);
                    }}
                  >
                    <Pencil className="h-4 w-4" />
                  </Button>
                  <Button
                    variant="outline"
                    size="icon"
                    onClick={async () => {
                      const confirmed = confirm('Delete this API key? This action cannot be undone.');
                      if (!confirmed) return;
                      try {
                        await apiKeyApi.deleteAPIKey(key.id);
                        refetch();
                      } catch (err) {
                        alert(err instanceof Error ? err.message : 'Failed to delete API key');
                      }
                    }}
                  >
                    <Trash2 className="h-4 w-4 text-red-600" />
                  </Button>
                </div>
              </Card>
            ))}
          </div>
        )}
      </main>

      <CreateApiKeyModal
        open={isCreateModalOpen}
        onOpenChange={setIsCreateModalOpen}
        onApiKeyCreated={handleKeyCreated}
      />

     <EditApiKeyModal
        open={isEditModalOpen}
        onOpenChange={setIsEditModalOpen}
        apiKey={selectedApiKey}
        onApiKeyUpdated={handleKeyCreated}
      />
    </>
  );
}