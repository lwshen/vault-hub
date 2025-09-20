'use client';

import { motion } from 'framer-motion';
import type { Variants } from 'framer-motion';
import { useState } from 'react';

const fadeUpVariants: Variants = {
  hidden: { opacity: 0, y: 30 },
  visible: (i: number) => ({
    opacity: 1,
    y: 0,
    transition: {
      duration: 0.8,
      delay: 0.1 + i * 0.1,
      ease: [0.25, 0.4, 0.25, 1],
    },
  }),
};

const fadeInVariants: Variants = {
  hidden: { opacity: 0 },
  visible: {
    opacity: 1,
    transition: {
      duration: 0.8,
      ease: [0.25, 0.4, 0.25, 1],
    },
  },
};

interface CodeBlockProps {
  code: string;
  language?: string;
  title?: string;
}

const CodeBlock = ({ code, language = 'bash', title }: CodeBlockProps) => (
  <div className="bg-muted/50 rounded-lg overflow-hidden">
    {title && (
      <div className="px-4 py-2 border-b border-border bg-muted/30">
        <span className="text-sm font-medium text-muted-foreground">{title}</span>
      </div>
    )}
    <div className="p-4">
      <pre className="text-sm font-mono text-foreground overflow-x-auto">
        <code className={`language-${language}`}>{code}</code>
      </pre>
    </div>
  </div>
);

interface SectionProps {
  id: string;
  title: string;
  children: React.ReactNode;
  delay?: number;
}

const Section = ({ id, title, children, delay = 0 }: SectionProps) => (
  <motion.section
    id={id}
    custom={delay}
    variants={fadeUpVariants}
    initial="hidden"
    whileInView="visible"
    viewport={{ once: true, margin: '-100px' }}
    className="mb-16"
  >
    <h2 className="text-3xl font-bold mb-6 text-foreground">{title}</h2>
    {children}
  </motion.section>
);

interface NavItem {
  id: string;
  title: string;
  children?: { id: string; title: string; }[];
}

const navigationItems: NavItem[] = [
  {
    id: 'cli-guide',
    title: 'CLI Guide',
    children: [
      { id: 'cli-installation', title: 'CLI Installation' },
      { id: 'cli-authentication', title: 'Authentication' },
      { id: 'cli-commands', title: 'Commands' },
    ],
  },
  {
    id: 'server-setup',
    title: 'Server Setup',
    children: [
      { id: 'server-installation', title: 'Installation' },
      { id: 'configuration', title: 'Configuration' },
      { id: 'first-vault', title: 'Creating Your First Vault' },
    ],
  },
  {
    id: 'api-reference',
    title: 'API Reference',
    children: [
      { id: 'authentication', title: 'Authentication' },
      { id: 'vault-operations', title: 'Vault Operations' },
      { id: 'api-keys', title: 'API Keys' },
    ],
  },
  {
    id: 'security',
    title: 'Security',
    children: [
      { id: 'encryption', title: 'Encryption' },
      { id: 'access-control', title: 'Access Control' },
      { id: 'audit-logs', title: 'Audit Logs' },
    ],
  },
];

export default function Documentation() {
  const [activeSection, setActiveSection] = useState('cli-guide');

  const scrollToSection = (sectionId: string) => {
    const element = document.getElementById(sectionId);
    if (element) {
      element.scrollIntoView({ behavior: 'smooth' });
      setActiveSection(sectionId);
    }
  };

  return (
    <div className="min-h-screen bg-background">
      {/* Hero Section */}
      <section className="relative py-20 md:py-32 overflow-hidden">
        <div className="absolute inset-0 bg-gradient-to-br from-emerald-500/[0.03] via-transparent to-blue-500/[0.03]" />

        <div className="container mx-auto px-4 md:px-6 relative z-10">
          <div className="max-w-4xl mx-auto text-center">
            <motion.div
              custom={0}
              variants={fadeUpVariants}
              initial="hidden"
              animate="visible"
              className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-muted border border-border mb-8"
            >
              <div className="flex items-center justify-center w-5 h-5 bg-blue-500/20 rounded-full">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="12"
                  height="12"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  className="text-blue-500"
                >
                  <path d="M4 19.5v-15A2.5 2.5 0 0 1 6.5 2H19a1 1 0 0 1 1 1v18a1 1 0 0 1-1 1H6.5a2.5 2.5 0 0 1 0-5H20" />
                </svg>
              </div>
              <span className="text-sm text-muted-foreground font-medium">Documentation</span>
            </motion.div>

            <motion.h1
              custom={1}
              variants={fadeUpVariants}
              initial="hidden"
              animate="visible"
              className="text-4xl sm:text-5xl md:text-6xl font-bold mb-6 tracking-tight"
            >
              <span className="bg-clip-text text-transparent bg-gradient-to-b from-foreground to-foreground/80">
                Get Started with
              </span>
              <br />
              <span className="bg-clip-text text-transparent bg-gradient-to-r from-emerald-500 via-primary to-blue-500">
                VaultHub
              </span>
            </motion.h1>

            <motion.p
              custom={2}
              variants={fadeUpVariants}
              initial="hidden"
              animate="visible"
              className="text-lg md:text-xl text-muted-foreground leading-relaxed max-w-3xl mx-auto"
            >
              Start with the CLI for quick access to your secrets, then explore server setup and API integration for team workflows.
            </motion.p>
          </div>
        </div>
      </section>

      {/* Documentation Content */}
      <div className="container mx-auto px-4 md:px-6 py-20">
        <div className="flex flex-col lg:flex-row gap-12">
          {/* Sidebar Navigation */}
          <motion.aside
            variants={fadeInVariants}
            initial="hidden"
            animate="visible"
            className="lg:w-64 lg:flex-shrink-0"
          >
            <div className="sticky top-24">
              <nav className="space-y-6">
                {navigationItems.map((section) => (
                  <div key={section.id}>
                    <button
                      type="button"
                      onClick={() => scrollToSection(section.id)}
                      className={`block w-full text-left font-medium text-sm transition-colors ${
                        activeSection === section.id
                          ? 'text-primary'
                          : 'text-foreground hover:text-primary'
                      }`}
                    >
                      {section.title}
                    </button>
                    {section.children && (
                      <ul className="mt-2 space-y-1 ml-4">
                        {section.children.map((child) => (
                          <li key={child.id}>
                            <button
                              type="button"
                              onClick={() => scrollToSection(child.id)}
                              className={`block w-full text-left text-sm transition-colors ${
                                activeSection === child.id
                                  ? 'text-primary'
                                  : 'text-muted-foreground hover:text-foreground'
                              }`}
                            >
                              {child.title}
                            </button>
                          </li>
                        ))}
                      </ul>
                    )}
                  </div>
                ))}
              </nav>
            </div>
          </motion.aside>

          {/* Main Content */}
          <main className="flex-1 max-w-4xl">
            {/* CLI Guide */}
            <Section id="cli-guide" title="CLI Guide" delay={0}>
              <div className="prose prose-lg max-w-none">
                <p className="text-muted-foreground leading-relaxed mb-8">
                  The VaultHub CLI is the primary way to interact with VaultHub. It provides secure, programmatic access
                  to your vaults and is perfect for development workflows, CI/CD pipelines, and automation.
                  Start here to quickly access your secrets without needing to set up a server.
                </p>

                <div id="cli-installation" className="mb-12">
                  <h3 className="text-2xl font-semibold mb-4 text-foreground">CLI Installation</h3>

                  <h4 className="text-lg font-medium mb-3 text-foreground">Download Pre-built Binaries</h4>
                  <p className="text-muted-foreground mb-4">
                    Download the latest CLI binary for your platform from the{' '}
                    <a
                      href="https://github.com/lwshen/vault-hub/releases/latest"
                      target="_blank"
                      rel="noopener noreferrer"
                      className="text-primary hover:text-primary/80 underline underline-offset-4"
                    >
                      GitHub releases page
                    </a>.
                  </p>

                  <div className="bg-primary/5 border border-primary/20 rounded-lg p-6 mb-6">
                    <div className="flex items-center justify-between">
                      <div>
                        <h5 className="font-semibold text-foreground mb-1">Latest Release</h5>
                        <p className="text-sm text-muted-foreground">Get the most recent version of VaultHub CLI</p>
                      </div>
                      <a
                        href="https://github.com/lwshen/vault-hub/releases/latest"
                        target="_blank"
                        rel="noopener noreferrer"
                        className="inline-flex items-center gap-2 px-4 py-2 bg-primary text-primary-foreground rounded-lg hover:bg-primary/90 transition-colors"
                      >
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                          <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
                          <polyline points="7,10 12,15 17,10" />
                          <line x1="12" x2="12" y1="15" y2="3" />
                        </svg>
                        Download Latest
                      </a>
                    </div>
                  </div>

                  <div className="grid md:grid-cols-3 gap-4 mb-6">
                    <div className="bg-card border border-border rounded-lg p-4 text-center">
                      <div className="w-8 h-8 mx-auto mb-2 bg-blue-500/10 rounded-lg flex items-center justify-center">
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" className="text-blue-500">
                          <rect width="20" height="14" x="2" y="3" rx="2" ry="2" />
                          <line x1="8" x2="16" y1="21" y2="21" />
                          <line x1="12" x2="12" y1="17" y2="21" />
                        </svg>
                      </div>
                      <div className="text-sm font-medium">Linux</div>
                      <div className="text-xs text-muted-foreground">amd64, arm64</div>
                    </div>
                    <div className="bg-card border border-border rounded-lg p-4 text-center">
                      <div className="w-8 h-8 mx-auto mb-2 bg-emerald-500/10 rounded-lg flex items-center justify-center">
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" className="text-emerald-500">
                          <rect width="20" height="14" x="2" y="3" rx="2" ry="2" />
                          <line x1="8" x2="16" y1="21" y2="21" />
                          <line x1="12" x2="12" y1="17" y2="21" />
                        </svg>
                      </div>
                      <div className="text-sm font-medium">Windows</div>
                      <div className="text-xs text-muted-foreground">amd64</div>
                    </div>
                    <div className="bg-card border border-border rounded-lg p-4 text-center">
                      <div className="w-8 h-8 mx-auto mb-2 bg-amber-500/10 rounded-lg flex items-center justify-center">
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" className="text-amber-500">
                          <rect width="20" height="14" x="2" y="3" rx="2" ry="2" />
                          <line x1="8" x2="16" y1="21" y2="21" />
                          <line x1="12" x2="12" y1="17" y2="21" />
                        </svg>
                      </div>
                      <div className="text-sm font-medium">macOS</div>
                      <div className="text-xs text-muted-foreground">amd64, arm64</div>
                    </div>
                  </div>

                  <h4 className="text-lg font-medium mb-3 text-foreground">Build from Source</h4>
                  <CodeBlock
                    title="Build CLI"
                    code={`# Clone the repository
git clone https://github.com/lwshen/vault-hub.git
cd vault-hub

# Build the CLI
go build -o vault-hub-cli ./apps/cli/main.go

# Make it executable and move to PATH (Linux/macOS)
chmod +x vault-hub-cli
sudo mv vault-hub-cli /usr/local/bin/vault-hub`}
                  />
                </div>

                <div id="cli-authentication" className="mb-12">
                  <h3 className="text-2xl font-semibold mb-4 text-foreground">Authentication</h3>
                  <p className="text-muted-foreground mb-4">
                    The CLI uses API keys for authentication. First, create an API key in the web interface:
                  </p>

                  <ol className="list-decimal list-inside text-muted-foreground space-y-2 mb-6">
                    <li>Log into the VaultHub web interface</li>
                    <li>Navigate to Dashboard → API Keys</li>
                    <li>Click "Create API Key" and give it a name</li>
                    <li>Copy the generated API key (starts with <code className="bg-muted px-1 rounded">vhub_</code>)</li>
                  </ol>

                  <CodeBlock
                    title="Set API Key"
                    code={`# Set the API key as an environment variable
export VAULT_HUB_API_KEY=vhub_your_api_key_here

# Or pass it directly to commands
vault-hub --api-key vhub_your_api_key_here list`}
                  />
                </div>

                <div id="cli-commands" className="mb-12">
                  <h3 className="text-2xl font-semibold mb-4 text-foreground">Commands</h3>

                  <div className="space-y-8">
                    <div>
                      <h4 className="text-lg font-medium mb-3 text-foreground">List Vaults</h4>
                      <CodeBlock
                        code={`# List all accessible vaults
vault-hub list

# Short form
vault-hub ls`}
                      />
                    </div>

                    <div>
                      <h4 className="text-lg font-medium mb-3 text-foreground">Get Vault Contents</h4>
                      <CodeBlock
                        code={`# Get vault by name
vault-hub get --name production-secrets

# Get vault by ID
vault-hub get --id vault-uuid-here

# Export to .env file
vault-hub get --name production-secrets --output .env

# Execute command with environment variables
vault-hub get --name production-secrets --exec "npm start"`}
                      />
                    </div>

                    <div>
                      <h4 className="text-lg font-medium mb-3 text-foreground">Version Information</h4>
                      <CodeBlock
                        code={`# Show version and build information
vault-hub version`}
                      />
                    </div>
                  </div>
                </div>
              </div>
            </Section>

            {/* Server Setup */}
            <Section id="server-setup" title="Server Setup" delay={1}>
              <div className="prose prose-lg max-w-none">
                <p className="text-muted-foreground leading-relaxed mb-8">
                  Set up and configure the VaultHub server for your team or organization.
                  The server provides the web interface and API endpoints for vault management.
                </p>

                <div id="server-installation" className="mb-12">
                  <h3 className="text-2xl font-semibold mb-4 text-foreground">Installation</h3>
                  <p className="text-muted-foreground mb-4">
                    VaultHub consists of a backend server and a web interface. You can run it locally or deploy it to your infrastructure.
                  </p>

                  <h4 className="text-lg font-medium mb-3 text-foreground">Prerequisites</h4>
                  <ul className="list-disc list-inside text-muted-foreground mb-6 space-y-1">
                    <li>Go 1.24+ for the backend server</li>
                    <li>Node.js 22+ and pnpm for the web interface (optional)</li>
                    <li>Database: SQLite (default), MySQL, or PostgreSQL</li>
                  </ul>

                  <h4 className="text-lg font-medium mb-3 text-foreground">Quick Start</h4>
                  <CodeBlock
                    title="Clone and Setup"
                    code={`# Clone the repository
git clone https://github.com/lwshen/vault-hub.git
cd vault-hub

# Set required environment variables
export JWT_SECRET=your-jwt-secret-here
export ENCRYPTION_KEY=$(openssl rand -base64 32)

# Run the server
go run ./apps/server/main.go`}
                  />
                </div>

                <div id="configuration" className="mb-12">
                  <h3 className="text-2xl font-semibold mb-4 text-foreground">Configuration</h3>
                  <p className="text-muted-foreground mb-4">
                    VaultHub can be configured using environment variables. Here are the essential settings:
                  </p>

                  <div className="bg-card border border-border rounded-lg p-6 mb-6">
                    <h4 className="text-lg font-medium mb-3 text-foreground">Required Variables</h4>
                    <div className="space-y-3 text-sm">
                      <div>
                        <code className="bg-muted px-2 py-1 rounded text-primary">JWT_SECRET</code>
                        <span className="text-muted-foreground ml-2">Secret key for JWT token signing</span>
                      </div>
                      <div>
                        <code className="bg-muted px-2 py-1 rounded text-primary">ENCRYPTION_KEY</code>
                        <span className="text-muted-foreground ml-2">AES-256 encryption key for vault data</span>
                      </div>
                    </div>
                  </div>

                  <div className="bg-card border border-border rounded-lg p-6">
                    <h4 className="text-lg font-medium mb-3 text-foreground">Optional Variables</h4>
                    <div className="space-y-3 text-sm">
                      <div>
                        <code className="bg-muted px-2 py-1 rounded text-primary">APP_PORT</code>
                        <span className="text-muted-foreground ml-2">Server port (default: 3000)</span>
                      </div>
                      <div>
                        <code className="bg-muted px-2 py-1 rounded text-primary">DATABASE_TYPE</code>
                        <span className="text-muted-foreground ml-2">sqlite|mysql|postgres (default: sqlite)</span>
                      </div>
                      <div>
                        <code className="bg-muted px-2 py-1 rounded text-primary">DATABASE_URL</code>
                        <span className="text-muted-foreground ml-2">Database connection string</span>
                      </div>
                    </div>
                  </div>
                </div>

                <div id="first-vault" className="mb-12">
                  <h3 className="text-2xl font-semibold mb-4 text-foreground">Creating Your First Vault</h3>
                  <p className="text-muted-foreground mb-4">
                    Once VaultHub is running, you can create your first vault through the web interface:
                  </p>

                  <ol className="list-decimal list-inside text-muted-foreground space-y-2 mb-6">
                    <li>Navigate to <code className="bg-muted px-2 py-1 rounded">http://localhost:3000</code></li>
                    <li>Register a new account or log in</li>
                    <li>Go to the Dashboard and click "Create Vault"</li>
                    <li>Enter a name and key-value pairs for your environment variables</li>
                    <li>Save your vault - all values are automatically encrypted</li>
                  </ol>

                  <div className="bg-emerald-500/10 border border-emerald-500/20 rounded-lg p-4">
                    <div className="flex items-start gap-3">
                      <svg className="w-5 h-5 text-emerald-500 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                      </svg>
                      <div>
                        <p className="text-sm font-medium text-emerald-600 dark:text-emerald-400">Security Note</p>
                        <p className="text-sm text-emerald-700 dark:text-emerald-300 mt-1">
                          All vault values are encrypted with AES-256-GCM before being stored in the database.
                          Your encryption key should be kept secure and backed up safely.
                        </p>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </Section>

            {/* API Reference */}
            <Section id="api-reference" title="API Reference" delay={2}>
              <div className="prose prose-lg max-w-none">
                <p className="text-muted-foreground leading-relaxed mb-8">
                  VaultHub provides a RESTful API with OpenAPI 3.0 specification.
                  All API endpoints use JSON for data exchange and require proper authentication.
                </p>

                <div id="authentication" className="mb-12">
                  <h3 className="text-2xl font-semibold mb-4 text-foreground">Authentication</h3>
                  <p className="text-muted-foreground mb-4">
                    VaultHub supports two authentication methods depending on the endpoint:
                  </p>

                  <div className="grid md:grid-cols-2 gap-6 mb-6">
                    <div className="bg-card border border-border rounded-lg p-6">
                      <h4 className="text-lg font-medium mb-3 text-foreground">JWT Authentication</h4>
                      <p className="text-sm text-muted-foreground mb-3">
                        Used for web interface and user management endpoints.
                      </p>
                      <CodeBlock
                        language="http"
                        code={`POST /api/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password"
}`}
                      />
                    </div>

                    <div className="bg-card border border-border rounded-lg p-6">
                      <h4 className="text-lg font-medium mb-3 text-foreground">API Key Authentication</h4>
                      <p className="text-sm text-muted-foreground mb-3">
                        Used for CLI and programmatic access to vault data.
                      </p>
                      <CodeBlock
                        language="http"
                        code={`GET /api/cli/vaults
Authorization: Bearer vhub_your_api_key_here`}
                      />
                    </div>
                  </div>
                </div>

                <div id="vault-operations" className="mb-12">
                  <h3 className="text-2xl font-semibold mb-4 text-foreground">Vault Operations</h3>

                  <div className="space-y-8">
                    <div>
                      <h4 className="text-lg font-medium mb-3 text-foreground">List Vaults (CLI)</h4>
                      <CodeBlock
                        language="http"
                        code={`GET /api/cli/vaults
Authorization: Bearer vhub_your_api_key_here

Response:
{
  "vaults": [
    {
      "uniqueId": "vault-uuid",
      "name": "production-secrets",
      "description": "Production environment variables",
      "createdAt": "2025-01-01T00:00:00Z"
    }
  ]
}`}
                      />
                    </div>

                    <div>
                      <h4 className="text-lg font-medium mb-3 text-foreground">Get Vault by Name (CLI)</h4>
                      <CodeBlock
                        language="http"
                        code={`GET /api/cli/vault/name/{name}
Authorization: Bearer vhub_your_api_key_here

Response:
{
  "uniqueId": "vault-uuid",
  "name": "production-secrets",
  "description": "Production environment variables",
  "value": {
    "API_KEY": "secret-api-key",
    "DATABASE_URL": "postgresql://...",
    "REDIS_URL": "redis://..."
  },
  "createdAt": "2025-01-01T00:00:00Z"
}`}
                      />
                    </div>

                    <div>
                      <h4 className="text-lg font-medium mb-3 text-foreground">Create Vault (Web)</h4>
                      <CodeBlock
                        language="http"
                        code={`POST /api/vaults
Authorization: Bearer jwt_token_here
Content-Type: application/json

{
  "name": "staging-secrets",
  "description": "Staging environment variables",
  "value": {
    "API_KEY": "staging-api-key",
    "DATABASE_URL": "postgresql://staging..."
  }
}`}
                      />
                    </div>
                  </div>
                </div>

                <div id="api-keys" className="mb-12">
                  <h3 className="text-2xl font-semibold mb-4 text-foreground">API Key Management</h3>
                  <p className="text-muted-foreground mb-4">
                    Manage API keys for programmatic access to vaults.
                  </p>

                  <div className="space-y-6">
                    <div>
                      <h4 className="text-lg font-medium mb-3 text-foreground">Create API Key</h4>
                      <CodeBlock
                        language="http"
                        code={`POST /api/api-keys
Authorization: Bearer jwt_token_here
Content-Type: application/json

{
  "name": "CI/CD Pipeline Key",
  "description": "Key for accessing production secrets in CI"
}

Response:
{
  "id": "key-uuid",
  "name": "CI/CD Pipeline Key",
  "key": "vhub_generated_api_key_here",
  "createdAt": "2025-01-01T00:00:00Z"
}`}
                      />
                    </div>

                    <div>
                      <h4 className="text-lg font-medium mb-3 text-foreground">List API Keys</h4>
                      <CodeBlock
                        language="http"
                        code={`GET /api/api-keys
Authorization: Bearer jwt_token_here

Response:
{
  "apiKeys": [
    {
      "id": "key-uuid",
      "name": "CI/CD Pipeline Key",
      "createdAt": "2025-01-01T00:00:00Z",
      "lastUsed": "2025-01-02T12:00:00Z"
    }
  ]
}`}
                      />
                    </div>
                  </div>
                </div>
              </div>
            </Section>

            {/* Security */}
            <Section id="security" title="Security" delay={3}>
              <div className="prose prose-lg max-w-none">
                <p className="text-muted-foreground leading-relaxed mb-8">
                  VaultHub implements multiple layers of security to protect your sensitive data.
                </p>

                <div id="encryption" className="mb-12">
                  <h3 className="text-2xl font-semibold mb-4 text-foreground">Encryption</h3>

                  <div className="bg-card border border-border rounded-lg p-6 mb-6">
                    <h4 className="text-lg font-medium mb-3 text-foreground">AES-256-GCM Encryption</h4>
                    <p className="text-muted-foreground mb-4">
                      All vault values are encrypted using AES-256-GCM before being stored in the database.
                      This provides both confidentiality and authenticity.
                    </p>

                    <div className="grid md:grid-cols-2 gap-4">
                      <div className="bg-muted/30 rounded-lg p-4">
                        <h5 className="font-medium text-foreground mb-2">Key Features</h5>
                        <ul className="text-sm text-muted-foreground space-y-1">
                          <li>• 256-bit encryption key</li>
                          <li>• Galois/Counter Mode (GCM)</li>
                          <li>• Authenticated encryption</li>
                          <li>• Unique IV per encryption</li>
                        </ul>
                      </div>
                      <div className="bg-muted/30 rounded-lg p-4">
                        <h5 className="font-medium text-foreground mb-2">Security Benefits</h5>
                        <ul className="text-sm text-muted-foreground space-y-1">
                          <li>• Data confidentiality</li>
                          <li>• Data integrity</li>
                          <li>• Authentication</li>
                          <li>• Resistance to tampering</li>
                        </ul>
                      </div>
                    </div>
                  </div>

                  <div className="bg-amber-500/10 border border-amber-500/20 rounded-lg p-4">
                    <div className="flex items-start gap-3">
                      <svg className="w-5 h-5 text-amber-500 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.732-.833-2.5 0L4.314 15.5c-.77.833.192 2.5 1.732 2.5z" />
                      </svg>
                      <div>
                        <p className="text-sm font-medium text-amber-600 dark:text-amber-400">Important</p>
                        <p className="text-sm text-amber-700 dark:text-amber-300 mt-1">
                          Keep your <code className="bg-amber-100 dark:bg-amber-900/30 px-1 rounded">ENCRYPTION_KEY</code> secure and backed up.
                          Without it, encrypted vault data cannot be recovered.
                        </p>
                      </div>
                    </div>
                  </div>
                </div>

                <div id="access-control" className="mb-12">
                  <h3 className="text-2xl font-semibold mb-4 text-foreground">Access Control</h3>

                  <div className="grid md:grid-cols-2 gap-6 mb-6">
                    <div className="bg-card border border-border rounded-lg p-6">
                      <h4 className="text-lg font-medium mb-3 text-foreground">Authentication Methods</h4>
                      <div className="space-y-3 text-sm">
                        <div className="flex items-center gap-3">
                          <div className="w-2 h-2 bg-emerald-500 rounded-full"></div>
                          <span className="text-foreground">JWT tokens for web access</span>
                        </div>
                        <div className="flex items-center gap-3">
                          <div className="w-2 h-2 bg-blue-500 rounded-full"></div>
                          <span className="text-foreground">API keys for programmatic access</span>
                        </div>
                        <div className="flex items-center gap-3">
                          <div className="w-2 h-2 bg-violet-500 rounded-full"></div>
                          <span className="text-foreground">Optional OIDC integration</span>
                        </div>
                      </div>
                    </div>

                    <div className="bg-card border border-border rounded-lg p-6">
                      <h4 className="text-lg font-medium mb-3 text-foreground">Route Protection</h4>
                      <div className="space-y-3 text-sm">
                        <div>
                          <span className="text-foreground font-medium">Public Routes:</span>
                          <div className="text-muted-foreground">Login, registration, static assets</div>
                        </div>
                        <div>
                          <span className="text-foreground font-medium">JWT Protected:</span>
                          <div className="text-muted-foreground">Web interface, user management</div>
                        </div>
                        <div>
                          <span className="text-foreground font-medium">API Key Protected:</span>
                          <div className="text-muted-foreground">CLI endpoints, vault access</div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>

                <div id="audit-logs" className="mb-12">
                  <h3 className="text-2xl font-semibold mb-4 text-foreground">Audit Logs</h3>
                  <p className="text-muted-foreground mb-6">
                    VaultHub maintains comprehensive audit logs of all operations for security monitoring and compliance.
                  </p>

                  <div className="bg-card border border-border rounded-lg p-6">
                    <h4 className="text-lg font-medium mb-4 text-foreground">Logged Events</h4>
                    <div className="grid md:grid-cols-2 gap-6">
                      <div>
                        <h5 className="font-medium text-foreground mb-2">Vault Operations</h5>
                        <ul className="text-sm text-muted-foreground space-y-1">
                          <li>• Vault creation and deletion</li>
                          <li>• Vault value updates</li>
                          <li>• Vault access (read operations)</li>
                          <li>• Permission changes</li>
                        </ul>
                      </div>
                      <div>
                        <h5 className="font-medium text-foreground mb-2">Authentication Events</h5>
                        <ul className="text-sm text-muted-foreground space-y-1">
                          <li>• User login and logout</li>
                          <li>• API key creation and usage</li>
                          <li>• Failed authentication attempts</li>
                          <li>• Session management</li>
                        </ul>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </Section>
          </main>
        </div>
      </div>
    </div>
  );
}
