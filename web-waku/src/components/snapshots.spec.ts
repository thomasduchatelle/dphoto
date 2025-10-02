import { test, expect } from "@playwright/test";
// we can't create tests asynchronously, thus using the sync-fetch lib
import fetch from "sync-fetch";

// URL where Ladle is served
const url = "http://localhost:61000";

// fetch Ladle's meta file
// https://ladle.dev/docs/meta
const stories = fetch(`${url}/meta.json`).json().stories;

// iterate through stories
Object.keys(stories).forEach((storyKey) => {
    // create a test for each story
    test(`${storyKey} - compare snapshots`, async ({ page }) => {
        // skip stories with `meta.skip` set to true
        test.skip(stories[storyKey].meta.skip, "meta.skip is true");

        // navigate to the story
        await page.goto(`${url}/?story=${storyKey}&mode=preview`);

        // stories are code-splitted, wait for them to be loaded
        await page.waitForSelector("[data-storyloaded]");

        // Wait for network to be idle to ensure all resources are loaded
        await page.waitForLoadState('networkidle');

        // Additional wait for any animations or transitions to complete
        const waitInMilliseconds = stories[storyKey].meta.skip ?? 0;
        if (waitInMilliseconds > 0) {
            await page.waitForTimeout(waitInMilliseconds);
        }

        // Hide any elements that might cause flakiness (like cursors, time-dependent content)
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

        // take a screenshot and compare it with the baseline
        await expect(page).toHaveScreenshot(`${storyKey}.png`, {
            // Use full page screenshots for consistency
            fullPage: true,
            // Add threshold for minor rendering differences
            threshold: 0.2,
        });
    });
});