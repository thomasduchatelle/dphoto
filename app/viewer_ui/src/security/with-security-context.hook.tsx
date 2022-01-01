import {ComponentType, ReactNode} from "react";
import {SecurityReactContext} from "./security.context";
import {SecurityContextType} from "./security.model";

export function withSecurityContext<T>(
  component: (t: T) => ReactNode,
): ComponentType<Omit<T, keyof SecurityContextType>> {

  return properties => (
    <SecurityReactContext.Consumer>
      {(securityContext: SecurityContextType) => {
        const props = {
          ...securityContext,
          ...properties,
        } as unknown as T

        const Component = component as ComponentType<T>
        return (
          <Component {...props} />
        )
      }}
    </SecurityReactContext.Consumer>
  )
}