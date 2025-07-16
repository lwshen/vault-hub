import { Card } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { 
  Search,
  Plus,
  MoreVertical,
  Lock,
  Unlock,
  Users
} from 'lucide-react';

export default function VaultsContent() {
  const vaults = [
    { name: 'Production API Keys', status: 'locked', members: 5, secrets: 12, lastAccessed: '2 hours ago' },
    { name: 'Database Credentials', status: 'unlocked', members: 3, secrets: 8, lastAccessed: '1 day ago' },
    { name: 'SSL Certificates', status: 'locked', members: 2, secrets: 4, lastAccessed: '3 days ago' },
    { name: 'OAuth Tokens', status: 'locked', members: 4, secrets: 15, lastAccessed: '1 week ago' },
    { name: 'AWS Secrets', status: 'locked', members: 6, secrets: 20, lastAccessed: '2 days ago' },
    { name: 'Third-party APIs', status: 'unlocked', members: 3, secrets: 7, lastAccessed: '5 hours ago' }
  ];
  
  return (
    <>
      {/* Top Header */}
      <header className="bg-card border-b border-border p-6 flex-shrink-0">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold tracking-tight">Vaults</h1>
            <p className="text-muted-foreground">
              Manage and organize your secret vaults
            </p>
          </div>
          <div className="flex items-center gap-3">
            <Button variant="outline" size="sm">
              <Search className="h-4 w-4 mr-2" />
              Search
            </Button>
            <Button size="sm">
              <Plus className="h-4 w-4 mr-2" />
              New Vault
            </Button>
          </div>
        </div>
      </header>
  
      {/* Main Content */}
      <main className="flex-1 overflow-y-auto p-6">
        <div className="grid gap-4">
          {vaults.map((vault, index) => (
            <Card key={index} className="p-6">
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-4">
                  {vault.status === 'locked' ? (
                    <Lock className="h-5 w-5 text-red-500" />
                  ) : (
                    <Unlock className="h-5 w-5 text-green-500" />
                  )}
                  <div>
                    <h3 className="text-lg font-semibold">{vault.name}</h3>
                    <div className="flex items-center gap-4 text-sm text-muted-foreground">
                      <span>{vault.secrets} secrets</span>
                      <span className="flex items-center gap-1">
                        <Users className="h-3 w-3" />
                        {vault.members} members
                      </span>
                      <span>Last accessed {vault.lastAccessed}</span>
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
      </main>
    </>
  );
}
