import AppNavComponent from "../components/app-nav.component";
import UserMenu from "../components/user.menu";
import {useGoogleSignInCase, useSecurityContext, useSignOutCase} from '../core/application'
import AuthenticatedRouter from "./authenticated.router";
import LoginPage from "./login";

export default () => {
  const {user, authenticationError} = useSecurityContext();
  const googleSignInCase = useGoogleSignInCase()
  const signOutCase = useSignOutCase()

  return user ? (
    <>
      <AppNavComponent
        rightContent={<UserMenu user={user} onLogout={signOutCase.logout}/>}
      />
      <AuthenticatedRouter/>
    </>
  ) : (
    <LoginPage googleSignIn={googleSignInCase.googleSignIn} authenticationError={authenticationError}/>
  )
}
