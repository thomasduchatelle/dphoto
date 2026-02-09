import type {Meta, StoryObj} from '@storybook/nextjs-vite';
import {NavigationLoadingIndicator} from './index';
import {Box, Typography} from '@mui/material';

const meta = {
    title: 'Shared/NavigationLoadingIndicator',
    component: NavigationLoadingIndicator,
    parameters: {
        layout: 'centered',
    },
    tags: ['autodocs'],
} satisfies Meta<typeof NavigationLoadingIndicator>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Inline: Story = {
    render: () => (
        <Box sx={{p: 2, display: 'flex', alignItems: 'center', gap: 2}}>
            <NavigationLoadingIndicator variant="inline"/>
            <Typography variant="body2">Loading content...</Typography>
        </Box>
    ),
};

export const Overlay: Story = {
    render: () => (
        <Box sx={{position: 'relative', height: 400, bgcolor: 'background.default'}}>
            <Typography variant="h6" sx={{p: 2}}>
                Page content underneath
            </Typography>
            <NavigationLoadingIndicator variant="overlay"/>
        </Box>
    ),
};
