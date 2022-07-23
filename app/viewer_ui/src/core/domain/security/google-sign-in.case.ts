import {AxiosError} from "axios";
import {OAuthService, UIStatePort} from "./security.domain";

interface ErrorBody {
  code: string
  error: string
}

export let authenticationTimeoutIds: NodeJS.Timeout[] = []

function lookupErrorMessage(code?: string, message?: string): string {
  switch (code) {
    case 'oauth.user-not-preregistered':
      return "You must be pre-registered to use this service."

    default:
      return message ?? "Sorry, we're unable to authenticate you using Google."
  }
}

/** googleSignIn use an identity token from Google to sign in DPhoto*/
export class GoogleSignInCase {
  constructor(readonly stateManager: UIStatePort,
              readonly oauthService: OAuthService) {
  }

  public googleSignIn = (identityToken: string): Promise<void> => {
    return this.oauthService.authenticateWithGoogleId(identityToken)
      .then(user => {
        const timeoutId = setTimeout(() => {
          authenticationTimeoutIds = authenticationTimeoutIds.filter(id => id !== timeoutId)
          this.googleSignIn(identityToken).then()

        }, (user.expiresIn - 60) * 1000)

        authenticationTimeoutIds.push(timeoutId)

        this.oauthService.dispatchAccessToken(user.accessToken)
        this.stateManager.storeUser(user)
      })
      .catch((err: AxiosError<ErrorBody>) => {
        console.log(`Authentication failed: ${JSON.stringify(err)}`)
        this.stateManager.displayAuthenticationError(lookupErrorMessage(err.response?.data?.code, err.response?.data?.error))
      })
  }
}
