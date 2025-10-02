import {defineConfig, devices} from '@playwright/test';

export default defineConfig({
    /* Reporter to use. See https://playwright.dev/docs/test-reporters */
    reporter: 'html',

    // /* Configure projects for major browsers */
    // projects: [
    //     {
    //         name: 'chromium',
    //         use: {
    //             ...devices['Desktop Chrome'],
    //             // Force consistent viewport for screenshots
    //             viewport: {width: 1280, height: 720},
    //         },
    //     },
    // ],

    /* Run your local dev server before starting the tests */
    webServer: {
        command: 'npm run ladle',
        url: 'http://localhost:61000',
        reuseExistingServer: !process.env.CI,
    },
});
