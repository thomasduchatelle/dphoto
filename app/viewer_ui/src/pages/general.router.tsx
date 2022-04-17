import {useGoogleSignInCase, useSecurityContext} from '../core/application'
import AuthenticatedRouter from "./authenticated/authenticated.router";
import LoginPage from "./login";

const GeneralRouter = () => {
  const {user, authenticationError} = useSecurityContext();
  const googleSignInCase = useGoogleSignInCase()

  return user ? (
    <AuthenticatedRouter/>
  ) : (
    <LoginPage googleSignIn={googleSignInCase.googleSignIn} authenticationError={authenticationError}/>
  )
}

export default GeneralRouter
