export interface AuthenticatedUser {
  name: string
  email: string
  picture?: string
  accessToken: string
  expiresIn: number
}

export interface UIStatePort {
  // update UI to remove any mention to the authenticated user
  clearUser(): void

  // update UI to display that the user is authenticated
  storeUser(user: AuthenticatedUser): void

  // update UI to notify an error while authenticating
  displayAuthenticationError(error: string): void
}

export interface OAuthService {
  // exchange Google ID token with a DPhoto access token
  authenticateWithGoogleId(googleIdToken: string): Promise<AuthenticatedUser>

  // store access token, to be used for all API calls
  dispatchAccessToken(accessToken: string): void;

  // remove store token(s) that would be used to request DPhoto API
  clearTokens(): void;

  // get direct access to DPhoto access token, could be used to generate download links
  getDPhotoAccessToken(): string
}
