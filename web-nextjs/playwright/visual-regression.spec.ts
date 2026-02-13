import {expect, test} from "@playwright/test";
import fetch from "sync-fetch";

const url = "http://localhost:6006";

// Storybook uses index.json (not meta.json like Ladle)
// Fetch the story index from Storybook
const indexJson = fetch(`${url}/index.json`).json();
const stories = indexJson.entries;

// Test viewports: desktop and mobile
const viewports = [
    {width: 1080, height: 800},
    {width: 414, height: 915}
];

// Iterate through all stories
Object.keys(stories).forEach((storyId) => {
    const story = stories[storyId];

    viewports.forEach(widthOption => {
        // Skip stories with 'no-tests' or 'skip-test' tags
        // To skip a story from visual testing, add tags: ['no-tests'] to the story or meta
        if (story.tags && (story.tags.includes('no-tests') || story.tags.includes('skip-test'))) {
            return;
        }

        const viewportId = `${widthOption.width}x${widthOption.height}`;
        test(`${storyId} - compare snapshots [${viewportId}]`, async ({page}) => {
            await page.setViewportSize(widthOption);

            // Navigate to the story in Storybook's iframe
            // Storybook URL format: iframe.html?id=story-id&viewMode=story
            await page.goto(`${url}/iframe.html?id=${storyId}&viewMode=story`);

            // Wait for Storybook to render the story
            await page.waitForSelector('#storybook-root', {timeout: 10000});

            // Wait for network to be idle to ensure all resources are loaded
            await page.waitForLoadState('networkidle');

            // Additional small delay to ensure rendering is complete
            await page.evaluate(() => {
                return new Promise(resolve => setTimeout(resolve, 100));
            });

            // Disable animations and transitions to prevent flakiness
            await page.addStyleTag({
                content: `
                *, *::before, *::after {
                    animation-duration: 0s !important;
                    animation-delay: 0s !important;
                    transition-duration: 0s !important;
                    transition-delay: 0s !important;
                }
                *:focus {
                    outline: none !important;
                }
            `
            });

            // Take a screenshot and compare it with the baseline
            await expect(page).toHaveScreenshot(`${storyId}-${viewportId}.png`, {
                // Use full page screenshots for consistency
                fullPage: true,
                // Add threshold for minor rendering differences (fonts, anti-aliasing, etc.)
                threshold: 0.2,
            });
        });
    })
});
