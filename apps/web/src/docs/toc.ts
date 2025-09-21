// Table of Contents configuration for documentation
import cliGuideContent from './cli-guide.md?raw';
import serverSetupContent from './server-setup.md?raw';
import apiReferenceContent from './api-reference.md?raw';
import securityContent from './security.md?raw';

export interface TOCItem {
  id: string;
  title: string;
  description?: string;
  content: string;
}

export const documentationTOC: TOCItem[] = [
  {
    id: 'cli-guide',
    title: 'CLI Guide',
    description: 'Get started with the VaultHub CLI',
    content: cliGuideContent,
  },
  {
    id: 'server-setup',
    title: 'Server Setup',
    description: 'Install and configure VaultHub server',
    content: serverSetupContent,
  },
  {
    id: 'api-reference',
    title: 'API Reference',
    description: 'Complete API documentation',
    content: apiReferenceContent,
  },
  {
    id: 'security',
    title: 'Security',
    description: 'Security features and best practices',
    content: securityContent,
  },
];

// Helper functions
export const getDocumentationItem = (id: string): TOCItem | undefined => {
  return documentationTOC.find(item => item.id === id);
};

export const getDefaultDocumentation = (): TOCItem => {
  return documentationTOC[0]; // Return first item (CLI Guide)
};
