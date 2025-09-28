import DashboardLayout from '@/components/layout/dashboard-layout';
import VaultDetailContent from '@/components/dashboard/vault-detail-content';

interface VaultDetailProps {
  vaultId: string;
}

export default function VaultDetail({ vaultId }: VaultDetailProps) {
  return (
    <DashboardLayout>
      <VaultDetailContent vaultId={vaultId} />
    </DashboardLayout>
  );
}
