import React, {ReactNode, useState} from "react";
import {GoogleLoginResponse, GoogleLoginResponseOffline} from "react-google-login";
import googleConfig from '../config/google.config'
import {authenticateWithGoogle, logoutFromGoogle} from "./google-authentication.service";
import LoginPage from "./login.page";
import LogoutButton from "./logout-button";
import {SecurityReactContext} from "./security.context";
import {SecurityContextType} from "./security.model";

const SecuredContent = ({children}: {
  children: ReactNode
}) => {

  const [err, setErr] = useState<string>('')
  const [securityContext, setSecurityContext] = useState<SecurityContextType>({})

  const onPopupFailure = ({error}: any) => {
    console.log(`Google popup didn't complete: ${error}`)
  }

  const loggedUser = securityContext.loggedUser

  const authenticate = (googleAnswer: (GoogleLoginResponse | GoogleLoginResponseOffline)) => {
    if (googleAnswer) {
      authenticateWithGoogle(googleAnswer)
        .then(newContext => {
          setSecurityContext(newContext)
        })
        .catch(err => setErr(err))
    }
  }

  const logout = () => {
    logoutFromGoogle().then(ctx => setSecurityContext(ctx))
  }

  return (
    <SecurityReactContext.Provider value={securityContext}>
      {(loggedUser && loggedUser.email) ?
        (
          <div style={{display: 'absolute', right: '20px'}}>
            <LogoutButton clientId={googleConfig.clientId} onLogoutSuccess={logout}/>
            {children}
          </div>
        )
        : (
          <LoginPage clientId={googleConfig.clientId}
                     isSignedIn={true}
                     authenticateWithGoogle={authenticate}
                     onFailure={onPopupFailure}
                     errorMessage={err}/>
        )
      }
    </SecurityReactContext.Provider>
  );
}

export default SecuredContent