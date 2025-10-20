import { TokenInitializer } from './TokenInitializer';

/**
 * AuthProvider wraps the application and ensures tokens are initialized.
 * Tokens are automatically read from cookies by the TokenInitializer.
 */
export function AuthProvider({ children }: { children: React.ReactNode }) {
  return (
    <>
      <TokenInitializer />
      {children}
    </>
  );
}
