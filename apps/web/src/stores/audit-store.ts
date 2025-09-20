import { auditApi } from '@/apis/api';
import type { AuditLog, AuditMetricsResponse } from '@lwshen/vault-hub-ts-fetch-client';
import { create } from 'zustand';

interface AuditLogState {
  auditLogs: AuditLog[];
  metrics: AuditMetricsResponse | null;
  isLoading: boolean;
  metricsLoading: boolean;
  error: string | null;
  currentPage: number;
  totalCount: number;
  totalPages: number;
  pageSize: number;
}

interface AuditLogActions {
  fetchMetrics: () => Promise<void>;
  fetchAuditLogs: (page: number) => Promise<void>;
  setPageSize: (pageSize: number) => void;
  setCurrentPage: (page: number) => void;
  clearError: () => void;
  reset: () => void;
}

type AuditLogStore = AuditLogState & AuditLogActions;

const initialState: AuditLogState = {
  auditLogs: [],
  metrics: null,
  isLoading: true,
  metricsLoading: true,
  error: null,
  currentPage: 1,
  totalCount: 0,
  totalPages: 0,
  pageSize: 20,
};

export const useAuditLogStore = create<AuditLogStore>((set, get) => ({
  ...initialState,

  fetchMetrics: async () => {
    set({ metricsLoading: true });
    try {
      const metricsData = await auditApi.getAuditMetrics();
      set({ metrics: metricsData });
    } catch (err) {
      console.error('Failed to fetch metrics:', err);
    } finally {
      set({ metricsLoading: false });
    }
  },

  fetchAuditLogs: async (page: number) => {
    const { pageSize } = get();
    set({ isLoading: true, error: null });

    try {
      const response = await auditApi.getAuditLogs(pageSize, page);
      const newLogs = response.auditLogs || [];

      set({
        auditLogs: newLogs,
        totalCount: response.totalCount || 0,
        totalPages: Math.ceil((response.totalCount || 0) / pageSize),
        currentPage: page,
      });
    } catch (err) {
      set({
        error: err instanceof Error ? err.message : 'Failed to fetch audit logs',
      });
    } finally {
      set({ isLoading: false });
    }
  },

  setPageSize: (pageSize: number) => {
    set({ pageSize, currentPage: 1 });
    // Fetch audit logs with new page size
    get().fetchAuditLogs(1);
  },

  setCurrentPage: (page: number) => {
    const { totalPages } = get();
    if (page >= 1 && page <= totalPages) {
      get().fetchAuditLogs(page);
    }
  },

  clearError: () => set({ error: null }),

  reset: () => set(initialState),
}));
