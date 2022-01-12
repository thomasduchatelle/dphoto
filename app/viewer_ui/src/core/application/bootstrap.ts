import {OAuthService} from "../domain/security";
import {OauthServiceImpl} from "../domain/security/adapters/oauthapi/oauth.service";

class AppContext {
  constructor(
    public oauthService: OAuthService
  ) {
  }

  static instance: AppContext

  static getInstance = (): AppContext => {
    if (!instance) {
      const oauthService: OAuthService = new OauthServiceImpl()
      this.instance = new AppContext(oauthService)
    }

    return this.instance
  }
}

let instance: AppContext;

export const getAppContext = AppContext.getInstance
