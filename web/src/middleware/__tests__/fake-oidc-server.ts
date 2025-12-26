import {http, HttpResponse} from 'msw';
import {setupServer} from 'msnow/node';

interface TokenResponse {
    access_token: string;
    refresh_token?: string;
    id_token?: string;
    expires_in: number;
    token_type: string;
}

interface TokenError {
    error: string;
    error_description?: string;
}

export class FakeOIDCServer {
    private server: ReturnType<typeof setupServer>;
    private issuerUrl: string;
    private tokenResponses = new Map<string, TokenResponse | TokenError>();
    private clientId: string;
    private clientSecret: string;

    constructor(issuerUrl: string, clientId: string, clientSecret: string) {
        this.issuerUrl = issuerUrl;
        this.clientId = clientId;
        this.clientSecret = clientSecret;

        this.server = setupServer(
            http.get(`${issuerUrl}/.well-known/openid-configuration`, () => {
                return HttpResponse.json({
                    issuer: issuerUrl,
                    authorization_endpoint: `${issuerUrl}/oauth2/authorize`,
                    token_endpoint: `${issuerUrl}/oauth2/token`,
                    jwks_uri: `${issuerUrl}/.well-known/jwks.json`,
                    response_types_supported: ['code'],
                    subject_types_supported: ['public'],
                    id_token_signing_alg_values_supported: ['RS256'],
                });
            }),

            http.post(`${issuerUrl}/oauth2/token`, async ({request}) => {
                const body = await request.text();
                const params = new URLSearchParams(body);
                const code = params.get('code');

                if (!code) {
                    return HttpResponse.json(
                        {error: 'invalid_request', error_description: 'Missing authorization code'},
                        {status: 400}
                    );
                }

                const response = this.tokenResponses.get(code);
                if (!response) {
                    return HttpResponse.json(
                        {error: 'invalid_grant', error_description: 'Invalid authorization code'},
                        {status: 400}
                    );
                }

                if ('error' in response) {
                    return HttpResponse.json(response, {status: 400});
                }

                return HttpResponse.json(response);
            }),
        );
    }

    start(): void {
        this.server.listen({onUnhandledRequest: 'bypass'});
    }

    stop(): void {
        this.server.close();
    }

    reset(): void {
        this.tokenResponses.clear();
    }

    setupSuccessfulTokenExchange(code: string, tokens: TokenResponse): void {
        this.tokenResponses.set(code, tokens);
    }

    setupTokenError(code: string, error: string, errorDescription?: string): void {
        this.tokenResponses.set(code, {
            error,
            error_description: errorDescription,
        });
    }
}

