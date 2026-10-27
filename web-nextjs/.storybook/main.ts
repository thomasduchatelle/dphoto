import type {StorybookConfig} from '@storybook/nextjs-vite';

const config: StorybookConfig = {
    "stories": [
        "../**/*.stories.@(js|jsx|mjs|ts|tsx)"
    ],
    "addons": [
        "@chromatic-com/storybook",
        "@storybook/addon-vitest",
        "@storybook/addon-a11y",
        "@storybook/addon-docs",
        "@storybook/addon-onboarding"
    ],
    "framework": "@storybook/nextjs-vite",
    "staticDirs": [
        "../public",
        "../stories/assets",
        "../../test/wiremock/__files/api",
    ]
};
export default config;