import { Card } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { 
  Search,
  Filter,
  Download,
  Activity,
  Plus,
  Users,
  Lock,
  Unlock,
  Key,
  Trash2,
  Edit,
  UserPlus,
  LogIn,
  LogOut
} from 'lucide-react';

export default function ActivityLogContent() {
  const activities = [
    {
      id: 1,
      type: 'vault_access',
      icon: Lock,
      iconColor: 'text-blue-500',
      title: 'Vault accessed',
      description: 'Production API Keys vault was accessed',
      user: 'john@example.com',
      timestamp: '2 minutes ago',
      details: 'Vault: Production API Keys'
    },
    {
      id: 2,
      type: 'secret_created',
      icon: Plus,
      iconColor: 'text-green-500',
      title: 'Secret created',
      description: 'New secret "AWS_ACCESS_KEY" was added',
      user: 'admin@example.com',
      timestamp: '15 minutes ago',
      details: 'Vault: AWS Secrets'
    },
    {
      id: 3,
      type: 'user_login',
      icon: LogIn,
      iconColor: 'text-purple-500',
      title: 'User login',
      description: 'User logged into the system',
      user: 'sarah@example.com',
      timestamp: '1 hour ago',
      details: 'IP: 192.168.1.100'
    },
    {
      id: 4,
      type: 'vault_created',
      icon: Plus,
      iconColor: 'text-green-500',
      title: 'Vault created',
      description: 'New vault "SSL Certificates" was created',
      user: 'admin@example.com',
      timestamp: '2 hours ago',
      details: 'Members: 3 users added'
    },
    {
      id: 5,
      type: 'secret_updated',
      icon: Edit,
      iconColor: 'text-yellow-500',
      title: 'Secret updated',
      description: 'Secret "DATABASE_PASSWORD" was modified',
      user: 'john@example.com',
      timestamp: '3 hours ago',
      details: 'Vault: Database Credentials'
    },
    {
      id: 6,
      type: 'user_invited',
      icon: UserPlus,
      iconColor: 'text-purple-500',
      title: 'User invited',
      description: 'New team member was invited',
      user: 'admin@example.com',
      timestamp: '4 hours ago',
      details: 'Invited: mike@example.com'
    },
    {
      id: 7,
      type: 'vault_unlocked',
      icon: Unlock,
      iconColor: 'text-orange-500',
      title: 'Vault unlocked',
      description: 'OAuth Tokens vault was unlocked',
      user: 'sarah@example.com',
      timestamp: '5 hours ago',
      details: 'Duration: 30 minutes'
    },
    {
      id: 8,
      type: 'secret_deleted',
      icon: Trash2,
      iconColor: 'text-red-500',
      title: 'Secret deleted',
      description: 'Secret "OLD_API_TOKEN" was removed',
      user: 'admin@example.com',
      timestamp: '6 hours ago',
      details: 'Vault: Third-party APIs'
    },
    {
      id: 9,
      type: 'user_logout',
      icon: LogOut,
      iconColor: 'text-gray-500',
      title: 'User logout',
      description: 'User logged out of the system',
      user: 'mike@example.com',
      timestamp: '1 day ago',
      details: 'Session duration: 4 hours'
    },
    {
      id: 10,
      type: 'vault_access',
      icon: Lock,
      iconColor: 'text-blue-500',
      title: 'Vault accessed',
      description: 'Database Credentials vault was accessed',
      user: 'sarah@example.com',
      timestamp: '1 day ago',
      details: 'Vault: Database Credentials'
    }
  ];

  const getActivityTypeLabel = (type: string) => {
    const labels: { [key: string]: string; } = {
      vault_access: 'Vault Access',
      vault_created: 'Vault Created',
      vault_unlocked: 'Vault Unlocked',
      secret_created: 'Secret Created',
      secret_updated: 'Secret Updated',
      secret_deleted: 'Secret Deleted',
      user_login: 'User Login',
      user_logout: 'User Logout',
      user_invited: 'User Invited'
    };
    return labels[type] || 'Unknown';
  };

  return (
    <>
      {/* Top Header */}
      <header className="bg-card border-b border-border p-6 flex-shrink-0">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold tracking-tight">Activity Log</h1>
            <p className="text-muted-foreground">
              Monitor all system activities and user actions
            </p>
          </div>
          <div className="flex items-center gap-3">
            <Button variant="outline" size="sm">
              <Search className="h-4 w-4 mr-2" />
              Search
            </Button>
            <Button variant="outline" size="sm">
              <Filter className="h-4 w-4 mr-2" />
              Filter
            </Button>
            <Button variant="outline" size="sm">
              <Download className="h-4 w-4 mr-2" />
              Export
            </Button>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="flex-1 overflow-y-auto p-6">
        <div className="space-y-4">
          {/* Activity Stats */}
          <div className="grid gap-4 md:grid-cols-4 mb-6">
            <Card className="p-4">
              <div className="flex items-center gap-3">
                <Activity className="h-8 w-8 text-blue-500" />
                <div>
                  <p className="text-2xl font-bold">24</p>
                  <p className="text-sm text-muted-foreground">Today</p>
                </div>
              </div>
            </Card>
            <Card className="p-4">
              <div className="flex items-center gap-3">
                <Lock className="h-8 w-8 text-green-500" />
                <div>
                  <p className="text-2xl font-bold">156</p>
                  <p className="text-sm text-muted-foreground">This Week</p>
                </div>
              </div>
            </Card>
            <Card className="p-4">
              <div className="flex items-center gap-3">
                <Users className="h-8 w-8 text-purple-500" />
                <div>
                  <p className="text-2xl font-bold">12</p>
                  <p className="text-sm text-muted-foreground">Active Users</p>
                </div>
              </div>
            </Card>
            <Card className="p-4">
              <div className="flex items-center gap-3">
                <Key className="h-8 w-8 text-orange-500" />
                <div>
                  <p className="text-2xl font-bold">89</p>
                  <p className="text-sm text-muted-foreground">This Month</p>
                </div>
              </div>
            </Card>
          </div>

          {/* Activity List */}
          <Card className="p-6">
            <div className="space-y-4">
              {activities.map((activity) => {
                const Icon = activity.icon;
                return (
                  <div key={activity.id} className="flex items-start gap-4 p-4 rounded-lg border bg-card hover:bg-muted/50 transition-colors">
                    <div className="flex-shrink-0">
                      <div className="w-10 h-10 rounded-full bg-muted flex items-center justify-center">
                        <Icon className={`h-5 w-5 ${activity.iconColor}`} />
                      </div>
                    </div>
                    <div className="flex-1 min-w-0">
                      <div className="flex items-start justify-between">
                        <div className="flex-1">
                          <h3 className="font-medium text-foreground">{activity.title}</h3>
                          <p className="text-sm text-muted-foreground mt-1">{activity.description}</p>
                          <div className="flex items-center gap-4 mt-2 text-xs text-muted-foreground">
                            <span>By {activity.user}</span>
                            <span>•</span>
                            <span>{activity.timestamp}</span>
                            <span>•</span>
                            <span>{activity.details}</span>
                          </div>
                        </div>
                        <div className="flex-shrink-0">
                          <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-muted text-muted-foreground">
                            {getActivityTypeLabel(activity.type)}
                          </span>
                        </div>
                      </div>
                    </div>
                  </div>
                );
              })}
            </div>

            {/* Load More */}
            <div className="mt-6 text-center">
              <Button variant="outline">
                Load More Activities
              </Button>
            </div>
          </Card>
        </div>
      </main>
    </>
  );
} 
