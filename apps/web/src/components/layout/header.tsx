'use client';

import { useState } from 'react';
import { Link, useLocation } from 'wouter';
import { cn } from '@/lib/utils';
import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Menu, X, ChevronDown, User, LogOut, Settings } from 'lucide-react';
import { ThemeToggle } from '@/components/theme-mode-toggle';
import useAuth from '@/hooks/use-auth';
import { PATH } from '@/const/path';
import { DASHBOARD_NAVIGATION } from '@/const/navigation';

export default function Header() {
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const [pathname, navigate] = useLocation();
  const { isAuthenticated, user, logout, isLoading } = useAuth();

  // Check if user is on dashboard pages
  const isDashboardPage = isAuthenticated && (
    pathname.startsWith('/dashboard') ||
    pathname.startsWith('/vaults') ||
    pathname.startsWith('/api-keys') ||
    pathname.startsWith('/audit-log')
  );

  const navigation = isAuthenticated ? [
    { name: 'Dashboard', href: PATH.DASHBOARD },
    { name: 'Features', href: '/features' },
    { name: 'Documentation', href: '/docs' },
  ] : [
    { name: 'Home', href: '/' },
    { name: 'Features', href: '/features' },
    { name: 'Documentation', href: '/docs' },
  ];


  const handleLogin = () => {
    navigate(PATH.LOGIN);
    setMobileMenuOpen(false);
  };

  const handleSignup = () => {
    navigate(PATH.SIGNUP);
    setMobileMenuOpen(false);
  };

  return (
    <header className="sticky top-0 left-0 right-0 z-50 bg-background/80 backdrop-blur-lg border-b border-border">
      <div className="container mx-auto px-4 sm:px-6">
        <div className="flex items-center justify-between h-16">
          {/* Logo */}
          <div className="flex-shrink-0 flex items-center">
            <Link href={isAuthenticated ? PATH.DASHBOARD : '/'} className="flex items-center gap-2">
              <div className="flex items-center justify-center w-8 h-8 bg-emerald-500/20 rounded-md">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="16"
                  height="16"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  className="text-emerald-400"
                >
                  <rect width="18" height="11" x="3" y="11" rx="2" ry="2" />
                  <path d="M7 11V7a5 5 0 0 1 10 0v4" />
                </svg>
              </div>
              <span className="text-lg font-semibold text-foreground">VaultHub</span>
            </Link>
          </div>

          {/* Desktop Navigation */}
          <nav className="hidden md:flex space-x-8">
            {navigation.map((item) => (
              <Link
                key={item.name}
                href={item.href}
                className={cn(
                  'text-sm font-medium transition-colors',
                  pathname === item.href
                    ? 'text-foreground'
                    : 'text-foreground/60 hover:text-foreground',
                )}
              >
                {item.name}
              </Link>
            ))}
          </nav>

          {/* Auth Buttons or User Menu */}
          <div className="hidden md:flex items-center gap-4">
            {/* Theme Toggle */}
            <ThemeToggle />

            {isLoading ? (
              <div className="h-8 w-24 animate-pulse rounded-md bg-muted" />
            ) : isAuthenticated ? (
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button
                    variant="ghost"
                    className="flex items-center gap-2 text-foreground/80 hover:text-foreground"
                  >
                    <User size={16} />
                    <span>{user?.name}</span>
                    <ChevronDown size={14} />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end" className="w-56 bg-popover border-border">
                  <div className="px-2 py-1.5 text-sm text-popover-foreground/80">
                    {user?.email}
                  </div>
                  <DropdownMenuSeparator className="bg-border" />
                  <DropdownMenuItem className="cursor-pointer focus:bg-accent">
                    <Settings size={16} className="mr-2" />
                    <span>Settings</span>
                  </DropdownMenuItem>
                  <DropdownMenuItem
                    className="cursor-pointer focus:bg-accent"
                    onClick={logout}
                  >
                    <LogOut size={16} className="mr-2" />
                    <span>Log out</span>
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            ) : (
              <>
                <Button
                  variant="ghost"
                  className="text-foreground/80 hover:text-foreground"
                  onClick={handleLogin}
                >
                  Log in
                </Button>
                <Button
                  className="bg-primary hover:bg-primary/90 text-primary-foreground"
                  onClick={handleSignup}
                >
                  Register
                </Button>
              </>
            )}
          </div>

          {/* Mobile menu button and theme toggle */}
          <div className="md:hidden flex items-center gap-2">
            <ThemeToggle />
            <button
              type="button"
              className="text-foreground/80 hover:text-foreground"
              onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
            >
              <span className="sr-only">Open main menu</span>
              {mobileMenuOpen ? <X size={24} /> : <Menu size={24} />}
            </button>
          </div>
        </div>
      </div>

      {/* Mobile menu */}
      {mobileMenuOpen && (
        <div className="md:hidden bg-background/95 border-t border-border">
          <div className="px-2 pt-2 pb-3 space-y-1 sm:px-3">
            {isDashboardPage ? (
              // Show dashboard navigation when on dashboard pages
              <>
                <div className="px-3 py-2 text-sm font-medium text-foreground/60 border-b border-border mb-2">
                  Dashboard Navigation
                </div>
                {DASHBOARD_NAVIGATION.map((item) => {
                  const Icon = item.icon;
                  return (
                    <Link
                      key={item.name}
                      href={item.href}
                      className={cn(
                        'flex items-center gap-3 px-3 py-2 rounded-md text-base font-medium',
                        pathname === item.href
                          ? 'text-foreground bg-accent'
                          : 'text-foreground/60 hover:text-foreground hover:bg-accent/50',
                      )}
                      onClick={() => setMobileMenuOpen(false)}
                    >
                      <Icon className="h-4 w-4" />
                      {item.name}
                    </Link>
                  );
                })}
              </>
            ) : (
              // Show regular navigation when not on dashboard pages
              navigation.map((item) => (
                <Link
                  key={item.name}
                  href={item.href}
                  className={cn(
                    'block px-3 py-2 rounded-md text-base font-medium',
                    pathname === item.href
                      ? 'text-foreground bg-accent'
                      : 'text-foreground/60 hover:text-foreground hover:bg-accent/50',
                  )}
                  onClick={() => setMobileMenuOpen(false)}
                >
                  {item.name}
                </Link>
              ))
            )}

            {/* Mobile auth buttons */}
            <div className="pt-4 pb-3 border-t border-border">
              {isAuthenticated ? (
                <>
                  <div className="px-3 py-2 text-foreground">
                    <div className="text-base font-medium">{user?.name}</div>
                    <div className="text-sm text-foreground/60">{user?.email}</div>
                  </div>
                  <div className="mt-3 space-y-1">
                    <Button
                      variant="ghost"
                      className="w-full justify-start text-foreground/60 hover:text-foreground hover:bg-accent/50"
                    >
                      Settings
                    </Button>
                    <Button
                      variant="ghost"
                      onClick={logout}
                      className="w-full justify-start text-foreground/60 hover:text-foreground hover:bg-accent/50"
                    >
                      Log out
                    </Button>
                  </div>
                </>
              ) : (
                <div className="px-3 space-y-2">
                  <Button
                    variant="ghost"
                    className="w-full justify-center text-foreground/80 hover:text-foreground hover:bg-accent"
                    onClick={handleLogin}
                  >
                    Log in
                  </Button>
                  <Button
                    className="w-full justify-center bg-primary hover:bg-primary/90 text-primary-foreground"
                    onClick={handleSignup}
                  >
                    Register
                  </Button>
                </div>
              )}
            </div>
          </div>
        </div>
      )}
    </header>
  );
}
