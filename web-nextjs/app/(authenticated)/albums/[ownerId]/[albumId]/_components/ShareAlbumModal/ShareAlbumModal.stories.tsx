import type {Meta, StoryObj} from '@storybook/nextjs-vite';
import {fn, userEvent, within} from 'storybook/test';
import {useState} from 'react';
import {Button} from '@mui/material';
import {ShareAlbumModal} from './index';
import {Sharing, UserDetails} from '@/domains/catalog/language';

const tonyStark: UserDetails = {
    name: 'Tony Stark',
    email: 'tony@stark-industries.com',
    picture: '/tonystark-profile.jpg',
};

const bruceBanner: UserDetails = {
    name: 'Bruce Banner',
    email: 'bruce@avengers.com',
};

const natasha: UserDetails = {
    name: 'Natasha Romanoff',
    email: 'natasha@avengers.com',
};

const sampleSharedWith: Sharing[] = [
    {user: tonyStark},
    {user: bruceBanner},
    {user: natasha},
];

const meta = {
    title: 'Catalog/ShareAlbumModal',
    component: ShareAlbumModal,
    parameters: {layout: 'fullscreen'},
    decorators: [
        (Story, context) => {
            const [open, setOpen] = useState(true);
            return (
                <>
                    <Button variant="contained" onClick={() => setOpen(true)}>Reopen dialog</Button>
                    <Story args={{...context.args, open, onClose: () => setOpen(false)}}/>
                </>
            );
        },
    ],
    args: {
        open: true,
        sharedWith: sampleSharedWith,
        suggestions: [],
        onClose: fn(),
        onGrant: fn(() => Promise.resolve()),
        onRevoke: fn(() => Promise.resolve()),
    },
} satisfies Meta<typeof ShareAlbumModal>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const EmptyAccessList: Story = {
    args: {
        sharedWith: [],
    },
};

export const EmailValidationError: Story = {
    play: async ({canvasElement}) => {
        const canvas = within(canvasElement.ownerDocument.body);
        const input = await canvas.findByLabelText('Email to grant access');
        await userEvent.click(input);
        await userEvent.type(input, 'not-an-email');
        await userEvent.tab();
    },
};

export const GrantError: Story = {
    args: {
        error: {
            type: 'grant',
            email: 'ghost@unknown.com',
            message: 'Failed to grant access, verify the email address or contact maintainers',
        },
    },
};

export const RevokeConfirmationOpen: Story = {
    play: async ({canvasElement}) => {
        const canvas = within(canvasElement.ownerDocument.body);
        const revokeButton = await canvas.findByLabelText('Revoke access for Bruce Banner');
        await userEvent.click(revokeButton);
    },
};

export const RevokeError: Story = {
    args: {
        error: {
            type: 'revoke',
            email: bruceBanner.email,
            message: `Couldn't revoke access of user ${bruceBanner.email}, try again later`,
        },
    },
};
