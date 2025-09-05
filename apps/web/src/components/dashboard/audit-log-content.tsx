import { auditApi } from '@/apis/api';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import DashboardHeader from '@/components/layout/dashboard-header';
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from '@/components/ui/pagination';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { AuditLogActionEnum, type AuditLog } from '@lwshen/vault-hub-ts-fetch-client';
import {
  Activity,
  AlertCircle,
  Edit,
  Key,
  Loader2,
  Lock,
  LogIn,
  LogOut,
  Plus,
  Trash2,
  UserPlus,
  Globe,
  TrendingUp,
  Shield,
} from 'lucide-react';
import { useCallback, useEffect, useState } from 'react';

// Icon mapping for different audit actions - using correct enum values
const getIconForAction = (action: AuditLogActionEnum) => {
  const iconMap: { [key in AuditLogActionEnum]: { icon: typeof Lock; color: string; }; } = {
    [AuditLogActionEnum.ReadVault]: { icon: Lock, color: 'text-blue-500' },
    [AuditLogActionEnum.CreateVault]: { icon: Plus, color: 'text-green-500' },
    [AuditLogActionEnum.UpdateVault]: { icon: Edit, color: 'text-yellow-500' },
    [AuditLogActionEnum.DeleteVault]: { icon: Trash2, color: 'text-red-500' },
    [AuditLogActionEnum.LoginUser]: { icon: LogIn, color: 'text-purple-500' },
    [AuditLogActionEnum.LogoutUser]: { icon: LogOut, color: 'text-gray-500' },
    [AuditLogActionEnum.RegisterUser]: { icon: UserPlus, color: 'text-purple-500' },
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

// Format timestamp to show precise date and time
const formatTimestamp = (timestamp: string | Date) => {
  const date = new Date(timestamp);
  const now = new Date();
  const diffInMs = now.getTime() - date.getTime();
  const diffInMinutes = Math.floor(diffInMs / (1000 * 60));
  const diffInHours = Math.floor(diffInMinutes / 60);
  const diffInDays = Math.floor(diffInHours / 24);

  // For very recent events, show relative time
  if (diffInMinutes < 1) return 'Just now';
  if (diffInMinutes < 60) return `${diffInMinutes} minute${diffInMinutes > 1 ? 's' : ''} ago`;

  // For older events, show precise date and time
  const timeString = date.toLocaleTimeString('en-US', {
    hour12: false,
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  });

  if (diffInHours < 24) {
    return `Today ${timeString}`;
  } else if (diffInDays === 1) {
    return `Yesterday ${timeString}`;
  } else {
    const dateString = date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
    });
    return `${dateString} ${timeString}`;
  }
};


export default function AuditLogContent() {
  const [auditLogs, setAuditLogs] = useState<AuditLog[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalCount, setTotalCount] = useState(0);
  const [totalPages, setTotalPages] = useState(0);
  const [pageSize, setPageSize] = useState(20);

  const fetchAuditLogs = useCallback(async (page: number) => {
    try {
      setIsLoading(true);
      setError(null);

      const response = await auditApi.getAuditLogs(pageSize, page);
      const newLogs = response.auditLogs || [];

      setAuditLogs(newLogs);
      setTotalCount(response.totalCount || 0);
      setTotalPages(Math.ceil((response.totalCount || 0) / pageSize));
      setCurrentPage(page);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch audit logs');
    } finally {
      setIsLoading(false);
    }
  }, [pageSize]);


  useEffect(() => {
    fetchAuditLogs(1);
  }, [fetchAuditLogs]);

  const handlePageChange = (page: number) => {
    if (page >= 1 && page <= totalPages) {
      fetchAuditLogs(page);
    }
  };

  const handlePageSizeChange = (newPageSize: string) => {
    const size = parseInt(newPageSize);
    setPageSize(size);
    setCurrentPage(1);
  };

  // Update fetchAuditLogs when pageSize changes
  useEffect(() => {
    fetchAuditLogs(1);
  }, [pageSize, fetchAuditLogs]);

  const getAuditTypeLabel = (action: AuditLogActionEnum) => {
    const labels: { [key in AuditLogActionEnum]: string; } = {
      [AuditLogActionEnum.ReadVault]: 'Vault Operation',
      [AuditLogActionEnum.CreateVault]: 'Vault Operation',
      [AuditLogActionEnum.UpdateVault]: 'Vault Operation',
      [AuditLogActionEnum.DeleteVault]: 'Vault Operation',
      [AuditLogActionEnum.LoginUser]: 'Authentication',
      [AuditLogActionEnum.LogoutUser]: 'Authentication',
      [AuditLogActionEnum.RegisterUser]: 'Authentication',
      [AuditLogActionEnum.CreateApiKey]: 'API Key Management',
      [AuditLogActionEnum.UpdateApiKey]: 'API Key Management',
      [AuditLogActionEnum.DeleteApiKey]: 'API Key Management',
    };
    return labels[action] || 'System';
  };

  const renderContent = () => {
    if (error) {
      return (
        <main className="flex-1 overflow-y-auto p-6">
          <Card className="p-6">
            <div className="flex items-center justify-center min-h-[200px] flex-col gap-4">
              <AlertCircle className="h-12 w-12 text-red-500" />
              <div className="text-center">
                <h3 className="text-lg font-semibold">Failed to load audit logs</h3>
                <p className="text-muted-foreground mb-4">{error}</p>
                <Button onClick={() => fetchAuditLogs(currentPage)}>Try Again</Button>
              </div>
            </div>
          </Card>
        </main>
      );
    }

    return (
      <main className="flex-1 overflow-y-auto p-6">
        <div className="space-y-4">
          {/* Audit Stats */}
          <div className="grid gap-4 md:grid-cols-4 mb-6">
            {/* Row 1 - Primary Metrics */}
            <Card className="p-4">
              <div className="flex items-center gap-3">
                <Activity className="h-8 w-8 text-blue-500" />
                <div>
                  <p className="text-2xl font-bold">1,247</p>
                  <p className="text-sm text-muted-foreground">Total Events</p>
                  <p className="text-xs text-muted-foreground">Last 30 days</p>
                </div>
              </div>
            </Card>
            <Card className="p-4">
              <div className="flex items-center gap-3">
                <Activity className="h-8 w-8 text-orange-500" />
                <div>
                  <p className="text-2xl font-bold">43</p>
                  <p className="text-sm text-muted-foreground">Last 24 Hours</p>
                  <p className="text-xs text-muted-foreground">Recent activity</p>
                </div>
              </div>
            </Card>
            <Card className="p-4">
              <div className="flex items-center gap-3">
                <Lock className="h-8 w-8 text-green-500" />
                <div>
                  <p className="text-2xl font-bold">892</p>
                  <p className="text-sm text-muted-foreground">Vault Events</p>
                  <p className="text-xs text-muted-foreground">Last 30 days</p>
                </div>
              </div>
            </Card>
            <Card className="p-4">
              <div className="flex items-center gap-3">
                <Key className="h-8 w-8 text-cyan-500" />
                <div>
                  <p className="text-2xl font-bold">156</p>
                  <p className="text-sm text-muted-foreground">API Key Events</p>
                  <p className="text-xs text-muted-foreground">Last 30 days</p>
                </div>
              </div>
            </Card>
          </div>

          {/* Secondary Metrics */}
          <div className="grid gap-4 md:grid-cols-4 mb-6">
            <Card className="p-4">
              <div className="flex items-center gap-3">
                <Shield className="h-8 w-8 text-red-500" />
                <div>
                  <p className="text-2xl font-bold">12</p>
                  <p className="text-sm text-muted-foreground">Critical Actions</p>
                  <p className="text-xs text-muted-foreground">Delete operations</p>
                </div>
              </div>
            </Card>
            <Card className="p-4">
              <div className="flex items-center gap-3">
                <Globe className="h-8 w-8 text-purple-500" />
                <div>
                  <p className="text-2xl font-bold">23</p>
                  <p className="text-sm text-muted-foreground">Unique IPs</p>
                  <p className="text-xs text-muted-foreground">Security monitoring</p>
                </div>
              </div>
            </Card>
            <Card className="p-4">
              <div className="flex items-center gap-3">
                <TrendingUp className="h-8 w-8 text-indigo-500" />
                <div>
                  <p className="text-2xl font-bold">prod-secrets</p>
                  <p className="text-sm text-muted-foreground">Most Active Vault</p>
                  <p className="text-xs text-muted-foreground">324 accesses</p>
                </div>
              </div>
            </Card>
            <Card className="p-4">
              <div className="flex items-center gap-3">
                <Activity className="h-8 w-8 text-emerald-500" />
                <div>
                  <p className="text-2xl font-bold">7</p>
                  <p className="text-sm text-muted-foreground">Active Sessions</p>
                  <p className="text-xs text-muted-foreground">Current users</p>
                </div>
              </div>
            </Card>
          </div>

          {/* Audit List */}
          <Card className="p-6">
            <div className="mb-4">
              {/* Header with title */}
              <div className="flex items-center justify-between mb-3">
                <h3 className="text-lg font-semibold">Audit Logs</h3>
              </div>

              {/* Controls - responsive layout */}
              <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
                <div className="flex items-center gap-2">
                  <span className="text-sm text-muted-foreground">Show</span>
                  <Select value={pageSize.toString()} onValueChange={handlePageSizeChange}>
                    <SelectTrigger className="w-20">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="10">10</SelectItem>
                      <SelectItem value="20">20</SelectItem>
                      <SelectItem value="50">50</SelectItem>
                      <SelectItem value="100">100</SelectItem>
                    </SelectContent>
                  </Select>
                  <span className="text-sm text-muted-foreground">per page</span>
                </div>

                {totalCount > 0 && (
                  <p className="text-sm text-muted-foreground">
                    Showing {((currentPage - 1) * pageSize) + 1}-{Math.min(currentPage * pageSize, totalCount)} of {totalCount} events
                  </p>
                )}
              </div>
            </div>
            {isLoading ? (
              <div className="flex items-center justify-center min-h-[200px] flex-col gap-4">
                <Loader2 className="h-8 w-8 animate-spin text-primary" />
                <p className="text-muted-foreground">Loading audit logs...</p>
              </div>
            ) : auditLogs.length === 0 ? (
              <div className="flex items-center justify-center min-h-[200px] flex-col gap-4">
                <Activity className="h-12 w-12 text-muted-foreground" />
                <div className="text-center">
                  <h3 className="text-lg font-semibold">No audit logs found</h3>
                  <p className="text-muted-foreground">No activity has been recorded yet.</p>
                </div>
              </div>
            ) : (
              <div className="space-y-4">
                {auditLogs.map((audit) => {
                  const { icon: Icon, color } = getIconForAction(audit.action);
                  return (
                    <div key={`${audit.action}-${audit.createdAt}`} className="flex items-start gap-4 p-4 rounded-lg border bg-card hover:bg-muted/50 transition-colors">
                      <div className="flex-shrink-0">
                        <div className="w-10 h-10 rounded-full bg-muted flex items-center justify-center">
                          <Icon className={`h-5 w-5 ${color}`} />
                        </div>
                      </div>
                      <div className="flex-1 min-w-0">
                        <div className="flex items-start justify-between">
                          <div className="flex-1">
                            <h3 className="font-medium text-foreground">{getActionTitle(audit.action)}</h3>
                            <div className="flex items-center gap-4 mt-2 text-xs text-muted-foreground">
                              <span>{formatTimestamp(audit.createdAt)}</span>
                              {audit.vault && (
                                <>
                                  <span>•</span>
                                  <span>Vault: {audit.vault.name} ({audit.vault.uniqueId})</span>
                                </>
                              )}
                              {audit.apiKey && (
                                <>
                                  <span>•</span>
                                  <span>API Key: {audit.apiKey.name} (ID: {audit.apiKey.id})</span>
                                </>
                              )}
                              {audit.ipAddress && (
                                <>
                                  <span>•</span>
                                  <span>IP: {audit.ipAddress}</span>
                                </>
                              )}
                            </div>
                          </div>
                          <div className="flex-shrink-0">
                            <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-muted text-muted-foreground">
                              {getAuditTypeLabel(audit.action)}
                            </span>
                          </div>
                        </div>
                      </div>
                    </div>
                  );
                })}

                {/* Pagination */}
                {totalPages > 1 && (
                  <div className="mt-6 flex justify-center">
                    <Pagination>
                      <PaginationContent>
                        <PaginationItem>
                          <PaginationPrevious
                            onClick={() => handlePageChange(currentPage - 1)}
                            className={currentPage <= 1 ? 'pointer-events-none opacity-50' : 'cursor-pointer'}
                          />
                        </PaginationItem>

                        {Array.from({ length: Math.min(5, totalPages) }, (_, i) => {
                          let pageNum;
                          if (totalPages <= 5) {
                            pageNum = i + 1;
                          } else if (currentPage <= 3) {
                            pageNum = i + 1;
                          } else if (currentPage >= totalPages - 2) {
                            pageNum = totalPages - 4 + i;
                          } else {
                            pageNum = currentPage - 2 + i;
                          }

                          return (
                            <PaginationItem key={pageNum}>
                              <PaginationLink
                                onClick={() => handlePageChange(pageNum)}
                                isActive={currentPage === pageNum}
                                className="cursor-pointer"
                              >
                                {pageNum}
                              </PaginationLink>
                            </PaginationItem>
                          );
                        })}

                        <PaginationItem>
                          <PaginationNext
                            onClick={() => handlePageChange(currentPage + 1)}
                            className={currentPage >= totalPages ? 'pointer-events-none opacity-50' : 'cursor-pointer'}
                          />
                        </PaginationItem>
                      </PaginationContent>
                    </Pagination>
                  </div>
                )}
              </div>
            )}
          </Card>
        </div>
      </main>
    );
  };

  return (
    <>
      <DashboardHeader
        title="Audit Log"
        description="Monitor audit logs"
      />
      {renderContent()}
    </>
  );
}
