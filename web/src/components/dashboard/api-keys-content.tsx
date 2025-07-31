import { apiKeyApi } from '@/apis/api';
import CreateApiKeyModal from '@/components/modals/create-api-key-modal';
import EditApiKeyModal from '@/components/modals/edit-api-key-modal';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import DashboardHeader from '@/components/layout/dashboard-header';
import { useApiKeys } from '@/hooks/use-api-keys';
import type { APIKey } from '@lwshen/vault-hub-ts-fetch-client';
import { AlertCircle, Key, Loader2, Pencil, Plus, Trash2 } from 'lucide-react';
import { useState } from 'react';


export default function ApiKeysContent() {
  const { apiKeys, isLoading, error, refetch } = useApiKeys();
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [selectedApiKey, setSelectedApiKey] = useState<APIKey | null>(null);

  const handleKeyCreated = () => {
    refetch();
  };

  const renderContent = () => {
    if (error) {
      return (
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
      );
    }

    return (
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
              <Card key={key.id} className="p-6">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-4">
                    <Key className="h-5 w-5 text-blue-500" />
                    <div>
                      <h3 className="text-lg font-semibold flex items-center gap-2">
                        {key.name}
                        {key.isActive === false && (
                          <span className="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-red-50 text-red-700 dark:bg-red-900/20 dark:text-red-300">
                            Disabled
                          </span>
                        )}
                      </h3>
                      <div className="flex items-center gap-4 text-sm text-muted-foreground flex-wrap">
                        <span>Created {key.createdAt.toLocaleString()}</span>
                        {key.expiresAt && <span>Expires {key.expiresAt.toLocaleString()}</span>}
                        {key.lastUsedAt && <span>Last Used {key.lastUsedAt.toLocaleString()}</span>}
                      </div>
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
        title="API Keys"
        description="Manage and create API keys"
        actions={
          <Button size="sm" onClick={() => setIsCreateModalOpen(true)}>
            <Plus className="h-4 w-4 mr-2" />
            New Key
          </Button>
        }
      />
      {renderContent()}

      <CreateApiKeyModal
        open={isCreateModalOpen}
        onOpenChange={setIsCreateModalOpen}
        onApiKeyCreated={handleKeyCreated}
      />

      {selectedApiKey && (
        <EditApiKeyModal
          open={isEditModalOpen}
          onOpenChange={setIsEditModalOpen}
          apiKey={selectedApiKey}
          onApiKeyUpdated={handleKeyCreated}
        />
      )}
    </>
  );
}
