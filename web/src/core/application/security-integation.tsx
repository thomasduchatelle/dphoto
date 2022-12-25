import {ReactNode, useState} from "react";
import {SecurityContext, SecurityContextPayloadType, SecurityContextType} from "./security.context";

export const SecurityIntegration = ({children}: {
  children?: ReactNode
}) => {
  const [state, setState] = useState<SecurityContextType>({
    payload: {},
    mutateContext(mutator: (current: SecurityContextPayloadType) => SecurityContextPayloadType) {
      setState(current => ({
        payload: mutator(current.payload),
        mutateContext: current.mutateContext
      }))
    }
  })

  return (
    <SecurityContext.Provider value={state}>
      {children}
    </SecurityContext.Provider>
  )
}
