import initStoryshots from '@storybook/addon-storyshots';
import {imageSnapshot} from '@storybook/addon-storyshots-puppeteer';

const getMatchOptions = ({context: {kind, story}, url}) => {
    return {
        comparisonMethod: 'ssim',
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

initStoryshots({
    test: imageSnapshot({
        storybookUrl: 'http://localhost:6006/',
        getMatchOptions,
        beforeScreenshot,
    }),
});