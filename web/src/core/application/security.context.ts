import {AxiosInstance} from "axios";
import {createContext, useContext, useMemo} from "react";
import {AuthenticatedUser, GoogleSignInCase, LogoutCase} from "../domain/security";
import {authenticatedAxios} from "../domain/security/adapters/oauthapi/oauth.service";
import {getAppContext} from "./bootstrap";

export interface SecurityContextPayloadType {
  user?: AuthenticatedUser
  authenticationError?: string
}

export interface SecurityContextType {
  payload: SecurityContextPayloadType

  mutateContext(mutator: (current: SecurityContextPayloadType) => SecurityContextPayloadType): void
}

export interface MustBeAuthenticated {
  loggedUser: AuthenticatedUser
  signOutCase: LogoutCase
  authenticatedAxios: AxiosInstance
  accessToken: string
}

export const SecurityContext = createContext<SecurityContextType>({
  payload: {},
  mutateContext(mutator: (current: SecurityContextPayloadType) => SecurityContextPayloadType) {
  }
})

export function useSecurityContext(): SecurityContextPayloadType {
  return useContext(SecurityContext).payload
}

export function useAuthenticatedUser(): AuthenticatedUser | undefined {
  return useContext(SecurityContext).payload.user
}

export function useMustBeAuthenticated(): MustBeAuthenticated {
  const user = useAuthenticatedUser();
  if (!user) {
    throw new Error("user must be authenticated to access this page")
  }

  const signOutCase = useSignOutCase()
  return useMemo(() => ({
    loggedUser: user,
    signOutCase: signOutCase,
    authenticatedAxios: authenticatedAxios(),
    accessToken: getAppContext().oauthService.getDPhotoAccessToken(),
  }), [user, signOutCase])
}

function newStateManager(securityContext: SecurityContextType) {
  return {
    clearUser(): void {
      securityContext.mutateContext(_ => ({}))
    },
    displayAuthenticationError(authenticationError: string): void {
      securityContext.mutateContext(current => ({...current, authenticationError}))
    },
    storeUser(user: AuthenticatedUser): void {
      securityContext.mutateContext(_ => ({user, authenticationError: undefined}))
    }
  };
}

export function useGoogleSignInCase(): GoogleSignInCase {
  const securityContext = useContext(SecurityContext);

  return new GoogleSignInCase(newStateManager(securityContext), getAppContext().oauthService)
}

export function useSignOutCase() {
  const securityContext = useContext(SecurityContext);

  return useMemo(() => new LogoutCase(newStateManager(securityContext), getAppContext().oauthService), [securityContext])
}