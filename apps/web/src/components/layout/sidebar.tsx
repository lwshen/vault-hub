import { Button } from '@/components/ui/button';
import { Link, useLocation } from 'wouter';
import { PATH } from '@/const/path';
import {
  Vault,
  Activity,
  Key,
} from 'lucide-react';

export default function Sidebar() {
  const [pathname] = useLocation();


  const navItems = [
    { href: PATH.DASHBOARD, icon: Activity, label: 'Dashboard' },
    { href: PATH.VAULTS, icon: Vault, label: 'Vaults' },
    { href: PATH.API_KEYS, icon: Key, label: 'API Keys' },
    { href: PATH.AUDIT_LOG, icon: Activity, label: 'Audit Log' },
  ];

  return (
    <div className="hidden md:flex w-64 bg-card border-r border-border flex-col h-full">
      {/* Navigation */}
      <nav className="flex-1 p-4 space-y-2">
        {navItems.map(({ href, icon: Icon, label }) => {
          const isActive = pathname === href;
          return (
            <Link key={href} href={href}>
              <Button
                variant={isActive ? 'default' : 'ghost'}
                className={`w-full justify-start ${
                  isActive
                    ? 'bg-primary/10 text-primary hover:bg-primary/20'
                    : ''
                }`}
              >
                <Icon className="h-4 w-4 mr-3" />
                {label}
              </Button>
            </Link>
          );
        })}
      </nav>

      {/* Quick Actions in Sidebar */}
      {/* <div className="p-4 border-t border-border flex-shrink-0">
        <Button className="w-full">
          <Plus className="h-4 w-4 mr-2" />
          New Vault
        </Button>
      </div> */}
    </div>
  );
}
