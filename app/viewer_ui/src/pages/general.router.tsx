import {useGoogleSignInCase, useSecurityContext} from '../core/application'
import AuthenticatedRouter from "./authenticated.router";
import LoginPage from "./login/login.page";

export default () => {
  const {user, authenticationError} = useSecurityContext();
  const googleSignInCase = useGoogleSignInCase()


  return user ? (
    <AuthenticatedRouter/>
  ) : (
    <LoginPage googleSignIn={googleSignInCase.googleSignIn} authenticationError={authenticationError}/>
  )
}