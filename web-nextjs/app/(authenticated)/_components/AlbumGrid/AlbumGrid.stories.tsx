import type {Meta, StoryObj} from '@storybook/nextjs-vite';
import {fn} from 'storybook/test';
import {AlbumGrid} from './index';
import {AlbumCard} from '../AlbumCard';
import {Box} from '@mui/material';

const sampleAlbums = Array.from({length: 12}, (_, i) => ({
    albumId: `album-${i + 1}`,
    ownerId: 'owner-1',
    name: `Album ${i + 1}`,
    startDate: '2026-01-15',
    endDate: '2026-02-28',
    mediaCount: 47,
}));

const meta = {
    title: 'Components/AlbumGrid',
    component: AlbumGrid,
    parameters: {
        layout: 'padded',
    },
    tags: ['autodocs'],
} satisfies Meta<typeof AlbumGrid>;

export default meta;
type Story = StoryObj<typeof meta>;

export const OneCard: Story = {
    render: () => (
        <AlbumGrid>
            <AlbumCard album={sampleAlbums[0]} onClick={fn()}/>
        </AlbumGrid>
    ),
};

export const ThreeCards: Story = {
    render: () => (
        <AlbumGrid>
            {sampleAlbums.slice(0, 3).map(album => (
                <AlbumCard key={album.albumId} album={album} onClick={fn()}/>
            ))}
        </AlbumGrid>
    ),
};

export const TwelveCards: Story = {
    render: () => (
        <AlbumGrid>
            {sampleAlbums.map(album => (
                <AlbumCard key={album.albumId} album={album} onClick={fn()}/>
            ))}
        </AlbumGrid>
    ),
};

export const Mobile: Story = {
    render: () => (
        <Box sx={{width: 400}}>
            <AlbumGrid>
                {sampleAlbums.slice(0, 3).map(album => (
                    <AlbumCard key={album.albumId} album={album} onClick={fn()}/>
                ))}
            </AlbumGrid>
        </Box>
    ),
};

export const Tablet: Story = {
    render: () => (
        <Box sx={{width: 700}}>
            <AlbumGrid>
                {sampleAlbums.slice(0, 4).map(album => (
                    <AlbumCard key={album.albumId} album={album} onClick={fn()}/>
                ))}
            </AlbumGrid>
        </Box>
    ),
};

export const Desktop: Story = {
    render: () => (
        <Box sx={{width: 1200}}>
            <AlbumGrid>
                {sampleAlbums.slice(0, 6).map(album => (
                    <AlbumCard key={album.albumId} album={album} onClick={fn()}/>
                ))}
            </AlbumGrid>
        </Box>
    ),
};

export const LargeDesktop: Story = {
    render: () => (
        <Box sx={{width: 1600}}>
            <AlbumGrid>
                {sampleAlbums.map(album => (
                    <AlbumCard key={album.albumId} album={album} onClick={fn()}/>
                ))}
            </AlbumGrid>
        </Box>
    ),
};
