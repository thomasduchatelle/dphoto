import {defineConfig, devices} from '@playwright/test';

export default defineConfig({
    testDir: './src',
    /* Run tests in files in parallel */
    fullyParallel: true,
    /* Fail the build on CI if you accidentally left test.only in the source code. */
    forbidOnly: !!process.env.CI,
    /* Retry on CI only */
    retries: process.env.CI ? 2 : 0,
    /* Opt out of parallel tests on CI. */
    workers: process.env.CI ? 1 : 4,
    /* Reporter to use. See https://playwright.dev/docs/test-reporters */
    reporter: 'html',
    /* Shared settings for all the projects below. See https://playwright.dev/docs/api/class-testoptions. */
    use: {
        /* Collect trace when retrying the failed test. See https://playwright.dev/docs/trace-viewer */
        trace: 'on-first-retry',

        /* Screenshot options for consistent results */
        screenshot: 'only-on-failure',
    },

    /* Configure projects for major browsers */
    projects: [
        {
            name: 'chromium',
            use: {
                ...devices['Desktop Chrome'],
                // Force consistent viewport for screenshots
                viewport: {width: 1280, height: 720},
            },
        },
    ],

    /* Run your local dev server before starting the tests */
    webServer: {
        command: 'npm run ladle',
        url: 'http://localhost:61000',
        reuseExistingServer: !process.env.CI,
        timeout: 120 * 1000,
    },
});
