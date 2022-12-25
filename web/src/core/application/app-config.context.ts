import {createContext, useContext} from "react";

export interface AppConfigContextType {
  googleClientId: string
  googleLoginUxMode: string
}

export const AppConfigContext = createContext<AppConfigContextType>({
  googleClientId: "",
  googleLoginUxMode: "",
})

export function useConfigContext(): AppConfigContextType {
  return useContext<AppConfigContextType>(AppConfigContext);
}