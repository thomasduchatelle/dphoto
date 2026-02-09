import type {Meta, StoryObj} from '@storybook/nextjs-vite';
import {fn} from 'storybook/test';
import {ErrorDisplay} from './index';
import {Box} from '@mui/material';

const meta = {
    title: 'Shared/ErrorDisplay',
    component: ErrorDisplay,
    parameters: {
        layout: 'centered',
    },
    tags: ['autodocs'],
} satisfies Meta<typeof ErrorDisplay>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
    args: {
        error: {
            message: 'Failed to load albums',
        },
        onRetry: fn(),
    },
};

export const WithTechnicalDetails: Story = {
    args: {
        error: {
            message: 'Failed to connect to the server',
            code: 'ERR_NETWORK',
            details: 'Network request failed: Connection timeout after 30000ms',
        },
        onRetry: fn(),
    },
};

export const NoRetry: Story = {
    args: {
        error: {
            message: 'You do not have permission to access this resource',
            code: 'ERR_FORBIDDEN',
        },
    },
};

export const LongMessage: Story = {
    args: {
        error: {
            message:
                'An unexpected error occurred while trying to process your request. The server encountered an internal error and was unable to complete your request. Please try again later or contact support if the problem persists.',
        },
        onRetry: fn(),
    },
};

export const Standalone: Story = {
    render: () => (
        <Box
            sx={{
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                minHeight: '100vh',
                p: 2,
            }}
        >
            <ErrorDisplay
                error={{
                    message: 'Page failed to load',
                }}
                onRetry={fn()}
            />
        </Box>
    ),
};

export const InlineContext: Story = {
    render: () => (
        <Box sx={{p: 3}}>
            <Box sx={{mb: 2}}>Some form content above...</Box>
            <ErrorDisplay
                error={{
                    message: 'Failed to save changes',
                }}
                onRetry={fn()}
                onDismiss={fn()}
            />
            <Box sx={{mt: 2}}>Some form content below...</Box>
        </Box>
    ),
};
