/** googleSignIn use an identity token from Google to sign in DPhoto*/
import {SecurityDependencies} from "./security.domain";

export const googleSignIn = (identityToken: string): Promise<void> => {
  console.log(`> googleSignIn(${identityToken})`)
  if (!SecurityDependencies.oauthService || !SecurityDependencies.stateManager) {
    return Promise.reject("oauthService and stateManager must be injected.")
  }

  return SecurityDependencies.oauthService.authenticateWithGoogleId(identityToken)
    .then(user => SecurityDependencies.stateManager?.storeUser(user))
    .catch(err => SecurityDependencies.stateManager?.displayAuthenticationError(err))
}