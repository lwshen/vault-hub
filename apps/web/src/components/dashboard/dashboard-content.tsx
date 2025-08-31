import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import {
  Activity,
  Key,
  Lock,
  MoreVertical,
  Plus,
  Search,
  Unlock,
  Users,
  Vault,
} from 'lucide-react';
import { versionApi } from '@/apis/api';
import { useEffect, useState } from 'react';

export default function DashboardContent() {
  const [version, setVersion] = useState<{ version: string; commit: string; } | null>(null);

  useEffect(() => {
    const fetchVersion = async () => {
      try {
        const response = await versionApi.getVersion();
        setVersion(response);
      } catch (error) {
        console.error('Failed to fetch version:', error);
      }
    };
    fetchVersion();
  }, []);

  const stats = [
    {
      title: 'Total Vaults',
      value: '12',
      icon: Vault,
      change: '+2 this month',
      changeType: 'positive' as const,
    },
    {
      title: 'Active Users',
      value: '24',
      icon: Users,
      change: '+5 this week',
      changeType: 'positive' as const,
    },
    {
      title: 'Secrets Stored',
      value: '156',
      icon: Key,
      change: '+12 today',
      changeType: 'positive' as const,
    },
    {
      title: 'Recent Activity',
      value: '8',
      icon: Activity,
      change: 'Last 24 hours',
      changeType: 'neutral' as const,
    },
  ];

  const recentVaults = [
    { name: 'Production API Keys', status: 'locked', lastAccessed: '2 hours ago' },
    { name: 'Database Credentials', status: 'unlocked', lastAccessed: '1 day ago' },
    { name: 'SSL Certificates', status: 'locked', lastAccessed: '3 days ago' },
    { name: 'OAuth Tokens', status: 'locked', lastAccessed: '1 week ago' },
  ];

  return (
    <>
      {/* Top Header */}
      <header className="bg-card border-b border-border p-6 flex-shrink-0">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold tracking-tight">Dashboard</h1>
            <p className="text-muted-foreground">
              Manage your vaults and monitor activity
            </p>
          </div>
          <div className="flex items-center gap-3">
            <Button variant="outline" size="sm">
              <Search className="h-4 w-4 mr-2" />
              Search
            </Button>
            <Button variant="outline" size="sm">
              <MoreVertical className="h-4 w-4" />
            </Button>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="flex-1 overflow-y-auto p-6 space-y-6">
        {/* Stats Grid */}
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
          {stats.map((stat) => {
            const Icon = stat.icon;
            return (
              <Card key={stat.title} className="p-6">
                <div className="flex items-center justify-between space-y-0 pb-2">
                  <h3 className="text-sm font-medium text-muted-foreground">
                    {stat.title}
                  </h3>
                  <Icon className="h-4 w-4 text-muted-foreground" />
                </div>
                <div className="space-y-1">
                  <div className="text-2xl font-bold">{stat.value}</div>
                  <p className={`text-xs ${
                    stat.changeType === 'positive'
                      ? 'text-green-600'
                      : 'text-muted-foreground'
                  }`}>
                    {stat.change}
                  </p>
                </div>
              </Card>
            );
          })}
        </div>

        {/* Main Content Grid */}
        <div className="grid gap-6 lg:grid-cols-3">
          {/* Recent Vaults */}
          <Card className="lg:col-span-2 p-6">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-semibold">Recent Vaults</h2>
              <Button variant="ghost" size="sm">
                <MoreVertical className="h-4 w-4" />
              </Button>
            </div>
            <div className="space-y-3">
              {recentVaults.map((vault) => (
                <div key={vault.name} className="flex items-center justify-between p-3 rounded-lg border">
                  <div className="flex items-center gap-3">
                    {vault.status === 'locked' ? (
                      <Lock className="h-4 w-4 text-red-500" />
                    ) : (
                      <Unlock className="h-4 w-4 text-green-500" />
                    )}
                    <div>
                      <p className="font-medium">{vault.name}</p>
                      <p className="text-sm text-muted-foreground">
                        Last accessed {vault.lastAccessed}
                      </p>
                    </div>
                  </div>
                  <Button variant="outline" size="sm">
                    Access
                  </Button>
                </div>
              ))}
            </div>
          </Card>

          {/* System Status */}
          <Card className="p-6">
            <h2 className="text-lg font-semibold mb-4">System Status</h2>
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <span className="text-sm font-medium">API Status</span>
                <div className="flex items-center gap-2">
                  <div className="h-2 w-2 bg-green-500 rounded-full"></div>
                  <span className="text-sm text-muted-foreground">Online</span>
                </div>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm font-medium">Database</span>
                <div className="flex items-center gap-2">
                  <div className="h-2 w-2 bg-green-500 rounded-full"></div>
                  <span className="text-sm text-muted-foreground">Healthy</span>
                </div>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm font-medium">Backup</span>
                <div className="flex items-center gap-2">
                  <div className="h-2 w-2 bg-yellow-500 rounded-full"></div>
                  <span className="text-sm text-muted-foreground">Running</span>
                </div>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm font-medium">Storage</span>
                <span className="text-sm text-muted-foreground">78% Used</span>
              </div>
              {version && (
                <div className="pt-3 border-t border-border">
                  <div className="flex items-center justify-between">
                    <span className="text-sm font-medium">Version</span>
                    <span className="text-sm text-muted-foreground">{version.version}</span>
                  </div>
                  <div className="flex items-center justify-between mt-1">
                    <span className="text-sm font-medium">Commit</span>
                    <span className="text-sm text-muted-foreground font-mono">{version.commit.substring(0, 7)}</span>
                  </div>
                </div>
              )}
            </div>
          </Card>
        </div>

        {/* Recent Activity */}
        <Card className="p-6">
          <h2 className="text-lg font-semibold mb-4">Recent Activity</h2>
          <div className="space-y-3">
            <div className="flex items-center gap-3 p-3 rounded-lg bg-muted/50">
              <Activity className="h-4 w-4 text-blue-500" />
              <div className="flex-1">
                <p className="font-medium">Vault "Production API Keys" was accessed</p>
                <p className="text-sm text-muted-foreground">by john@example.com • 2 hours ago</p>
              </div>
            </div>
            <div className="flex items-center gap-3 p-3 rounded-lg bg-muted/50">
              <Plus className="h-4 w-4 text-green-500" />
              <div className="flex-1">
                <p className="font-medium">New vault "SSL Certificates" created</p>
                <p className="text-sm text-muted-foreground">by admin@example.com • 1 day ago</p>
              </div>
            </div>
            <div className="flex items-center gap-3 p-3 rounded-lg bg-muted/50">
              <Users className="h-4 w-4 text-purple-500" />
              <div className="flex-1">
                <p className="font-medium">Team member invited</p>
                <p className="text-sm text-muted-foreground">sarah@example.com invited by admin@example.com • 2 days ago</p>
              </div>
            </div>
          </div>
        </Card>
      </main>
    </>
  );
}
