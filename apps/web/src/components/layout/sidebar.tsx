import { Button } from '@/components/ui/button';
import { Link, useLocation } from 'wouter';
import { DASHBOARD_NAVIGATION } from '@/const/navigation';

export default function Sidebar() {
  const [pathname] = useLocation();

  return (
    <div className="hidden md:flex w-64 bg-card border-r border-border flex-col h-full">
      {/* Navigation */}
      <nav className="flex-1 p-4 space-y-2">
        {DASHBOARD_NAVIGATION.map(({ href, icon: Icon, label }) => {
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
