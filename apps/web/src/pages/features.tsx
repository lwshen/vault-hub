'use client';

import { motion } from 'framer-motion';
import type { Variants } from 'framer-motion';

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

interface FeatureCardProps {
  icon: React.ReactNode;
  title: string;
  description: string;
  delay: number;
}

const FeatureCard = ({ icon, title, description, delay }: FeatureCardProps) => (
  <motion.div
    custom={delay}
    variants={fadeUpVariants}
    initial="hidden"
    whileInView="visible"
    viewport={{ once: true, margin: '-100px' }}
    className="group relative p-6 md:p-8 rounded-2xl bg-card border border-border hover:border-border/60 transition-all duration-300 hover:shadow-lg hover:shadow-primary/5"
  >
    <div className="absolute inset-0 bg-gradient-to-br from-primary/5 via-transparent to-secondary/5 rounded-2xl opacity-0 group-hover:opacity-100 transition-opacity duration-300" />
    <div className="relative z-10">
      <div className="flex items-center justify-center w-12 h-12 mb-4 bg-primary/10 rounded-xl group-hover:bg-primary/15 transition-colors duration-300">
        {icon}
      </div>
      <h3 className="text-xl font-semibold mb-3 text-foreground group-hover:text-primary transition-colors duration-300">
        {title}
      </h3>
      <p className="text-muted-foreground leading-relaxed">
        {description}
      </p>
    </div>
  </motion.div>
);

interface SecurityFeatureProps {
  title: string;
  description: string;
  icon: React.ReactNode;
  delay: number;
}

const SecurityFeature = ({ title, description, icon, delay }: SecurityFeatureProps) => (
  <motion.div
    custom={delay}
    variants={fadeUpVariants}
    initial="hidden"
    whileInView="visible"
    viewport={{ once: true, margin: '-100px' }}
    className="flex items-start gap-4 p-4 rounded-lg hover:bg-muted/50 transition-colors duration-200"
  >
    <div className="flex items-center justify-center w-10 h-10 bg-emerald-500/10 rounded-lg flex-shrink-0">
      {icon}
    </div>
    <div>
      <h4 className="font-medium text-foreground mb-1">{title}</h4>
      <p className="text-sm text-muted-foreground leading-relaxed">{description}</p>
    </div>
  </motion.div>
);

export default function Features() {
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
              <div className="flex items-center justify-center w-5 h-5 bg-primary/20 rounded-full">
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
                  className="text-primary"
                >
                  <polygon points="13,2 3,14 12,14 11,22 21,10 12,10 13,2" />
                </svg>
              </div>
              <span className="text-sm text-muted-foreground font-medium">Powerful Features</span>
            </motion.div>

            <motion.h1
              custom={1}
              variants={fadeUpVariants}
              initial="hidden"
              animate="visible"
              className="text-4xl sm:text-5xl md:text-6xl font-bold mb-6 tracking-tight"
            >
              <span className="bg-clip-text text-transparent bg-gradient-to-b from-foreground to-foreground/80">
                Everything You Need for
              </span>
              <br />
              <span className="bg-clip-text text-transparent bg-gradient-to-r from-emerald-500 via-primary to-blue-500">
                Secure Configuration
              </span>
            </motion.h1>

            <motion.p
              custom={2}
              variants={fadeUpVariants}
              initial="hidden"
              animate="visible"
              className="text-lg md:text-xl text-muted-foreground leading-relaxed max-w-3xl mx-auto"
            >
              VaultHub provides enterprise-grade security with developer-friendly tools
              to manage your environment variables, API keys, and sensitive configuration data.
            </motion.p>
          </div>
        </div>
      </section>

      {/* Core Features Grid */}
      <section className="py-20 md:py-32">
        <div className="container mx-auto px-4 md:px-6">
          <motion.div
            variants={fadeInVariants}
            initial="hidden"
            whileInView="visible"
            viewport={{ once: true }}
            className="text-center mb-16"
          >
            <h2 className="text-3xl md:text-4xl font-bold mb-4">Core Features</h2>
            <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
              Built with security, performance, and developer experience in mind
            </p>
          </motion.div>

          <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
            <FeatureCard
              delay={0}
              icon={
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="24"
                  height="24"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  className="text-primary"
                >
                  <rect width="18" height="11" x="3" y="11" rx="2" ry="2" />
                  <path d="m7 11 0-4a5 5 0 0 1 10 0v4" />
                </svg>
              }
              title="AES-256-GCM Encryption"
              description="Military-grade encryption ensures your secrets are protected at rest and in transit with advanced cryptographic standards."
            />

            <FeatureCard
              delay={1}
              icon={
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="24"
                  height="24"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  className="text-primary"
                >
                  <rect width="20" height="14" x="2" y="3" rx="2" ry="2" />
                  <line x1="8" x2="16" y1="21" y2="21" />
                  <line x1="12" x2="12" y1="17" y2="21" />
                </svg>
              }
              title="Web Interface"
              description="Modern, responsive React application with intuitive design for managing your vaults and configuration data."
            />

            <FeatureCard
              delay={2}
              icon={
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="24"
                  height="24"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  className="text-primary"
                >
                  <polyline points="4,17 10,11 4,5" />
                  <line x1="12" x2="20" y1="19" y2="19" />
                </svg>
              }
              title="Command Line Interface"
              description="Powerful CLI tool for automation, CI/CD integration, and seamless workflow integration across platforms."
            />

            <FeatureCard
              delay={3}
              icon={
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="24"
                  height="24"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  className="text-primary"
                >
                  <path d="M9 12l2 2 4-4" />
                  <path d="M21 12c0 1-1 2-2 2h-1c-1 0-2-1-2-2s1-2 2-2h1c1 0 2 1 2 2z" />
                  <path d="M3 12c0-1 1-2 2-2h1c1 0 2 1 2 2s-1 2-2 2H5c-1 0-2-1-2-2z" />
                  <path d="M18 6h1c1 0 2 1 2 2v8c0 1-1 2-2 2h-1" />
                  <path d="M6 6H5c-1 0-2 1-2 2v8c0 1 1 2 2 2h1" />
                </svg>
              }
              title="API Key Authentication"
              description="Secure API keys for programmatic access with fine-grained permissions and access control."
            />

            <FeatureCard
              delay={4}
              icon={
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="24"
                  height="24"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  className="text-primary"
                >
                  <path d="M12 20h9" />
                  <path d="M16.5 3.5a2.12 2.12 0 0 1 3 3L7 19l-4 1 1-4Z" />
                </svg>
              }
              title="Complete Audit Trail"
              description="Comprehensive logging of all operations with detailed timestamps and user tracking for compliance."
            />

            <FeatureCard
              delay={5}
              icon={
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="24"
                  height="24"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  className="text-primary"
                >
                  <path d="M14.5 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7.5L14.5 2z" />
                  <polyline points="14,2 14,8 20,8" />
                  <line x1="16" x2="8" y1="13" y2="13" />
                  <line x1="16" x2="8" y1="17" y2="17" />
                  <polyline points="10,9 9,9 8,9" />
                </svg>
              }
              title="Environment File Export"
              description="Export vault contents as environment files (.env) for seamless integration with your development workflow."
            />
          </div>
        </div>
      </section>

      {/* Security Features */}
      <section className="py-20 md:py-32 bg-muted/20">
        <div className="container mx-auto px-4 md:px-6">
          <div className="grid lg:grid-cols-2 gap-16 items-center">
            <motion.div
              variants={fadeInVariants}
              initial="hidden"
              whileInView="visible"
              viewport={{ once: true }}
            >
              <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-emerald-500/10 border border-emerald-500/20 mb-6">
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
                  className="text-emerald-500"
                >
                  <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10" />
                  <path d="m9 12 2 2 4-4" />
                </svg>
                <span className="text-sm font-medium text-emerald-600 dark:text-emerald-400">Security First</span>
              </div>

              <h2 className="text-3xl md:text-4xl font-bold mb-6">
                Enterprise-Grade Security
              </h2>
              <p className="text-lg text-muted-foreground mb-8 leading-relaxed">
                VaultHub implements multiple layers of security to protect your most sensitive data.
                From encryption at rest to secure authentication mechanisms.
              </p>

              <div className="space-y-4">
                <SecurityFeature
                  delay={0}
                  title="End-to-End Encryption"
                  description="All data is encrypted before storage using AES-256-GCM with unique encryption keys"
                  icon={
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
                      className="text-emerald-500"
                    >
                      <rect width="18" height="11" x="3" y="11" rx="2" ry="2" />
                      <path d="m7 11 0-4a5 5 0 0 1 10 0v4" />
                    </svg>
                  }
                />

                <SecurityFeature
                  delay={1}
                  title="JWT & OIDC Authentication"
                  description="Secure authentication with JWT tokens and optional OIDC integration"
                  icon={
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
                      className="text-emerald-500"
                    >
                      <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2" />
                      <circle cx="9" cy="7" r="4" />
                      <path d="m22 22-5-5" />
                      <path d="m17 22 5-5" />
                    </svg>
                  }
                />

                <SecurityFeature
                  delay={2}
                  title="Access Control"
                  description="Fine-grained permissions and API key-based access control for different vault scopes"
                  icon={
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
                      className="text-emerald-500"
                    >
                      <rect width="18" height="18" x="3" y="3" rx="2" />
                      <path d="M9 9h6v6H9z" />
                    </svg>
                  }
                />
              </div>
            </motion.div>

            <motion.div
              variants={fadeUpVariants}
              initial="hidden"
              whileInView="visible"
              viewport={{ once: true }}
              custom={0}
              className="relative"
            >
              <div className="relative bg-card border border-border rounded-2xl p-8 md:p-12">
                <div className="absolute inset-0 bg-gradient-to-br from-emerald-500/5 via-transparent to-blue-500/5 rounded-2xl" />
                <div className="relative z-10">
                  <div className="text-center mb-8">
                    <div className="inline-flex items-center justify-center w-16 h-16 bg-emerald-500/10 rounded-2xl mb-4">
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        width="32"
                        height="32"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        strokeWidth="2"
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        className="text-emerald-500"
                      >
                        <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10" />
                      </svg>
                    </div>
                    <h3 className="text-xl font-semibold mb-2">Security Stats</h3>
                  </div>

                  <div className="grid grid-cols-2 gap-6">
                    <div className="text-center">
                      <div className="text-2xl font-bold text-emerald-500 mb-1">256-bit</div>
                      <div className="text-sm text-muted-foreground">AES Encryption</div>
                    </div>
                    <div className="text-center">
                      <div className="text-2xl font-bold text-blue-500 mb-1">100%</div>
                      <div className="text-sm text-muted-foreground">Audit Coverage</div>
                    </div>
                    <div className="text-center">
                      <div className="text-2xl font-bold text-violet-500 mb-1">Zero</div>
                      <div className="text-sm text-muted-foreground">Plain Text Storage</div>
                    </div>
                    <div className="text-center">
                      <div className="text-2xl font-bold text-amber-500 mb-1">Multi</div>
                      <div className="text-sm text-muted-foreground">Auth Methods</div>
                    </div>
                  </div>
                </div>
              </div>
            </motion.div>
          </div>
        </div>
      </section>

      {/* Developer Experience */}
      <section className="py-20 md:py-32">
        <div className="container mx-auto px-4 md:px-6">
          <motion.div
            variants={fadeInVariants}
            initial="hidden"
            whileInView="visible"
            viewport={{ once: true }}
            className="text-center mb-16"
          >
            <h2 className="text-3xl md:text-4xl font-bold mb-4">Developer Experience</h2>
            <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
              Built by developers, for developers. Integrate seamlessly with your existing workflow.
            </p>
          </motion.div>

          <div className="grid lg:grid-cols-3 gap-8">
            {/* CLI Integration */}
            <motion.div
              variants={fadeUpVariants}
              initial="hidden"
              whileInView="visible"
              viewport={{ once: true }}
              custom={0}
              className="bg-card border border-border rounded-xl p-6"
            >
              <div className="flex items-center gap-3 mb-4">
                <div className="w-10 h-10 bg-blue-500/10 rounded-lg flex items-center justify-center">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="20"
                    height="20"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    strokeWidth="2"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    className="text-blue-500"
                  >
                    <polyline points="4,17 10,11 4,5" />
                    <line x1="12" x2="20" y1="19" y2="19" />
                  </svg>
                </div>
                <h3 className="text-lg font-semibold">CLI Integration</h3>
              </div>
              <p className="text-sm text-muted-foreground mb-4">
                Powerful command-line interface for automation, CI/CD pipelines, and developer workflows.
              </p>
              <div className="bg-muted/50 rounded-lg p-3 font-mono text-sm mb-4">
                <div className="text-emerald-500">$ vault-hub get --name prod-secrets</div>
                <div className="text-muted-foreground">API_KEY=***</div>
                <div className="text-muted-foreground">DATABASE_URL=***</div>
              </div>
              <div className="space-y-2 text-xs text-muted-foreground">
                <div>• List and manage vaults</div>
                <div>• Execute commands with env vars</div>
                <div>• Cross-platform binaries</div>
              </div>
            </motion.div>

            {/* OpenAPI Integration */}
            <motion.div
              variants={fadeUpVariants}
              initial="hidden"
              whileInView="visible"
              viewport={{ once: true }}
              custom={1}
              className="bg-card border border-border rounded-xl p-6"
            >
              <div className="flex items-center gap-3 mb-4">
                <div className="w-10 h-10 bg-violet-500/10 rounded-lg flex items-center justify-center">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="20"
                    height="20"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    strokeWidth="2"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    className="text-violet-500"
                  >
                    <path d="M8 2v4" />
                    <path d="M16 2v4" />
                    <rect width="18" height="18" x="3" y="4" rx="2" />
                    <path d="M3 10h18" />
                  </svg>
                </div>
                <h3 className="text-lg font-semibold">OpenAPI Integration</h3>
              </div>
              <p className="text-sm text-muted-foreground mb-4">
                Complete OpenAPI 3.0 specification with auto-generated client libraries for seamless integration.
              </p>
              <div className="bg-muted/50 rounded-lg p-3 font-mono text-sm mb-4">
                <div className="text-violet-500">{'// Auto-generated TypeScript client'}</div>
                <div className="text-muted-foreground">const response = await api.getVault({'{'}name{'}'});</div>
                <div className="text-muted-foreground">console.log(response.data.value);</div>
              </div>
              <div className="space-y-2 text-xs text-muted-foreground">
                <div>• TypeScript & Go clients</div>
                <div>• Auto-generated from spec</div>
                <div>• Type-safe API interactions</div>
              </div>
            </motion.div>

            {/* Cross-Platform Support */}
            <motion.div
              variants={fadeUpVariants}
              initial="hidden"
              whileInView="visible"
              viewport={{ once: true }}
              custom={2}
              className="bg-card border border-border rounded-xl p-6"
            >
              <div className="flex items-center gap-3 mb-4">
                <div className="w-10 h-10 bg-emerald-500/10 rounded-lg flex items-center justify-center">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="20"
                    height="20"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    strokeWidth="2"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    className="text-emerald-500"
                  >
                    <circle cx="12" cy="12" r="10" />
                    <path d="M12 2a14.5 14.5 0 0 0 0 20 14.5 14.5 0 0 0 0-20" />
                    <path d="M2 12h20" />
                  </svg>
                </div>
                <h3 className="text-lg font-semibold">Cross-Platform</h3>
              </div>
              <p className="text-sm text-muted-foreground mb-4">
                Native binaries for all major operating systems and architectures.
              </p>
              <div className="grid grid-cols-3 gap-3 mb-4">
                <div className="text-center p-3 rounded-lg bg-muted/30">
                  <div className="w-6 h-6 mx-auto mb-1 bg-blue-500/10 rounded flex items-center justify-center">
                    <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" className="text-blue-500">
                      <rect width="20" height="14" x="2" y="3" rx="2" ry="2" />
                      <line x1="8" x2="16" y1="21" y2="21" />
                      <line x1="12" x2="12" y1="17" y2="21" />
                    </svg>
                  </div>
                  <div className="text-xs font-medium">Linux</div>
                </div>
                <div className="text-center p-3 rounded-lg bg-muted/30">
                  <div className="w-6 h-6 mx-auto mb-1 bg-emerald-500/10 rounded flex items-center justify-center">
                    <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" className="text-emerald-500">
                      <rect width="20" height="14" x="2" y="3" rx="2" ry="2" />
                      <line x1="8" x2="16" y1="21" y2="21" />
                      <line x1="12" x2="12" y1="17" y2="21" />
                    </svg>
                  </div>
                  <div className="text-xs font-medium">Windows</div>
                </div>
                <div className="text-center p-3 rounded-lg bg-muted/30">
                  <div className="w-6 h-6 mx-auto mb-1 bg-amber-500/10 rounded flex items-center justify-center">
                    <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" className="text-amber-500">
                      <rect width="20" height="14" x="2" y="3" rx="2" ry="2" />
                      <line x1="8" x2="16" y1="21" y2="21" />
                      <line x1="12" x2="12" y1="17" y2="21" />
                    </svg>
                  </div>
                  <div className="text-xs font-medium">macOS</div>
                </div>
              </div>
              <div className="space-y-2 text-xs text-muted-foreground">
                <div>• amd64 & arm64 support</div>
                <div>• Statically linked binaries</div>
                <div>• Single binary deployment</div>
              </div>
            </motion.div>
          </div>
        </div>
      </section>
    </div>
  );
}
