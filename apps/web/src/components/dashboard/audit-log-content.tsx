import DashboardHeader from '@/components/layout/dashboard-header';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
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
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip';
import { useAuditLogStore } from '@/stores/audit-store';
import { formatTimestamp, getActionTitle, getIconForAction } from '@/utils/audit-log';
import {
  Activity,
  AlertCircle,
  Key,
  Loader2,
  Lock,
} from 'lucide-react';
import { useEffect } from 'react';

export default function AuditLogContent() {
  const {
    auditLogs,
    metrics,
    isLoading,
    metricsLoading,
    error,
    currentPage,
    totalCount,
    totalPages,
    pageSize,
    fetchMetrics,
    fetchAuditLogs,
    setPageSize,
    setCurrentPage,
  } = useAuditLogStore();

  useEffect(() => {
    fetchMetrics();
    fetchAuditLogs(1);
  }, [fetchMetrics, fetchAuditLogs]);

  const handlePageChange = (page: number) => {
    setCurrentPage(page);
  };

  const handlePageSizeChange = (newPageSize: string) => {
    const size = parseInt(newPageSize);
    setPageSize(size);
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
                  <p className="text-2xl font-bold">
                    {metricsLoading ? (
                      <Loader2 className="h-6 w-6 animate-spin inline" />
                    ) : (
                      metrics?.totalEventsLast30Days?.toLocaleString() || '-'
                    )}
                  </p>
                  <p className="text-sm text-muted-foreground">Total Events</p>
                  <p className="text-xs text-muted-foreground">Last 30 days</p>
                </div>
              </div>
            </Card>
            <Card className="p-4">
              <div className="flex items-center gap-3">
                <Activity className="h-8 w-8 text-orange-500" />
                <div>
                  <p className="text-2xl font-bold">
                    {metricsLoading ? (
                      <Loader2 className="h-6 w-6 animate-spin inline" />
                    ) : (
                      metrics?.eventsCountLast24Hours?.toLocaleString() || '-'
                    )}
                  </p>
                  <p className="text-sm text-muted-foreground">Last 24 Hours</p>
                  <p className="text-xs text-muted-foreground">Recent audit logs</p>
                </div>
              </div>
            </Card>
            <Card className="p-4">
              <div className="flex items-center gap-3">
                <Lock className="h-8 w-8 text-green-500" />
                <div>
                  <p className="text-2xl font-bold">
                    {metricsLoading ? (
                      <Loader2 className="h-6 w-6 animate-spin inline" />
                    ) : (
                      metrics?.vaultEventsLast30Days?.toLocaleString() || '-'
                    )}
                  </p>
                  <p className="text-sm text-muted-foreground">Vault Events</p>
                  <p className="text-xs text-muted-foreground">Last 30 days</p>
                </div>
              </div>
            </Card>
            <Card className="p-4">
              <div className="flex items-center gap-3">
                <Key className="h-8 w-8 text-cyan-500" />
                <div>
                  <p className="text-2xl font-bold">
                    {metricsLoading ? (
                      <Loader2 className="h-6 w-6 animate-spin inline" />
                    ) : (
                      metrics?.apiKeyEventsLast30Days?.toLocaleString() || '-'
                    )}
                  </p>
                  <p className="text-sm text-muted-foreground">API Key Events</p>
                  <p className="text-xs text-muted-foreground">Last 30 days</p>
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
              <>
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead className="w-12"></TableHead>
                      <TableHead>Action</TableHead>
                      <TableHead className="hidden sm:table-cell">Resource</TableHead>
                      <TableHead className="hidden md:table-cell">Method</TableHead>
                      <TableHead className="hidden lg:table-cell">IP Address</TableHead>
                      <TableHead>Time</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {auditLogs.map((audit) => {
                      const { icon: Icon, color } = getIconForAction(audit.action);
                      return (
                        <TableRow key={`${audit.action}-${audit.createdAt}`}>
                          <TableCell>
                            <div className="w-8 h-8 rounded-full bg-muted flex items-center justify-center">
                              <Icon className={`h-4 w-4 ${color}`} />
                            </div>
                          </TableCell>
                          <TableCell>
                            <div className="flex flex-col">
                              <span className="font-medium">{getActionTitle(audit.action)}</span>
                              <span className="text-xs text-muted-foreground sm:hidden">
                                {audit.vault && `Vault: ${audit.vault.name}`}
                                {audit.apiKey && `API Key: ${audit.apiKey.name}`}
                                {!audit.vault && !audit.apiKey && 'User Account'}
                              </span>
                            </div>
                          </TableCell>
                          <TableCell className="hidden sm:table-cell">
                            {audit.vault && (
                              <div className="flex flex-col">
                                <span className="font-medium">{audit.vault.name}</span>
                                <span className="text-xs text-muted-foreground">{audit.vault.uniqueId}</span>
                              </div>
                            )}
                            {audit.apiKey && (
                              <div className="flex flex-col">
                                <span className="font-medium">{audit.apiKey.name}</span>
                                <span className="text-xs text-muted-foreground">ID: {audit.apiKey.id}</span>
                              </div>
                            )}
                            {!audit.vault && !audit.apiKey && (
                              <span className="text-muted-foreground">User Account</span>
                            )}
                          </TableCell>
                          <TableCell className="hidden md:table-cell">
                            {audit.apiKey ? (
                              <Tooltip>
                                <TooltipTrigger asChild>
                                  <span className={'inline-flex items-center px-2 py-1 rounded-full text-xs font-medium cursor-help bg-muted text-muted-foreground'}>
                                    CLI/API
                                  </span>
                                </TooltipTrigger>
                                <TooltipContent>
                                  <div className="space-y-1">
                                    <p>API Key: {audit.apiKey.name}</p>
                                    <p className={`text-sm ${
                                      audit.apiKey.isActive
                                        ? 'text-green-600 dark:text-green-400'
                                        : 'text-red-600 dark:text-red-400'
                                    }`}>
                                      Status: {audit.apiKey.isActive ? 'Active' : 'Inactive'}
                                    </p>
                                    {!audit.apiKey.isActive && (
                                      <p className="text-xs text-orange-600 dark:text-orange-400">
                                        ⚠ This API key was deleted
                                      </p>
                                    )}
                                  </div>
                                </TooltipContent>
                              </Tooltip>
                            ) : (
                              <span className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-muted text-muted-foreground">
                                Web UI
                              </span>
                            )}
                          </TableCell>
                          <TableCell className="hidden lg:table-cell">
                            <span className="text-sm text-muted-foreground">
                              {audit.ipAddress || '-'}
                            </span>
                          </TableCell>
                          <TableCell>
                            <div className="flex flex-col">
                              <span className="text-sm">{formatTimestamp(audit.createdAt)}</span>
                              <span className="text-xs text-muted-foreground lg:hidden">
                                {audit.ipAddress && `IP: ${audit.ipAddress}`}
                              </span>
                            </div>
                          </TableCell>
                        </TableRow>
                      );
                    })}
                  </TableBody>
                </Table>

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
              </>
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
