import {ReactNode, useMemo, useState} from "react";
import {SecurityContext, SecurityContextType} from "./security.context";

export const SecurityIntegration = ({children}: {
  children?: ReactNode
}) => {
  type State = Omit<SecurityContextType, "mutateContext">
  const [state, setState] = useState<State>({})

  const mutateContext = useMemo(() => (mutator: (current: State) => State) => {
    setState(current => ({
      ...mutator(current),
      mutateContext: context.mutateContext,
    }))
  }, [setState])

  const context = {
    ...state,
    mutateContext
  }

  return (
    <SecurityContext.Provider value={context}>
      {children}
    </SecurityContext.Provider>
  )
}
