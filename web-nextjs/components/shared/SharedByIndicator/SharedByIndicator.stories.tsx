import type {Meta, StoryObj} from '@storybook/nextjs-vite';
import {SharedByIndicator} from './index';
import {Box, Typography} from '@mui/material';

const users = [
    {name: 'Alice Johnson', email: 'alice@example.com', picture: 'https://i.pravatar.cc/150?img=1'},
    {name: 'Bob Smith', email: 'bob@example.com', picture: 'https://i.pravatar.cc/150?img=2'},
    {name: 'Carol White', email: 'carol@example.com', picture: 'https://i.pravatar.cc/150?img=3'},
    {name: 'David Brown', email: 'david@example.com'},
    {name: 'Eve Davis', email: 'eve@example.com'},
    {name: 'Frank Miller', email: 'frank@example.com'},
    {name: 'Grace Wilson', email: 'grace@example.com'},
    {name: 'Henry Moore', email: 'henry@example.com'},
    {name: 'Ivy Taylor', email: 'ivy@example.com'},
    {name: 'Jack Anderson', email: 'jack@example.com'},
];

const meta = {
    title: 'Shared/SharedByIndicator',
    component: SharedByIndicator,
    parameters: {
        layout: 'centered',
    },
    tags: ['autodocs'],
} satisfies Meta<typeof SharedByIndicator>;

export default meta;
type Story = StoryObj<typeof meta>;

export const OneUser: Story = {
    render: () => (
        <Box sx={{p: 2}}>
            <Typography variant="body2" gutterBottom>
                Shared with 1 user:
            </Typography>
            <SharedByIndicator users={users.slice(0, 1)}/>
        </Box>
    ),
};

export const ThreeUsers: Story = {
    render: () => (
        <Box sx={{p: 2}}>
            <Typography variant="body2" gutterBottom>
                Shared with 3 users (all visible):
            </Typography>
            <SharedByIndicator users={users.slice(0, 3)}/>
        </Box>
    ),
};

export const FiveUsers: Story = {
    render: () => (
        <Box sx={{p: 2}}>
            <Typography variant="body2" gutterBottom>
                Shared with 5 users (3 visible + &quot;+2&quot;):
            </Typography>
            <SharedByIndicator users={users.slice(0, 5)}/>
        </Box>
    ),
};

export const TenUsers: Story = {
    render: () => (
        <Box sx={{p: 2}}>
            <Typography variant="body2" gutterBottom>
                Shared with 10 users (3 visible + &quot;+7&quot;):
            </Typography>
            <SharedByIndicator users={users}/>
        </Box>
    ),
};

export const HoverTooltip: Story = {
    render: () => (
        <Box sx={{p: 2}}>
            <Typography variant="body2" gutterBottom>
                Hover over &quot;+7&quot; to see tooltip with names:
            </Typography>
            <SharedByIndicator users={users}/>
            <Typography variant="caption" sx={{mt: 2, display: 'block', fontStyle: 'italic'}}>
                (Hover or focus on the +7 avatar to see remaining names)
            </Typography>
        </Box>
    ),
};
