import {OAuthService} from "../domain/security";
import {OauthServiceImpl} from "../domain/security/adapters/oauthapi/oauth.service";

export interface AppContext {
  oauthService: OAuthService
}

let instance: AppContext;

export function getAppContext() {
  if (!instance) {
    instance = bootstrap()
  }
  return instance
}

export function bootstrap(): AppContext {
  const oauthService: OAuthService = new OauthServiceImpl()

  return {
    oauthService,
  }
}