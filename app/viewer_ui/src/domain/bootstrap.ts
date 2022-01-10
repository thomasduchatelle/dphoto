import {SecurityDependencies} from "./security";
import {OauthServiceImpl} from "./security/adalpters/oauthapi/oauth.service";

export function bootstrap() {
  SecurityDependencies.oauthService = OauthServiceImpl
}