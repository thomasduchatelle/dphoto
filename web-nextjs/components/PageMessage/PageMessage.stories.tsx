import type {Meta, StoryObj} from '@storybook/nextjs-vite';
import {fn} from 'storybook/test';
import {PageMessage} from './index';
import {Button} from '@mui/material';
import {AppBackground} from '@/components/AppLayout/AppBackground';
import CollectionsIcon from '@mui/icons-material/Collections';
import CheckCircleOutlineIcon from '@mui/icons-material/CheckCircleOutline';

const meta = {
    title: 'Components/PageMessage',
    component: PageMessage,
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
        icon: <CollectionsIcon/>,
        title: 'No Albums Found',
        message: 'Create your first album to get started organizing your photos.',
    },
} satisfies Meta<typeof PageMessage>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const InfoWithActions: Story = {
    args: {
        icon: <CollectionsIcon/>,
        title: 'No Albums Found',
        message: 'Create your first album to get started organizing your photos.',
    },
    render: (args) => (
        <PageMessage {...args}>
            <Button variant="contained" onClick={fn()}>
                Create Album
            </Button>
            <Button variant="outlined" onClick={fn()}>
                Learn More
            </Button>
        </PageMessage>
    ),
};

export const Success: Story = {
    args: {
        variant: 'success',
        icon: <CheckCircleOutlineIcon/>,
        title: 'Successfully Logged Out',
        message: 'You have been successfully logged out of your account.',
    },
};

export const SuccessWithAction: Story = {
    args: {
        variant: 'success',
        icon: <CheckCircleOutlineIcon/>,
        title: 'Successfully Logged Out',
        message: 'You have been successfully logged out of your account.',
    },
    render: (args) => (
        <PageMessage {...args}>
            <Button variant="contained" onClick={fn()}>
                Sign In Again
            </Button>
        </PageMessage>
    ),
};
