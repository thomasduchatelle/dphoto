export interface OAuthSession {
  sessionId: string;
  originalUrl: string;
  nonce: string;
  codeVerifier: string;
  createdAt: number;
}

export interface SessionStore {
  saveSession(session: OAuthSession): Promise<void>;
  getSession(sessionId: string): Promise<OAuthSession | null>;
  deleteSession(sessionId: string): Promise<void>;
}

class InMemorySessionStore implements SessionStore {
  private sessions: Map<string, OAuthSession> = new Map();
  private readonly TTL_MS = 10 * 60 * 1000; // 10 minutes

  async saveSession(session: OAuthSession): Promise<void> {
    this.sessions.set(session.sessionId, session);
    
    // Auto-cleanup after TTL
    setTimeout(() => {
      this.sessions.delete(session.sessionId);
    }, this.TTL_MS);
  }

  async getSession(sessionId: string): Promise<OAuthSession | null> {
    const session = this.sessions.get(sessionId);
    
    if (!session) {
      return null;
    }

    // Check if session has expired
    if (Date.now() - session.createdAt > this.TTL_MS) {
      this.sessions.delete(sessionId);
      return null;
    }

    return session;
  }

  async deleteSession(sessionId: string): Promise<void> {
    this.sessions.delete(sessionId);
  }
}

// Singleton instance
let sessionStoreInstance: SessionStore | null = null;

export function getSessionStore(): SessionStore {
  if (!sessionStoreInstance) {
    sessionStoreInstance = new InMemorySessionStore();
  }
  return sessionStoreInstance;
}
