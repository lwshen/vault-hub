import { PATH } from './path';
import { Activity, Vault, Key } from 'lucide-react';

export const DASHBOARD_NAVIGATION = [
  {
    name: 'Dashboard',
    href: PATH.DASHBOARD,
    icon: Activity,
    label: 'Dashboard',
  },
  {
    name: 'Vaults',
    href: PATH.VAULTS,
    icon: Vault,
    label: 'Vaults',
  },
  {
    name: 'API Keys',
    href: PATH.API_KEYS,
    icon: Key,
    label: 'API Keys',
  },
  {
    name: 'Audit Log',
    href: PATH.AUDIT_LOG,
    icon: Activity,
    label: 'Audit Log',
  },
] as const;
