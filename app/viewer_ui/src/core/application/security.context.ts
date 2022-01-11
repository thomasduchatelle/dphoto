import {createContext, useContext} from "react";
import {AuthenticatedUser, GoogleSignInCase} from "../domain/security";
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

export function useGoogleSignInCase(): GoogleSignInCase {
  const securityContext = useContext(SecurityContext);

  return new GoogleSignInCase({
      clearUser(): void {
        securityContext.mutateContext(_ => ({}))
      },
      displayAuthenticationError(authenticationError: string): void {
        securityContext.mutateContext(current => ({...current, authenticationError}))
      },
      storeUser(user: AuthenticatedUser): void {
        securityContext.mutateContext(_ => ({user, authenticationError: undefined}))
      }
    },
    getAppContext().oauthService)
}