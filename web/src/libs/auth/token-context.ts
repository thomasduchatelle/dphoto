// Token context for passing tokens from SSR to client
// Tokens are validated on the server and passed to client in a controlled way

export interface TokenInfo {
  accessToken: string;
  expiresAt: number;
}

let clientAccessToken: TokenInfo | null = null;

export function setClientAccessToken(token: TokenInfo): void {
  clientAccessToken = token;
}

export function getClientAccessToken(): string | null {
  if (!clientAccessToken) {
    return null;
  }

  // Check if token is expired
  if (Date.now() >= clientAccessToken.expiresAt) {
    clientAccessToken = null;
    return null;
  }

  return clientAccessToken.accessToken;
}

export function clearClientAccessToken(): void {
  clientAccessToken = null;
}
