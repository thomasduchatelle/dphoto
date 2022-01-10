export interface AuthenticatedUser {
  name: string
  email: string
  picture?: string
  accessToken: string
}

export interface NavigationPort {
  gotoLoginPage(redirectPath: string): void
}

export interface StatePort {
  clearUser(): void

  storeUser(user: AuthenticatedUser): void

  displayAuthenticationError(error: string): void
}

export interface OAuthService {
  authenticateWithGoogleId(googleIdToken: string): Promise<AuthenticatedUser>
}

class SecurityDependenciesClass {
  constructor(public navigationManager?: NavigationPort,
              public stateManager?: StatePort,
              public oauthService?: OAuthService) {
  }
}

export const SecurityDependencies = new SecurityDependenciesClass()
