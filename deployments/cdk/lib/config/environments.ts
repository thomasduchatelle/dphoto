export interface EnvironmentConfig {
    production: boolean
    cliAccessKeys?: string[]
    keybaseUser?: string
    rootDomain: string
    domainName: string
    certificateEmail: string
    googleLoginClientId: string
}

export const environments: Record<string, EnvironmentConfig> = {
    live: {
        production: true,
        cliAccessKeys: ['2025-07'],
        rootDomain: 'duchatelle.me',
        domainName: 'dphoto.duchatelle.me',
        certificateEmail: 'duchatelle.thomas@gmail.com',
        googleLoginClientId: '841197197570-e78pmti86rg3gjrl03cd93th4tuiml8a.apps.googleusercontent.com'
    },
    next: {
        production: false,
        cliAccessKeys: ['2025-07'],
        rootDomain: 'duchatelle.me',
        domainName: 'nextcdk.duchatelle.me',
        certificateEmail: 'duchatelle.thomas@gmail.com',
        googleLoginClientId: '841197197570-7hlq9e86d6u37eoq8nsd8af4aaisl5gb.apps.googleusercontent.com'
    },
    test: {
        production: true,
        cliAccessKeys: ['2024-04'],
        rootDomain: 'exmaple.com',
        domainName: 'dphoto.example.com',
        certificateEmail: 'dphoto@example.com',
        googleLoginClientId: 'test-google-client-id'
    }
};
