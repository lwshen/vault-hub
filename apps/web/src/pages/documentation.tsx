'use client';

import { motion } from 'framer-motion';
import type { Variants } from 'framer-motion';
import { useState } from 'react';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';

// Import markdown content
import cliGuideContent from '@/docs/cli-guide.md?raw';
import serverSetupContent from '@/docs/server-setup.md?raw';
import apiReferenceContent from '@/docs/api-reference.md?raw';
import securityContent from '@/docs/security.md?raw';

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

interface NavItem {
  id: string;
  title: string;
  content: string;
}

const navigationItems: NavItem[] = [
  {
    id: 'cli-guide',
    title: 'CLI Guide',
    content: cliGuideContent,
  },
  {
    id: 'server-setup',
    title: 'Server Setup',
    content: serverSetupContent,
  },
  {
    id: 'api-reference',
    title: 'API Reference',
    content: apiReferenceContent,
  },
  {
    id: 'security',
    title: 'Security',
    content: securityContent,
  },
];

export default function Documentation() {
  const [activeSection, setActiveSection] = useState('cli-guide');

  const currentContent = navigationItems.find(item => item.id === activeSection)?.content || cliGuideContent;

  const handleSectionChange = (sectionId: string) => {
    setActiveSection(sectionId);
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
              <nav className="space-y-2">
                {navigationItems.map((section) => (
                  <button
                    key={section.id}
                    type="button"
                    onClick={() => handleSectionChange(section.id)}
                    className={`block w-full text-left font-medium text-sm transition-colors px-3 py-2 rounded-lg ${
                      activeSection === section.id
                        ? 'text-primary bg-primary/10'
                        : 'text-foreground hover:text-primary hover:bg-muted/50'
                    }`}
                  >
                    {section.title}
                  </button>
                ))}
              </nav>
            </div>
          </motion.aside>

          {/* Main Content */}
          <motion.main
            key={activeSection}
            variants={fadeUpVariants}
            initial="hidden"
            animate="visible"
            custom={0}
            className="flex-1 max-w-4xl"
          >
            <div className="prose prose-lg max-w-none dark:prose-invert">
              <ReactMarkdown remarkPlugins={[remarkGfm]}>
                {currentContent}
              </ReactMarkdown>
            </div>
          </motion.main>
        </div>
      </div>
    </div>
  );
}
