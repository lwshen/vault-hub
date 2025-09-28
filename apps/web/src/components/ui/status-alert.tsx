import { cn } from '@/lib/utils';

export type AlertVariant = 'info' | 'warning' | 'error' | 'success';

interface StatusAlertProps {
  variant: AlertVariant;
  title?: string;
  children: React.ReactNode;
  className?: string;
}

const variantStyles: Record<AlertVariant, string> = {
  info: 'text-blue-700 bg-blue-50 border-blue-200 dark:text-blue-400 dark:bg-blue-900/20 dark:border-blue-800',
  warning: 'text-amber-700 bg-amber-50 border-amber-200 dark:text-amber-400 dark:bg-amber-900/20 dark:border-amber-800',
  error: 'text-red-600 bg-red-50 border-red-200 dark:text-red-400 dark:bg-red-900/20 dark:border-red-800',
  success: 'text-green-700 bg-green-50 border-green-200 dark:text-green-400 dark:bg-green-900/20 dark:border-green-800',
};

export function StatusAlert({ variant, title, children, className }: StatusAlertProps) {
  return (
    <div
      className={cn(
        'p-3 text-sm border rounded-md',
        variantStyles[variant],
        className
      )}
    >
      {title && <strong>{title}: </strong>}
      {children}
    </div>
  );
}