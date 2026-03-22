import {expect, test} from "@playwright/test";
import fetch from "sync-fetch";

const url = "http://localhost:6006";

const entries: Record<string, {id: string; title: string; name: string; type: string; tags?: string[]}> = fetch(`${url}/index.json`).json().entries;

Object.values(entries)
    .filter(entry => entry.type === "story")
    .forEach((story) => {
        [{width: 1080, height: 800}, {width: 414, height: 915}].forEach(viewport => {
            const viewportId = `${viewport.width}x${viewport.height}`;
            test(`${story.id} - compare snapshots [${viewportId}]`, async ({page}) => {
                await page.setViewportSize(viewport);
                await page.goto(`${url}/iframe.html?id=${story.id}&viewMode=story`);

                await page.waitForLoadState('networkidle');

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

                await expect(page).toHaveScreenshot(`${story.id}-${viewportId}.png`, {
                    fullPage: true,
                    threshold: 0.2,
                });
            });
        });
    });
