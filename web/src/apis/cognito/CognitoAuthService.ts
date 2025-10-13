import {Issuer, Client, generators, TokenSet} from 'openid-client';

export interface CognitoConfig {
    region: string;
    userPoolId: string;
    clientId: string;
    clientSecret: string;
    redirectUri: string;
    domain: string;
}

export interface AuthSession {
    sessionId: string;
    originalUrl: string;
    nonce: string;
    codeVerifier: string;
    state: string;
}

export class CognitoAuthService {
    private client: Client | null = null;
    private issuer: Issuer<Client> | null = null;
    private config: CognitoConfig;

    constructor(config: CognitoConfig) {
        this.config = config;
    }

    async initialize(): Promise<void> {
        if (this.client) {
            return;
        }

        const issuerUrl = `https://cognito-idp.${this.config.region}.amazonaws.com/${this.config.userPoolId}`;
        this.issuer = await Issuer.discover(issuerUrl);

        this.client = new this.issuer.Client({
            client_id: this.config.clientId,
            client_secret: this.config.clientSecret,
            redirect_uris: [this.config.redirectUri],
            response_types: ['code'],
        });
    }

    async getAuthorizationUrl(originalUrl: string = '/'): Promise<{ url: string; session: AuthSession }> {
        await this.initialize();
        if (!this.client) {
            throw new Error('Client not initialized');
        }

        const codeVerifier = generators.codeVerifier();
        const codeChallenge = generators.codeChallenge(codeVerifier);
        const state = generators.state();
        const nonce = generators.nonce();

        const sessionId = generators.random();

        const authUrl = this.client.authorizationUrl({
            scope: 'openid email profile',
            code_challenge: codeChallenge,
            code_challenge_method: 'S256',
            state: state,
            nonce: nonce,
        });

        return {
            url: authUrl,
            session: {
                sessionId,
                originalUrl,
                nonce,
                codeVerifier,
                state,
            },
        };
    }

    async handleCallback(params: URLSearchParams, session: AuthSession): Promise<TokenSet> {
        await this.initialize();
        if (!this.client) {
            throw new Error('Client not initialized');
        }

        // Verify state
        const state = params.get('state');
        if (state !== session.state) {
            throw new Error('State mismatch');
        }

        const tokenSet = await this.client.callback(
            this.config.redirectUri,
            Object.fromEntries(params),
            {
                code_verifier: session.codeVerifier,
                state: session.state,
                nonce: session.nonce,
            }
        );

        return tokenSet;
    }

    async refreshTokens(refreshToken: string): Promise<TokenSet> {
        await this.initialize();
        if (!this.client) {
            throw new Error('Client not initialized');
        }

        const tokenSet = await this.client.refresh(refreshToken);
        return tokenSet;
    }

    async validateToken(accessToken: string): Promise<boolean> {
        await this.initialize();
        if (!this.client || !this.issuer) {
            throw new Error('Client not initialized');
        }

        try {
            const userinfo = await this.client.userinfo(accessToken);
            return !!userinfo;
        } catch (err) {
            return false;
        }
    }
}

export function getCognitoConfig(): CognitoConfig {
    // These would typically come from environment variables or a config endpoint
    return {
        region: process.env.COGNITO_REGION || 'us-east-1',
        userPoolId: process.env.COGNITO_USER_POOL_ID || '',
        clientId: process.env.COGNITO_CLIENT_ID || '',
        clientSecret: process.env.COGNITO_CLIENT_SECRET || '',
        redirectUri: process.env.COGNITO_REDIRECT_URI || `${typeof window !== 'undefined' ? window.location.origin : ''}/auth/callback`,
        domain: process.env.COGNITO_DOMAIN || '',
    };
}
