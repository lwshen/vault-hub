import { type ReactNode } from 'react';

interface DashboardHeaderProps {
  title: string;
  description: string;
  actions?: ReactNode;
}

export default function DashboardHeader({ title, description, actions }: DashboardHeaderProps) {
  return (
    <header className="bg-card border-b border-border p-6 flex-shrink-0">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">{title}</h1>
          <p className="text-muted-foreground">{description}</p>
        </div>
        {actions && (
          <div className="flex items-center gap-3">
            {actions}
          </div>
        )}
      </div>
    </header>
  );
}
