import {OAuthService, UIStatePort} from "./security.domain";

export class LogoutCase {
  constructor(readonly stateManager: UIStatePort,
              readonly oauthService: OAuthService) {
  }

  public logout(): Promise<void> {
    this.oauthService.clearTokens()
    this.stateManager.clearUser()
    return Promise.resolve()
  }
}