import {getOidcConfigFromEnv, oidcConfig} from './oidc-config';
import {redirectUrl} from "@/libs/requests";

/**
 * Generates the Cognito logout URL with the logout_uri parameter
 */
export async function getLogoutUrl(): Promise<string> {
    const oidcEnvConfig = getOidcConfigFromEnv();
    const config = await oidcConfig(oidcEnvConfig);

    const logoutUri = await redirectUrl("/auth/logout");

    return `${config.serverMetadata().end_session_endpoint}?client_id=${oidcEnvConfig.clientId}&logout_uri=${encodeURIComponent(logoutUri.toString())}`;
}
