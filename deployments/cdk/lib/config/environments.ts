export interface EnvironmentConfig {
    production: boolean
    // Pairs of access keys used by the CLI
    cliAccessKeys?: string[]
    // Domain registered in Route53
    rootDomain: string
    // (sub)Domain used to expose the root of the application.
    domainName: string
    // (subsub)Domain used for Cognito hosted UI (authentication)
    cognitoDomainName: string
    // Other domains to allow redirection to after login (Cognito Hosted UI), and logout
    cognitoExtraRedirectDomains: string[]
    // Email used for SSL certificate registration automated by let's encrypt
    certificateEmail: string
    // OAuth2 Client ID for Google SSO, used by Cognito
    googleLoginClientId: string
}

export const environments: Record<string, EnvironmentConfig> = {
    live: {
        production: true,
        cliAccessKeys: ['2025-07'],
        rootDomain: 'duchatelle.me',
        domainName: 'dphoto.duchatelle.me',
        cognitoDomainName: 'login.dphoto.duchatelle.me',
        cognitoExtraRedirectDomains: [],
        certificateEmail: 'duchatelle.thomas@gmail.com',
        googleLoginClientId: '841197197570-1o0or8ioo9c4m31405q2h2k8hvdb5enh.apps.googleusercontent.com',
    },
    next: {
        production: false,
        cliAccessKeys: ['2025-07'],
        rootDomain: 'duchatelle.me',
        domainName: 'next.duchatelle.me',
        cognitoDomainName: 'login.next.duchatelle.me',
        cognitoExtraRedirectDomains: ['localhost:3000'],
        certificateEmail: 'duchatelle.thomas@gmail.com',
        googleLoginClientId: '841197197570-7hlq9e86d6u37eoq8nsd8af4aaisl5gb.apps.googleusercontent.com',
    },
    test: {
        production: true,
        cliAccessKeys: ['2024-04'],
        rootDomain: 'exmaple.com',
        domainName: 'dphoto.example.com',
        cognitoDomainName: 'login.dphoto.example.com',
        cognitoExtraRedirectDomains: ["localhost:3210"],
        certificateEmail: 'dphoto@example.com',
        googleLoginClientId: 'test-google-client-id',
    }
};
