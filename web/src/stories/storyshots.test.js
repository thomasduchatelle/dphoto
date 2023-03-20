import initStoryshots from '@storybook/addon-storyshots';
import {imageSnapshot} from '@storybook/addon-storyshots-puppeteer';
import path from "path";

const getMatchOptions = () => {
    return {
        comparisonMethod: 'ssim',
        customDiffConfig: {
            ssim: 'fast',
        },
        failureThreshold: 0.015,
        failureThresholdType: 'percent',
    };
};

const beforeScreenshot = (page, {context: {parameters}}) => {
    if (parameters["delay"] > 0) {
        return new Promise((resolve) =>
            setTimeout(() => {
                resolve();
            }, parameters["delay"])
        );
    }
};

const storybookUrl = process.env.CI === "true" ? `file://${path.resolve(__dirname, '../../storybook-static')}` : 'http://localhost:6006/'

initStoryshots({
    test: imageSnapshot({
        storybookUrl,
        getMatchOptions,
        beforeScreenshot,
    }),
});


// ----------------------------------------------------------------------------
// Below configuration kept for reference but do not worth the extra complexity
// Works better with options:
// <Box sx={{
//     textRendering: "geometricprecision !important",
// }}>
//     {children}
// </Box>
const getCustomBrowser = () => {
    const puppeteer = require('puppeteer');
    return puppeteer.launch({
        // , '--force-color-profile=srgb', '--enable-font-antialiasing', '--disable-gpu'
        args: [
            '--no-sandbox ',
            '--disable-setuid-sandbox',
            '--disable-dev-shm-usage',
            "--window-size=1440,1080",
            "--font-render-hinting=none",
            "--force-color-profile=generic-rgb",
            "--disable-gpu",
            "--disable-translate",
            "--disable-extensions",
            "--disable-accelerated-2d-canvas",
            // "--deterministic-mode",
            "--disable-skia-runtime-opts",
            "--force-device-scale-factor=1",
            "--js-flags=--random-seed=1157259157",
            "--disable-partial-raster",
            "--use-gl=swiftshader"
        ],
        // args: ['--no-sandbox ', '--disable-setuid-sandbox', '--disable-dev-shm-usage', '--font-render-hinting=medium'],
        // executablePath: chromeExecutablePath,
    })
};

const customizePage = (page) => {
    return page.setUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.121 Safari/537.36");
}