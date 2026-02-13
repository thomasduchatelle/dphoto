import type {Meta, StoryObj} from '@storybook/nextjs-vite';
import {fn} from 'storybook/test';
import {AlbumGrid} from './index';
import {AlbumCard} from '../AlbumCard';
import {Box, Typography} from '@mui/material';
import {AppBackground} from '@/components/AppLayout/AppBackground';
import {Album, AlbumId} from '@/domains/catalog/language/catalog-state';

const createAlbumId = (owner: string, folderName: string): AlbumId => ({owner, folderName});

// Sample albums with diverse properties matching the actual Album domain type
const sampleAlbums: Album[] = [
    {
        albumId: createAlbumId('owner-1', 'santa-cruz-beach'),
        name: 'Santa Cruz Beach Trip',
        start: new Date('2026-07-15'),
        end: new Date('2026-07-22'),
        totalCount: 47,
        temperature: 8.2,
        relativeTemperature: 0.45,
        sharedWith: [{user: {name: 'Alice Johnson', email: 'alice@example.com', picture: 'https://i.pravatar.cc/150?img=1'}}],
    },
    {
        albumId: createAlbumId('owner-1', 'iceland-adventure'),
        name: 'Iceland Adventure',
        start: new Date('2025-08-01'),
        end: new Date('2025-08-14'),
        totalCount: 203,
        temperature: 15.6,
        relativeTemperature: 0.92,
        sharedWith: [],
    },
    {
        albumId: createAlbumId('owner-1', 'summer-family'),
        name: 'Summer Family Gathering',
        start: new Date('2025-06-01'),
        end: new Date('2025-08-31'),
        totalCount: 124,
        temperature: 11.3,
        relativeTemperature: 0.68,
        sharedWith: [
            {user: {name: 'Bob Smith', email: 'bob@example.com', picture: 'https://i.pravatar.cc/150?img=2'}},
            {user: {name: 'Carol White', email: 'carol@example.com', picture: 'https://i.pravatar.cc/150?img=3'}},
        ],
    },
    {
        albumId: createAlbumId('owner-1', 'paris-weekend'),
        name: 'Paris Weekend',
        start: new Date('2025-03-12'),
        end: new Date('2025-03-15'),
        totalCount: 89,
        temperature: 9.1,
        relativeTemperature: 0.52,
        sharedWith: [],
    },
    {
        albumId: createAlbumId('owner-1', 'christmas-2024'),
        name: 'Christmas 2024',
        start: new Date('2024-12-23'),
        end: new Date('2024-12-26'),
        totalCount: 156,
        temperature: 14.8,
        relativeTemperature: 0.88,
        sharedWith: [
            {user: {name: 'David Lee', email: 'david@example.com', picture: 'https://i.pravatar.cc/150?img=4'}},
            {user: {name: 'Emma Brown', email: 'emma@example.com', picture: 'https://i.pravatar.cc/150?img=5'}},
            {user: {name: 'Frank Wilson', email: 'frank@example.com', picture: 'https://i.pravatar.cc/150?img=6'}},
        ],
    },
    {
        albumId: createAlbumId('owner-1', 'alps-hiking'),
        name: 'Hiking in the Alps',
        start: new Date('2024-09-05'),
        end: new Date('2024-09-12'),
        totalCount: 234,
        temperature: 16.2,
        relativeTemperature: 0.96,
        sharedWith: [{user: {name: 'Grace Taylor', email: 'grace@example.com', picture: 'https://i.pravatar.cc/150?img=7'}}],
    },
    {
        albumId: createAlbumId('owner-1', 'new-york-city'),
        name: 'New York City',
        start: new Date('2024-11-18'),
        end: new Date('2024-11-23'),
        totalCount: 178,
        temperature: 13.4,
        relativeTemperature: 0.81,
        sharedWith: [],
    },
    {
        albumId: createAlbumId('owner-1', 'tokyo-spring'),
        name: 'Tokyo Spring',
        start: new Date('2024-04-01'),
        end: new Date('2024-04-10'),
        totalCount: 312,
        temperature: 17.5,
        relativeTemperature: 1.0,
        sharedWith: [{user: {name: 'Henry Davis', email: 'henry@example.com', picture: 'https://i.pravatar.cc/150?img=8'}}],
    },
    {
        albumId: createAlbumId('owner-1', 'family-reunion'),
        name: 'Family Reunion',
        start: new Date('2024-07-01'),
        end: new Date('2024-07-04'),
        totalCount: 93,
        temperature: 9.7,
        relativeTemperature: 0.56,
        sharedWith: [
            {user: {name: 'Iris Chen', email: 'iris@example.com', picture: 'https://i.pravatar.cc/150?img=9'}},
            {user: {name: 'Jack Martinez', email: 'jack@example.com', picture: 'https://i.pravatar.cc/150?img=10'}},
        ],
    },
    {
        albumId: createAlbumId('owner-2', 'bobs-vacation'),
        name: 'Bob\'s Vacation Shared',
        start: new Date('2024-05-10'),
        end: new Date('2024-05-20'),
        totalCount: 67,
        temperature: 7.8,
        relativeTemperature: 0.42,
        ownedBy: {name: 'Bob Smith', email: 'bob@example.com', picture: 'https://i.pravatar.cc/150?img=2'},
        sharedWith: [],
    },
    {
        albumId: createAlbumId('owner-1', 'winter-ski-trip'),
        name: 'Winter Ski Trip',
        start: new Date('2024-01-15'),
        end: new Date('2024-01-22'),
        totalCount: 145,
        temperature: 12.1,
        relativeTemperature: 0.73,
        sharedWith: [],
    },
    {
        albumId: createAlbumId('owner-1', 'autumn-foliage'),
        name: 'Autumn Foliage Tour',
        start: new Date('2023-10-01'),
        end: new Date('2023-10-15'),
        totalCount: 189,
        temperature: 14.3,
        relativeTemperature: 0.85,
        sharedWith: [{user: {name: 'Karen Wilson', email: 'karen@example.com', picture: 'https://i.pravatar.cc/150?img=11'}}],
    },
];

const meta = {
    title: 'Components/AlbumGrid',
    component: AlbumGrid,
    parameters: {
        layout: 'fullscreen',
    },
    tags: ['autodocs'],
    globals: {viewport: {}},
} satisfies Meta<typeof AlbumGrid>;

export default meta;
type Story = StoryObj<typeof meta>;

// ============================================================================
// 1. Default State - Full Grid
// ============================================================================

export const Default: Story = {
    render: () => (
        <AppBackground>
            <Box sx={{p: 6}}>
                <Typography
                    variant="h1"
                    sx={{
                        mb: 4,
                        fontSize: 32,
                        fontWeight: 300,
                        color: '#4a9ece',
                        textTransform: 'uppercase',
                        letterSpacing: '0.08em',
                    }}
                >
                    Your Albums
                </Typography>
                <AlbumGrid>
                    {sampleAlbums.map(album => (
                        <AlbumCard key={`${album.albumId.owner}-${album.albumId.folderName}`} album={album} onClick={fn()} />
                    ))}
                </AlbumGrid>
            </Box>
        </AppBackground>
    ),
};

// ============================================================================
// 2. Responsive Breakpoints
// ============================================================================

export const ResponsiveBreakpoints: Story = {
    render: () => (
        <AppBackground>
            <Box sx={{p: 6}}>
                {/* Mobile - 1 column */}
                <Box sx={{mb: 8}}>
                    <Typography
                        sx={{
                            mb: 2,
                            fontSize: 14,
                            fontWeight: 500,
                            color: '#4a9ece',
                            textTransform: 'uppercase',
                            letterSpacing: '0.1em',
                        }}
                    >
                        Mobile (xs) - 1 Column
                    </Typography>
                    <Box sx={{maxWidth: 400, border: '1px dashed rgba(74, 158, 206, 0.3)', p: 2}}>
                        <AlbumGrid>
                            {sampleAlbums.slice(0, 3).map(album => (
                                <AlbumCard key={`${album.albumId.owner}-${album.albumId.folderName}`} album={album} onClick={fn()} />
                            ))}
                        </AlbumGrid>
                    </Box>
                </Box>

                {/* Tablet - 2 columns */}
                <Box sx={{mb: 8}}>
                    <Typography
                        sx={{
                            mb: 2,
                            fontSize: 14,
                            fontWeight: 500,
                            color: '#4a9ece',
                            textTransform: 'uppercase',
                            letterSpacing: '0.1em',
                        }}
                    >
                        Tablet (sm) - 2 Columns
                    </Typography>
                    <Box sx={{maxWidth: 700, border: '1px dashed rgba(74, 158, 206, 0.3)', p: 2}}>
                        <AlbumGrid>
                            {sampleAlbums.slice(0, 4).map(album => (
                                <AlbumCard key={`${album.albumId.owner}-${album.albumId.folderName}`} album={album} onClick={fn()} />
                            ))}
                        </AlbumGrid>
                    </Box>
                </Box>

                {/* Desktop - 3 columns */}
                <Box sx={{mb: 8}}>
                    <Typography
                        sx={{
                            mb: 2,
                            fontSize: 14,
                            fontWeight: 500,
                            color: '#4a9ece',
                            textTransform: 'uppercase',
                            letterSpacing: '0.1em',
                        }}
                    >
                        Desktop (md) - 3 Columns
                    </Typography>
                    <Box sx={{maxWidth: 1100, border: '1px dashed rgba(74, 158, 206, 0.3)', p: 2}}>
                        <AlbumGrid>
                            {sampleAlbums.slice(0, 6).map(album => (
                                <AlbumCard key={`${album.albumId.owner}-${album.albumId.folderName}`} album={album} onClick={fn()} />
                            ))}
                        </AlbumGrid>
                    </Box>
                </Box>

                {/* Large Desktop - 4 columns */}
                <Box>
                    <Typography
                        sx={{
                            mb: 2,
                            fontSize: 14,
                            fontWeight: 500,
                            color: '#4a9ece',
                            textTransform: 'uppercase',
                            letterSpacing: '0.1em',
                        }}
                    >
                        Large Desktop (lg) - 4 Columns
                    </Typography>
                    <Box sx={{maxWidth: 1600, border: '1px dashed rgba(74, 158, 206, 0.3)', p: 2}}>
                        <AlbumGrid>
                            {sampleAlbums.slice(0, 8).map(album => (
                                <AlbumCard key={`${album.albumId.owner}-${album.albumId.folderName}`} album={album} onClick={fn()} />
                            ))}
                        </AlbumGrid>
                    </Box>
                </Box>
            </Box>
        </AppBackground>
    ),
};

// ============================================================================
// 3. Mobile Viewport Story
// ============================================================================

export const DefaultMobile: Story = {
    render: () => (
        <AppBackground>
            <Box sx={{p: 3}}>
                <Typography
                    variant="h1"
                    sx={{
                        mb: 3,
                        fontSize: 24,
                        fontWeight: 300,
                        color: '#4a9ece',
                        textTransform: 'uppercase',
                        letterSpacing: '0.08em',
                    }}
                >
                    Your Albums
                </Typography>
                <AlbumGrid>
                    {sampleAlbums.slice(0, 6).map(album => (
                        <AlbumCard key={`${album.albumId.owner}-${album.albumId.folderName}`} album={album} onClick={fn()} />
                    ))}
                </AlbumGrid>
            </Box>
        </AppBackground>
    ),
    globals: {viewport: {value: 'mobile2', isRotated: false}},
};
