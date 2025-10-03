import {defineConfig} from '@playwright/test';

export default defineConfig({
    testDir: './playwright',
    snapshotPathTemplate: `{testDir}/{testFilePath}-${process.env.CI ? 'snapshots' : 'local'}/{arg}-{platform}{ext}`,
    workers: process.env.CI ? 2 : undefined,
    fullyParallel: true,

    /* Run your local dev server before starting the tests */
    webServer: {
        command: 'npm run ladle',
        url: 'http://localhost:61000',
        reuseExistingServer: !process.env.CI,
    },
});
