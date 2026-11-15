import type {Meta, StoryObj} from '@storybook/nextjs-vite';
import {Box} from '@mui/material';
import {AppNavView} from './AppNav';
import {AppBackground} from './AppBackground';

const meta = {
    title: 'Components/AppNav',
    component: AppNavView,
    parameters: {layout: 'fullscreen'},
    decorators: [
        (Story) => (
            <AppBackground>
                <Box sx={{p: 3}}>
                    <Story/>
                </Box>
            </AppBackground>
        ),
    ],
    args: {
        activeSection: 'albums',
    },
} satisfies Meta<typeof AppNavView>;

export default meta;
type Story = StoryObj<typeof meta>;

export const AlbumsActive: Story = {};

export const TagsActive: Story = {
    args: {activeSection: 'tags'},
};
