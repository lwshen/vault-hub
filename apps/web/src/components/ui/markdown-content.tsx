import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import rehypeHighlight from 'rehype-highlight';
import { cn } from '@/lib/utils';

interface MarkdownContentProps {
  children: string;
  className?: string;
  size?: 'sm' | 'base' | 'lg' | 'xl';
}

const sizeClasses = {
  sm: 'prose prose-sm',
  base: 'prose',
  lg: 'prose prose-lg',
  xl: 'prose prose-xl',
};

export function MarkdownContent({ 
  children, 
  className,
  size = 'lg' 
}: MarkdownContentProps) {
  const proseClasses = cn(
    // Base prose classes
    sizeClasses[size],
    'max-w-none',
    'dark:prose-invert',
    
    // Custom styling can be added here if needed in the future
    className
  );

  return (
    <div className={proseClasses}>
      <ReactMarkdown
        remarkPlugins={[remarkGfm]}
        rehypePlugins={[rehypeHighlight]}
      >
        {children}
      </ReactMarkdown>
    </div>
  );
}