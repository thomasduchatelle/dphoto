'use client';

import { useRouter } from 'waku';
import { ReactNode, useEffect, useState } from 'react';

export interface RouterContextValue {
  path: string;
  params: Record<string, string>;
  query: URLSearchParams;
  navigate: (path: string) => void;
  replace: (path: string) => void;
}

export function useClientRouter(): RouterContextValue {
  const router = useRouter();
  const [path, setPath] = useState('/');
  const [params, setParams] = useState<Record<string, string>>({});
  const [query, setQuery] = useState<URLSearchParams>(new URLSearchParams());

  useEffect(() => {
    if (typeof window !== 'undefined') {
      const updateRoute = () => {
        const currentPath = window.location.pathname;
        const currentSearch = new URLSearchParams(window.location.search);
        setPath(currentPath);
        setQuery(currentSearch);
        
        // Parse params from path
        const parsedParams = parseParams(currentPath);
        setParams(parsedParams);
      };

      updateRoute();

      // Listen for popstate events (back/forward)
      window.addEventListener('popstate', updateRoute);
      
      return () => {
        window.removeEventListener('popstate', updateRoute);
      };
    }
  }, []);

  const navigate = (newPath: string) => {
    router.push(newPath);
    const parsedParams = parseParams(newPath);
    setParams(parsedParams);
    setPath(newPath.split('?')[0]);
    const searchPart = newPath.split('?')[1];
    setQuery(new URLSearchParams(searchPart || ''));
  };

  const replace = (newPath: string) => {
    router.replace(newPath);
    const parsedParams = parseParams(newPath);
    setParams(parsedParams);
    setPath(newPath.split('?')[0]);
    const searchPart = newPath.split('?')[1];
    setQuery(new URLSearchParams(searchPart || ''));
  };

  return {
    path,
    params,
    query,
    navigate,
    replace,
  };
}

function parseParams(path: string): Record<string, string> {
  const cleanPath = path.split('?')[0];
  const parts = cleanPath.split('/').filter(p => p);
  
  // Match /albums/:owner/:album/:encodedId/:filename
  if (parts[0] === 'albums') {
    if (parts.length >= 5) {
      return {
        owner: parts[1],
        album: parts[2],
        encodedId: parts[3],
        filename: parts[4],
      };
    } else if (parts.length >= 3) {
      return {
        owner: parts[1],
        album: parts[2],
      };
    }
  }
  
  return {};
}

export function ClientLink({ to, children, className, onClick }: { 
  to: string; 
  children: ReactNode; 
  className?: string;
  onClick?: (e: React.MouseEvent) => void;
}) {
  const { navigate } = useClientRouter();

  const handleClick = (e: React.MouseEvent<HTMLAnchorElement>) => {
    e.preventDefault();
    if (onClick) {
      onClick(e);
    }
    navigate(to);
  };

  return (
    <a href={to} onClick={handleClick} className={className}>
      {children}
    </a>
  );
}
