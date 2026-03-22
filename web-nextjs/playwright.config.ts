import {defineConfig} from '@playwright/test';

export default defineConfig({
    testDir: './playwright',
    snapshotPathTemplate: `{testDir}/{testFilePath}-${process.env.CI ? 'snapshots' : 'local'}/{arg}-{platform}{ext}`,
    workers: process.env.CI ? 2 : undefined,
    fullyParallel: true,

    webServer: {
        command: 'npm run storybook',
        url: 'http://localhost:6006',
        reuseExistingServer: !process.env.CI,
    },
});
