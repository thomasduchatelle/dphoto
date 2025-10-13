export interface EnvironmentConfig {
    production: boolean
    cliAccessKeys?: string[]
    keybaseUser?: string
    rootDomain: string
    domainName: string
    certificateEmail: string
    googleLoginClientId: string
    googleClientSecret: string
}

export const environments: Record<string, EnvironmentConfig> = {
    live: {
        production: true,
        cliAccessKeys: ['2025-07'],
        rootDomain: 'duchatelle.me',
        domainName: 'dphoto.duchatelle.me',
        certificateEmail: 'duchatelle.thomas@gmail.com',
        googleLoginClientId: '841197197570-1o0or8ioo9c4m31405q2h2k8hvdb5enh.apps.googleusercontent.com',
        googleClientSecret: process.env.GOOGLE_CLIENT_SECRET_LIVE || ''
    },
    next: {
        production: false,
        cliAccessKeys: ['2025-07'],
        rootDomain: 'duchatelle.me',
        domainName: 'next.duchatelle.me',
        certificateEmail: 'duchatelle.thomas@gmail.com',
        googleLoginClientId: '841197197570-7hlq9e86d6u37eoq8nsd8af4aaisl5gb.apps.googleusercontent.com',
        googleClientSecret: process.env.GOOGLE_CLIENT_SECRET_NEXT || ''
    },
    test: {
        production: true,
        cliAccessKeys: ['2024-04'],
        rootDomain: 'exmaple.com',
        domainName: 'dphoto.example.com',
        certificateEmail: 'dphoto@example.com',
        googleLoginClientId: 'test-google-client-id',
        googleClientSecret: 'test-google-client-secret'
    }
};
