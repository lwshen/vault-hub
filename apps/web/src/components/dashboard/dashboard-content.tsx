import { auditApi, versionApi } from '@/apis/api';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import type { AuditLog, AuditMetricsResponse } from '@lwshen/vault-hub-ts-fetch-client';
import { AuditLogActionEnum } from '@lwshen/vault-hub-ts-fetch-client';
import {
  Activity,
  Key,
  Loader2,
  Lock,
  MoreVertical,
  Plus,
  Unlock,
  Users,
  Vault,
  Eye,
  Edit,
  Trash2,
} from 'lucide-react';
import { useEffect, useState } from 'react';

export default function DashboardContent() {
  const [version, setVersion] = useState<{ version: string; commit: string; } | null>(null);
  const [metrics, setMetrics] = useState<AuditMetricsResponse | null>(null);
  const [recentAuditLogs, setRecentAuditLogs] = useState<AuditLog[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const versionResponse = await versionApi.getVersion();
        setVersion(versionResponse);
        const metricsResponse = await auditApi.getAuditMetrics();
        setMetrics(metricsResponse);
        // Fetch recent audit logs (first 5 items from page 1)
        const auditResponse = await auditApi.getAuditLogs(5, 1);
        setRecentAuditLogs(auditResponse.auditLogs || []);
      } catch (error) {
        console.error('Failed to fetch data:', error);
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, []);

  // Icon mapping for different audit actions - using correct enum values
  const getIconForAction = (action: AuditLogActionEnum) => {
    const iconMap: { [key in AuditLogActionEnum]: { icon: typeof Lock; color: string; }; } = {
      [AuditLogActionEnum.ReadVault]: { icon: Eye, color: 'text-blue-500' },
      [AuditLogActionEnum.CreateVault]: { icon: Plus, color: 'text-green-500' },
      [AuditLogActionEnum.UpdateVault]: { icon: Edit, color: 'text-yellow-500' },
      [AuditLogActionEnum.DeleteVault]: { icon: Trash2, color: 'text-red-500' },
      [AuditLogActionEnum.LoginUser]: { icon: Users, color: 'text-purple-500' },
      [AuditLogActionEnum.LogoutUser]: { icon: Users, color: 'text-gray-500' },
      [AuditLogActionEnum.RegisterUser]: { icon: Users, color: 'text-purple-500' },
      [AuditLogActionEnum.CreateApiKey]: { icon: Key, color: 'text-green-500' },
      [AuditLogActionEnum.UpdateApiKey]: { icon: Key, color: 'text-yellow-500' },
      [AuditLogActionEnum.DeleteApiKey]: { icon: Key, color: 'text-red-500' },
    };

    return iconMap[action] || { icon: Activity, color: 'text-gray-500' };
  };

  // Convert action to readable title
  const getActionTitle = (action: AuditLogActionEnum) => {
    const titleMap: { [key in AuditLogActionEnum]: string; } = {
      [AuditLogActionEnum.ReadVault]: 'Vault accessed',
      [AuditLogActionEnum.CreateVault]: 'Vault created',
      [AuditLogActionEnum.UpdateVault]: 'Vault updated',
      [AuditLogActionEnum.DeleteVault]: 'Vault deleted',
      [AuditLogActionEnum.LoginUser]: 'User logged in',
      [AuditLogActionEnum.LogoutUser]: 'User logged out',
      [AuditLogActionEnum.RegisterUser]: 'User registered',
      [AuditLogActionEnum.CreateApiKey]: 'API key created',
      [AuditLogActionEnum.UpdateApiKey]: 'API key updated',
      [AuditLogActionEnum.DeleteApiKey]: 'API key deleted',
    };

    return titleMap[action] || action;
  };

  const formatTimestamp = (timestamp: string | Date) => {
    const date = new Date(timestamp);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffMins = Math.floor(diffMs / 60000);
    const diffHours = Math.floor(diffMins / 60);
    const diffDays = Math.floor(diffHours / 24);

    if (diffMins < 1) return 'Just now';
    if (diffMins < 60) return `${diffMins} minute${diffMins === 1 ? '' : 's'} ago`;
    if (diffHours < 24) return `${diffHours} hour${diffHours === 1 ? '' : 's'} ago`;
    if (diffDays < 7) return `${diffDays} day${diffDays === 1 ? '' : 's'} ago`;
    return date.toLocaleDateString();
  };

  const stats = [
    {
      title: 'Total Events (30 days)',
      value: metrics?.totalEventsLast30Days?.toString() || '-',
      icon: Activity,
      change: 'Last 30 days',
      changeType: 'neutral' as const,
    },
    {
      title: 'Events (24 hours)',
      value: metrics?.eventsCountLast24Hours?.toString() || '-',
      icon: Users,
      change: 'Last 24 hours',
      changeType: 'positive' as const,
    },
    {
      title: 'Vault Events (30 days)',
      value: metrics?.vaultEventsLast30Days?.toString() || '-',
      icon: Vault,
      change: 'Last 30 days',
      changeType: 'positive' as const,
    },
    {
      title: 'API Key Events (30 days)',
      value: metrics?.apiKeyEventsLast30Days?.toString() || '-',
      icon: Key,
      change: 'Last 30 days',
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
                  {loading ? (
                    <Loader2 className="h-4 w-4 animate-spin inline" />
                  ) : (
                    <div className="text-2xl font-bold">{stat.value}</div>
                  )}
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

        {/* Recent Audit Logs */}
        <Card className="p-6">
          <h2 className="text-lg font-semibold mb-4">Recent Audit Logs</h2>
          <div className="space-y-3">
            {loading ? (
              <div className="flex items-center justify-center py-8">
                <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
                <span className="ml-2 text-sm text-muted-foreground">Loading recent audit logs...</span>
              </div>
            ) : recentAuditLogs.length > 0 ? (
              recentAuditLogs.map((log) => {
                const { icon: ActionIcon, color } = getIconForAction(log.action);
                const actionTitle = getActionTitle(log.action);
                const resourceName = log.vault?.name || log.apiKey?.name;
                const uniqueKey = `${log.action}-${log.createdAt}-${log.vault?.uniqueId || log.apiKey?.id || 'user'}`;
                return (
                  <div key={uniqueKey} className="flex items-center gap-3 p-3 rounded-lg bg-muted/50">
                    <ActionIcon className={`h-4 w-4 ${color}`} />
                    <div className="flex-1">
                      <p className="font-medium">
                        {actionTitle}{resourceName && ` (${resourceName})`}
                      </p>
                      <p className="text-sm text-muted-foreground">
                        {log.apiKey ? 'via API Key' : 'via Web UI'} • {formatTimestamp(log.createdAt)}
                      </p>
                    </div>
                  </div>
                );
              })
            ) : (
              <div className="flex items-center justify-center py-8 text-center">
                <div>
                  <Activity className="h-8 w-8 text-muted-foreground mx-auto mb-2" />
                  <p className="text-sm text-muted-foreground">No recent audit logs</p>
                </div>
              </div>
            )}
          </div>
        </Card>
      </main>
    </>
  );
}
