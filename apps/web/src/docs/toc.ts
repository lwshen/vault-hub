// Table of Contents configuration for documentation
import cliGuideContent from './cli-guide.md?raw';
import serverSetupContent from './server-setup.md?raw';
import apiReferenceContent from './api-reference.md?raw';
import securityContent from './security.md?raw';

export interface TOCItem {
  id: string;
  title: string;
  content: string;
}

export const documentationTOC: TOCItem[] = [
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

// Helper functions
export const getDocumentationItem = (id: string): TOCItem | undefined => {
  return documentationTOC.find(item => item.id === id);
};

export const getDefaultDocumentation = (): TOCItem => {
  return documentationTOC[0]; // Return first item (CLI Guide)
};
