import { Button } from '@/components/ui/button';
import { Link, useLocation } from 'wouter';
import { PATH } from '@/const/path';
import { 
  Vault, 
  Users, 
  Key, 
  Activity, 
  Plus, 
  Search 
} from 'lucide-react';

export default function Sidebar() {
  const [pathname] = useLocation();

  return (
    <div className="w-64 bg-card border-r border-border flex flex-col h-full">
      {/* Navigation */}
      <nav className="flex-1 p-4 space-y-2">
        <div className="space-y-1">
          <div className="px-3 py-2 text-sm font-medium text-muted-foreground">
            MAIN
          </div>
          <Link href={PATH.DASHBOARD}>
            <Button 
              variant={pathname === PATH.DASHBOARD ? 'default' : 'ghost'} 
              className={`w-full justify-start ${
                pathname === PATH.DASHBOARD 
                  ? 'bg-primary/10 text-primary hover:bg-primary/20' 
                  : ''
              }`}
            >
              <Activity className="h-4 w-4 mr-3" />
              Dashboard
            </Button>
          </Link>
          <Link href={PATH.VAULTS}>
            <Button 
              variant={pathname === PATH.VAULTS ? 'default' : 'ghost'} 
              className={`w-full justify-start ${
                pathname === PATH.VAULTS 
                  ? 'bg-primary/10 text-primary hover:bg-primary/20' 
                  : ''
              }`}
            >
              <Vault className="h-4 w-4 mr-3" />
              Vaults
            </Button>
          </Link>
          <Button variant="ghost" className="w-full justify-start">
            <Key className="h-4 w-4 mr-3" />
            Secrets
          </Button>
          <Button variant="ghost" className="w-full justify-start">
            <Users className="h-4 w-4 mr-3" />
            Team
          </Button>
        </div>

        <div className="space-y-1 pt-4">
          <div className="px-3 py-2 text-sm font-medium text-muted-foreground">
            TOOLS
          </div>
          <Button variant="ghost" className="w-full justify-start">
            <Search className="h-4 w-4 mr-3" />
            Search
          </Button>
          <Button variant="ghost" className="w-full justify-start">
            <Activity className="h-4 w-4 mr-3" />
            Activity Log
          </Button>
        </div>

        <div className="space-y-1 pt-4">
          <div className="px-3 py-2 text-sm font-medium text-muted-foreground">
            SETTINGS
          </div>
          <Button variant="ghost" className="w-full justify-start">
            <Key className="h-4 w-4 mr-3" />
            API Keys
          </Button>
        </div>
      </nav>

      {/* Quick Actions in Sidebar */}
      <div className="p-4 border-t border-border flex-shrink-0">
        <Button className="w-full">
          <Plus className="h-4 w-4 mr-2" />
          New Vault
        </Button>
      </div>
    </div>
  );
} 
