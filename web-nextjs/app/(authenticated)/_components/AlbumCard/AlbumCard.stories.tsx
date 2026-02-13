import type {Meta, StoryObj} from '@storybook/nextjs-vite';
import {fn} from 'storybook/test';
import {AlbumCard} from './index';
import {Box} from '@mui/material';

const sampleAlbum = {
    albumId: 'album-1',
    ownerId: 'owner-1',
    name: 'Summer Vacation 2026',
    startDate: '2026-01-15',
    endDate: '2026-02-28',
    mediaCount: 47,
};

const sharedUsers = [
    {name: 'Alice Johnson', email: 'alice@example.com', picture: 'https://i.pravatar.cc/150?img=1'},
    {name: 'Bob Smith', email: 'bob@example.com'},
    {name: 'Carol White', email: 'carol@example.com', picture: 'https://i.pravatar.cc/150?img=3'},
];

const meta = {
    title: 'Components/AlbumCard',
    component: AlbumCard,
    parameters: {
        layout: 'centered',
    },
    tags: ['autodocs'],
    args: {onClick: fn()},
} satisfies Meta<typeof AlbumCard>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
    render: () => (
        <Box sx={{width: 400}}>
            <AlbumCard album={sampleAlbum} onClick={fn()}/>
        </Box>
    ),
};

export const HighDensity: Story = {
    render: () => (
        <Box sx={{width: 400}}>
            <AlbumCard
                album={{
                    ...sampleAlbum,
                    name: 'Daily Photography Challenge',
                    mediaCount: 150,
                    startDate: '2026-01-01',
                    endDate: '2026-01-10',
                }}
                onClick={fn()}
            />
        </Box>
    ),
};

export const LowDensity: Story = {
    render: () => (
        <Box sx={{width: 400}}>
            <AlbumCard
                album={{
                    ...sampleAlbum,
                    name: 'Occasional Moments',
                    mediaCount: 15,
                    startDate: '2026-01-01',
                    endDate: '2026-01-30',
                }}
                onClick={fn()}
            />
        </Box>
    ),
};

export const SharedAlbum: Story = {
    render: () => (
        <Box sx={{width: 400}}>
            <AlbumCard
                album={sampleAlbum}
                owner={{
                    name: 'John Doe',
                    email: 'john.doe@example.com',
                    picture: 'https://i.pravatar.cc/150?img=12',
                }}
                onClick={fn()}
            />
        </Box>
    ),
};

export const AlbumIShared: Story = {
    render: () => (
        <Box sx={{width: 400}}>
            <AlbumCard album={sampleAlbum} sharedWith={sharedUsers} onClick={fn()}/>
        </Box>
    ),
};

export const LongName: Story = {
    render: () => (
        <Box sx={{width: 400}}>
            <AlbumCard
                album={{
                    ...sampleAlbum,
                    name: 'My Amazing Summer Vacation Trip to Europe Including France, Italy, and Spain',
                }}
                onClick={fn()}
            />
        </Box>
    ),
};

export const Mobile: Story = {
    render: () => (
        <Box sx={{width: 300}}>
            <AlbumCard album={sampleAlbum} onClick={fn()}/>
        </Box>
    ),
};

export const Desktop: Story = {
    render: () => (
        <Box sx={{width: 400}}>
            <AlbumCard album={sampleAlbum} onClick={fn()}/>
        </Box>
    ),
};
