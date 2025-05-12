"use client"

import { useState } from "react"
import { Link, useLocation } from "wouter"
import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { Menu, X, ChevronDown, User, LogOut, Settings } from "lucide-react"

// This would be replaced with actual auth logic in a real app
const useAuth = () => {
  // Mock authentication state - in a real app, this would come from your auth provider
  const [isAuthenticated, setIsAuthenticated] = useState(false)
  const user = isAuthenticated ? { name: "Demo User", email: "user@example.com" } : null

  const login = () => setIsAuthenticated(true)
  const logout = () => setIsAuthenticated(false)

  return { isAuthenticated, user, login, logout }
}

export default function Header() {
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false)
  const [pathname, _] = useLocation()
  const { isAuthenticated, user, login, logout } = useAuth()

  const navigation = [
    { name: "Home", href: "/" },
    { name: "Features", href: "/features" },
    { name: "Pricing", href: "/pricing" },
    { name: "Documentation", href: "/docs" },
  ]

  return (
    <header className="fixed top-0 left-0 right-0 z-50 bg-black/30 backdrop-blur-lg border-b border-white/[0.08]">
      <div className="container mx-auto px-4 sm:px-6">
        <div className="flex items-center justify-between h-16">
          {/* Logo */}
          <div className="flex-shrink-0 flex items-center">
            <Link href="/" className="flex items-center gap-2">
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
              <span className="text-lg font-semibold text-white">VaultHub</span>
            </Link>
          </div>

          {/* Desktop Navigation */}
          <nav className="hidden md:flex space-x-8">
            {navigation.map((item) => (
              <Link
                key={item.name}
                href={item.href}
                className={cn(
                  "text-sm font-medium transition-colors",
                  pathname === item.href ? "text-white" : "text-white/60 hover:text-white",
                )}
              >
                {item.name}
              </Link>
            ))}
          </nav>

          {/* Auth Buttons or User Menu */}
          <div className="hidden md:flex items-center gap-4">
            {isAuthenticated ? (
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button variant="ghost" className="flex items-center gap-2 text-white/80 hover:text-white">
                    <User size={16} />
                    <span>{user?.name}</span>
                    <ChevronDown size={14} />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end" className="w-56 bg-zinc-900 border-white/10">
                  <div className="px-2 py-1.5 text-sm text-white/60">{user?.email}</div>
                  <DropdownMenuSeparator className="bg-white/10" />
                  <DropdownMenuItem className="text-white/80 hover:text-white focus:text-white cursor-pointer">
                    <Settings size={16} className="mr-2" />
                    <span>Settings</span>
                  </DropdownMenuItem>
                  <DropdownMenuItem
                    className="text-white/80 hover:text-white focus:text-white cursor-pointer"
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
                  className="text-white/80 hover:text-white"
                  onClick={login} // For demo purposes
                >
                  Log in
                </Button>
                <Button className="bg-emerald-600 hover:bg-emerald-500 text-white">Register</Button>
              </>
            )}
          </div>

          {/* Mobile menu button */}
          <div className="md:hidden flex items-center">
            <button
              type="button"
              className="text-white/80 hover:text-white"
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
        <div className="md:hidden bg-black/95 border-t border-white/[0.08]">
          <div className="px-2 pt-2 pb-3 space-y-1 sm:px-3">
            {navigation.map((item) => (
              <Link
                key={item.name}
                href={item.href}
                className={cn(
                  "block px-3 py-2 rounded-md text-base font-medium",
                  pathname === item.href
                    ? "text-white bg-white/[0.08]"
                    : "text-white/60 hover:text-white hover:bg-white/[0.04]",
                )}
                onClick={() => setMobileMenuOpen(false)}
              >
                {item.name}
              </Link>
            ))}

            {/* Mobile auth buttons */}
            <div className="pt-4 pb-3 border-t border-white/[0.08]">
              {isAuthenticated ? (
                <>
                  <div className="px-3 py-2 text-white">
                    <div className="text-base font-medium">{user?.name}</div>
                    <div className="text-sm text-white/60">{user?.email}</div>
                  </div>
                  <div className="mt-3 space-y-1">
                    <button className="block w-full text-left px-3 py-2 text-base font-medium text-white/60 hover:text-white hover:bg-white/[0.04] rounded-md">
                      Settings
                    </button>
                    <button
                      onClick={logout}
                      className="block w-full text-left px-3 py-2 text-base font-medium text-white/60 hover:text-white hover:bg-white/[0.04] rounded-md"
                    >
                      Log out
                    </button>
                  </div>
                </>
              ) : (
                <div className="px-3 space-y-2">
                  <Button
                    variant="ghost"
                    className="w-full justify-center text-white/80 hover:text-white hover:bg-white/10"
                    onClick={login} // For demo purposes
                  >
                    Log in
                  </Button>
                  <Button className="w-full justify-center bg-emerald-600 hover:bg-emerald-500 text-white">
                    Register
                  </Button>
                </div>
              )}
            </div>
          </div>
        </div>
      )}
    </header>
  )
}
