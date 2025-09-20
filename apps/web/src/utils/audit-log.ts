import { AuditLogActionEnum } from '@lwshen/vault-hub-ts-fetch-client';
import {
  Activity,
  Edit,
  Eye,
  Key,
  Lock,
  LogIn,
  LogOut,
  Plus,
  Trash2,
  UserPlus,
} from 'lucide-react';

// Icon mapping for different audit actions - using correct enum values
export const getIconForAction = (action: AuditLogActionEnum) => {
  const iconMap: { [key in AuditLogActionEnum]: { icon: typeof Lock; color: string; }; } = {
    [AuditLogActionEnum.ReadVault]: { icon: Eye, color: 'text-blue-500' },
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
export const getActionTitle = (action: AuditLogActionEnum) => {
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
export const formatTimestamp = (timestamp: string | Date) => {
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
