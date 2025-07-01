export interface EnvironmentConfig {
    production: boolean;
    cliAccessKeys?: string[];
    keybaseUser?: string;
}

export const environments: Record<string, EnvironmentConfig> = {
    dev: {
        production: false,
        cliAccessKeys: ['2025-07'],
        keybaseUser: 'keybase:thomasduchatelle'
    },
    live: {
        production: true,
        cliAccessKeys: ['2025-07'],
        keybaseUser: 'keybase:thomasduchatelle'
    },
    next: {
        production: false,
        cliAccessKeys: ['2025-07'],
        keybaseUser: 'keybase:thomasduchatelle'
    },
    test: {
        production: true,
        cliAccessKeys: ['2024-04'],
        keybaseUser: 'keybase:thomasduchatelle'
    }
};
