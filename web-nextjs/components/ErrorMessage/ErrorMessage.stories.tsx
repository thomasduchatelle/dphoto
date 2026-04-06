import type {Meta, StoryObj} from '@storybook/nextjs-vite';
import {ErrorMessage} from './index';
import {AppBackground} from '@/components/AppLayout/AppBackground';

const meta = {
    title: 'Components/ErrorMessage',
    component: ErrorMessage,
    parameters: {
        layout: 'fullscreen',
    },
    decorators: [
        (Story) => (
            <AppBackground>
                <Story/>
            </AppBackground>
        ),
    ],
    args: {
        error: new Error('An unexpected error occurred while loading your content.'),
    },
} satisfies Meta<typeof ErrorMessage>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const WithCustomTitle: Story = {
    args: {
        error: new Error('An unexpected error occurred while loading your albums.'),
        title: 'Failed to Load Albums',
    },
};

const errorWithoutStack = new Error('Network request failed');
errorWithoutStack.stack = undefined;

export const WithoutStackTrace: Story = {
    args: {
        error: errorWithoutStack,
        title: 'Failed to Load Albums',
    },
};
