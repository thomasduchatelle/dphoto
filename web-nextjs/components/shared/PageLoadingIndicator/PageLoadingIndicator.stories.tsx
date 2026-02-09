import type {Meta, StoryObj} from '@storybook/nextjs-vite';
import {PageLoadingIndicator} from './index';
import {Box} from '@mui/material';

const meta = {
    title: 'Shared/PageLoadingIndicator',
    component: PageLoadingIndicator,
    parameters: {
        layout: 'centered',
    },
    tags: ['autodocs'],
} satisfies Meta<typeof PageLoadingIndicator>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
    render: () => (
        <Box sx={{position: 'relative', height: 200}}>
            <PageLoadingIndicator/>
        </Box>
    ),
};

export const CustomMessage: Story = {
    render: () => (
        <Box sx={{position: 'relative', height: 200}}>
            <PageLoadingIndicator message="Loading albums..."/>
        </Box>
    ),
};

export const NoMessage: Story = {
    render: () => (
        <Box sx={{position: 'relative', height: 200}}>
            <PageLoadingIndicator message=""/>
        </Box>
    ),
};
