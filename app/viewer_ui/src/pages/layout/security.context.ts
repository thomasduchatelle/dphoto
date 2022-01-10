import {createContext} from "react";
import {AuthenticatedUser} from "../../domain/security";


export interface SecurityContextType {
  user?: AuthenticatedUser
  authenticationError?: string

  signInWithGoogle(): void

  signOut(): void
}

export const SecurityContext = createContext<SecurityContextType>({
  signInWithGoogle: () => {
  },
  signOut: () => {
  },
})