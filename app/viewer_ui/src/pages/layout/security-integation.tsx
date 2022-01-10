import {ReactNode, useEffect, useMemo, useState} from "react";
import {GoogleLoginResponse, GoogleLoginResponseOffline, useGoogleLogin, useGoogleLogout} from "react-google-login";
import {useLocation, useNavigate} from "react-router-dom";
import googleConfig from "../../config/google.config";
import {AuthenticatedUser, googleSignIn, SecurityDependencies, welcomeAnonymous} from "../../domain/security";
import {SecurityContext, SecurityContextType} from "./security.context";


function isGoogleLoginResponse(value: GoogleLoginResponse | GoogleLoginResponseOffline): value is GoogleLoginResponse {
  return value.hasOwnProperty('profileObj');
}

interface SecurityIntegrationState {
  ready: boolean
  user?: AuthenticatedUser
  authenticationError?: string
}

export default ({loading, children}: {
  loading: ReactNode
  children?: ReactNode
}) => {
  const navigate = useNavigate()
  const location = useLocation()
  const [state, setState] = useState<SecurityIntegrationState>({ready: false})

  useEffect(() => {
    // security domain: dependency injection
    SecurityDependencies.navigationManager = {
      gotoLoginPage(redirectTo: string): void {
        navigate(redirectTo, {replace: true})
      }
    }

    SecurityDependencies.stateManager = {
      clearUser(): void {
        setState({ready: true})
      }, displayAuthenticationError(error: string): void {
        setState({ready: true, authenticationError: error})
      }, storeUser(user: AuthenticatedUser): void {
        setState({ready: true, user})
      }

    }
  }, [navigate])


  const {signIn, loaded} = useGoogleLogin({
    autoLoad: false,
    clientId: googleConfig.clientId,
    onAutoLoadFinished(successLogin: boolean): void {
      console.log(`> onAutoLoadFinished(${successLogin})`);
      if (!successLogin) {
        welcomeAnonymous(location.pathname + location.search)
          .then(() => setState(current => ({...current, ready: true})))
      }
    },
    onFailure(error: any): void {
      console.log(`> onFailure(${JSON.stringify(error)})`)
    },
    onSuccess(response: GoogleLoginResponse | GoogleLoginResponseOffline): void {
      console.log(`> onSuccess(${JSON.stringify(response)})`);
      if (isGoogleLoginResponse(response)) {
        googleSignIn(response.tokenId)
          .then(() => setState(current => ({...current, ready: true})))
      }
    },
    isSignedIn: true,
    uxMode: googleConfig.uxMode,
  });

  const {signOut, loaded: signOutLoaded} = useGoogleLogout({
    clientId: googleConfig.clientId,
    uxMode: googleConfig.uxMode,
  });

  const context = useMemo<SecurityContextType>(() => {
    return {
      user: state.user,
      signInWithGoogle: (): void => {
        signIn()
        console.log("after sign-in!")
      },
      signOut: (): void => {
        signOut()
      },
    }
  }, [state.user])


  const ready = loaded && signOutLoaded && state.ready

  return (
    <SecurityContext.Provider value={context}>
      {ready && (children) || (loading)}
      {/*<GoogleLogin clientId='foo' />*/}
    </SecurityContext.Provider>
  )
}