import {createContext, useContext} from "react";
import {AuthenticatedUser, GoogleSignInCase, LogoutCase} from "../domain/security";
import {getAppContext} from "./bootstrap";


export interface SecurityContextType {
  user?: AuthenticatedUser
  authenticationError?: string

  mutateContext(mutator: (current: Omit<SecurityContextType, "mutateContext">) => Omit<SecurityContextType, "mutateContext">): void
}

export const SecurityContext = createContext<SecurityContextType>({
  mutateContext(mutator: (current: SecurityContextType) => SecurityContextType) {
  }
})

export function useSecurityContext(): Omit<SecurityContextType, "mutateContext"> {
  return useContext(SecurityContext)
}

export function useAuthenticatedUser(): AuthenticatedUser | undefined {
  return useContext(SecurityContext).user
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

  return new LogoutCase(newStateManager(securityContext), getAppContext().oauthService)
}