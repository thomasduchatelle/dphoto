import axios from "axios";
import {ReactNode, useEffect, useState} from "react";
import {AppConfigContext, AppConfigContextType} from "./app-config.context";

interface ConfigFile {
  googleClientId: string
}

export const AppConfigIntegration = ({children}: {
  children?: ReactNode
}) => {
  const [context, setContext] = useState<AppConfigContextType>({
    googleClientId: "",
    googleLoginUxMode: "popup",
  })

  useEffect(() => {
    axios.get<ConfigFile>("/env-config.json").then(cfg => {
      setContext(ctx => ({
        ...ctx,
        googleClientId: cfg.data.googleClientId,
      }))
    })
  }, [])

  if (!context.googleClientId) {
    return null
  }

  return (
    <AppConfigContext.Provider value={context}>
      {children}
    </AppConfigContext.Provider>
  );
}
