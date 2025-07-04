export interface EnvironmentConfig {
    production: boolean
    cliAccessKeys?: string[]
    keybaseUser?: string
    rootDomain: string
    domainName: string
    certificateEmail: string
}

export const environments: Record<string, EnvironmentConfig> = {
    live: {
        production: true,
        cliAccessKeys: ['2025-07'],
        rootDomain: 'duchatelle.me',
        domainName: 'dphoto.duchatelle.me',
        certificateEmail: 'duchatelle.thomas@gmail.com'
    },
    next: {
        production: false,
        cliAccessKeys: ['2025-07'],
        rootDomain: 'duchatelle.me',
        domainName: 'nextcdk.duchatelle.me',
        certificateEmail: 'duchatelle.thomas@gmail.com'
    },
    test: {
        production: true,
        cliAccessKeys: ['2024-04'],
        rootDomain: 'exmaple.com',
        domainName: 'dphoto.example.com',
        certificateEmail: 'dphoto@example.com'
    }
};