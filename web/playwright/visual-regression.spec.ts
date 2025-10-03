import {expect, test} from "@playwright/test";
import fetch from "sync-fetch";

const url = "http://localhost:61000";

const stories = fetch(`${url}/meta.json`).json().stories;

// iterate through stories
Object.keys(stories).forEach((storyKey) => {
    [{width: 1080, height: 800}, {width: 414, height: 915}].forEach(widthOption => {

        if (stories[storyKey].meta.skipSnapshot) {
            return
        }

        const viewportId = `${widthOption.width}x${widthOption.height}`;
        test(`${storyKey} - compare snapshots [${viewportId}]`, async ({page}) => {
            await page.setViewportSize(widthOption)
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
            await expect(page).toHaveScreenshot(`${storyKey}-${viewportId}.png`, {
                // Use full page screenshots for consistency
                fullPage: true,
                // Add threshold for minor rendering differences
                threshold: 0.2,
            });
        });
    })
});