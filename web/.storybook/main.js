module.exports = {
    "stories": ["../src/**/*.mdx", "../src/**/*.stories.@(js|jsx|ts|tsx)"],

    "addons": [
        "@storybook/addon-links",
        "@storybook/addon-essentials",
        "@storybook/addon-interactions",
        "@storybook/preset-create-react-app"
    ],

    "framework": {
        name: "@storybook/react-webpack5",
        options: {}
    },

    "staticDirs": ['../public', '../src/images', "../../test/wiremock/__files/api/static"],

    docs: {},

    typescript: {
        reactDocgen: "react-docgen-typescript"
    }
}