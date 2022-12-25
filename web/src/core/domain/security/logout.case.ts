import {authenticationTimeoutIds} from "./google-sign-in.case";
import {OAuthService, UIStatePort} from "./security.domain";

export class LogoutCase {
  constructor(readonly stateManager: UIStatePort,
              readonly oauthService: OAuthService) {
  }

  public logout = (): Promise<void> => {
    authenticationTimeoutIds.forEach(timeoutId => {
      clearTimeout(timeoutId)
    })
    authenticationTimeoutIds.splice(0)

    this.oauthService.clearTokens()
    this.stateManager.clearUser()
    return Promise.resolve()
  }
}