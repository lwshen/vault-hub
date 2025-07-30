import { useState } from 'react';
import { Card } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Key, Plus, Loader2 } from 'lucide-react';

export default function ApiKeysContent() {
  const [isLoading] = useState(false);
  const [apiKeys] = useState<Array<{ id: number; name: string; createdAt: string }>>([]);

  return (
    <>
      {/* Top Header */}
      <header className="bg-card border-b border-border p-6 flex-shrink-0">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold tracking-tight">API Keys</h1>
            <p className="text-muted-foreground">
              Manage and create API keys for programmatic access
            </p>
          </div>
          <div className="flex items-center gap-3">
            <Button size="sm">
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
              <p className="text-muted-foreground mb-4">
                Create your first key to get started
              </p>
              <Button size="sm">
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
                  <p className="font-medium">{key.name}</p>
                  <p className="text-sm text-muted-foreground">
                    Created {new Date(key.createdAt).toLocaleDateString()}
                  </p>
                </div>
                <Button variant="outline" size="sm">
                  Revoke
                </Button>
              </Card>
            ))}
          </div>
        )}
      </main>
    </>
  );
}