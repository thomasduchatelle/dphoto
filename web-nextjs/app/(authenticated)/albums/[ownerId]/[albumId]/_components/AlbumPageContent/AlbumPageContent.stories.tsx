import type {Meta, StoryObj} from '@storybook/nextjs-vite';
import {Box} from '@mui/material';
import {AlbumPageContent} from './index';
import {AppBackground} from '@/components/AppLayout/AppBackground';
import {AppHeader} from '@/components/AppLayout/AppHeader';
import {Album, AlbumId, MediaType} from '@/domains/catalog/language';
import {
    loadedStateWithTwoAlbums,
    someMediasByDays,
} from '@/domains/catalog/tests/test-helper-state';

// ---------------------------------------------------------------------------
// Sample albums — same thumbnails as AlbumGrid.stories.tsx
// ---------------------------------------------------------------------------

const createAlbumId = (owner: string, folderName: string): AlbumId => ({owner, folderName});

const sampleAlbums: Album[] = [
    {
        albumId: createAlbumId('sandfall', 'clair-obscur'),
        name: 'Clair Obscur',
        start: new Date('2025-04-24'),
        end: new Date('2025-06-01'),
        totalCount: 47,
        temperature: 6.7,
        relativeTemperature: 1,
        sharedWith: [],
        thumbnails: [
            '/thumbnails/clair-obscur-1.jpg',
            '/thumbnails/clair-obscur-2.jpg',
            '/thumbnails/clair-obscur-3.jpg',
            '/thumbnails/clair-obscur-4.jpg',
        ],
    },
    {
        albumId: createAlbumId('sony', 'astro-bot'),
        name: 'Astro Bot',
        start: new Date('2024-09-06'),
        end: new Date('2024-09-12'),
        totalCount: 23,
        temperature: 12.1,
        relativeTemperature: 0.72,
        sharedWith: [{user: {name: 'Tony Stark', email: 'ironman@avenger.com', picture: '/static/tonystark-profile.jpg'}}],
        thumbnails: [
            '/thumbnails/astro-bot-01.jpg',
            '/thumbnails/astro-bot-02.jpg',
            '/thumbnails/astro-bot-03.jpg',
        ],
    },
    {
        albumId: createAlbumId('cdprojekt', 'the-witcher-3'),
        name: 'The Witcher 3: Wild Hunt',
        start: new Date('2024-05-30'),
        end: new Date('2024-06-14'),
        totalCount: 189,
        temperature: 9.4,
        relativeTemperature: 0.55,
        sharedWith: [],
        thumbnails: [
            '/thumbnails/the-witcher-3-01.jpg',
            '/thumbnails/the-witcher-3-02.jpg',
            '/thumbnails/the-witcher-3-03.jpg',
        ],
    },
    {
        albumId: createAlbumId('kojima', 'death-stranding-1'),
        name: 'Death Stranding',
        start: new Date('2023-11-08'),
        end: new Date('2023-11-22'),
        totalCount: 312,
        temperature: 17.5,
        relativeTemperature: 1.0,
        ownedBy: {name: 'Kojima', users: [{name: 'Tony Stark', email: 'ironman@avenger.com', picture: '/static/tonystark-profile.jpg'}]},
        sharedWith: [],
        thumbnails: [
            '/thumbnails/death-stranding-1-01.jpg',
            '/thumbnails/death-stranding-1-02.jpg',
            '/thumbnails/death-stranding-1-03.jpg',
            '/thumbnails/death-stranding-1-04.jpg',
        ],
    },
    {
        albumId: createAlbumId('kojima', 'death-stranding-2'),
        name: 'Death Stranding 2: On The Beach',
        start: new Date('2025-06-05'),
        end: new Date('2025-06-30'),
        totalCount: 78,
        temperature: 4.1,
        relativeTemperature: 0.22,
        sharedWith: [{user: {name: 'Tony Stark', email: 'ironman@avenger.com', picture: '/static/tonystark-profile.jpg'}}],
        thumbnails: ['/thumbnails/death-stranding-2-01.jpg'],
    },
    {
        albumId: createAlbumId('sandfall', 'clair-obscur-dlc'),
        name: 'Clair Obscur DLC',
        start: new Date('2025-09-10'),
        end: new Date('2025-09-15'),
        totalCount: 11,
        temperature: 1.2,
        relativeTemperature: 0.06,
        sharedWith: [],
        thumbnails: [],
    },
];

// ---------------------------------------------------------------------------
// Multi-day medias — thumbnails from the sample albums
// ---------------------------------------------------------------------------

const multiDayMedias = [
    {
        day: new Date(2025, 3, 25),
        medias: [
            {id: 'm1', type: MediaType.IMAGE, time: new Date('2025-04-25T09:00:00Z'), uiRelativePath: 'm1.jpg', contentPath: '/nextjs/thumbnails/clair-obscur-1.jpg', source: ''},
            {id: 'm2', type: MediaType.IMAGE, time: new Date('2025-04-25T10:00:00Z'), uiRelativePath: 'm2.jpg', contentPath: '/nextjs/thumbnails/clair-obscur-2.jpg', source: ''},
            {id: 'm3', type: MediaType.VIDEO, time: new Date('2025-04-25T11:00:00Z'), uiRelativePath: 'm3.mp4', contentPath: '/nextjs/thumbnails/clair-obscur-3.jpg', source: ''},
            {id: 'm4', type: MediaType.IMAGE, time: new Date('2025-04-25T12:00:00Z'), uiRelativePath: 'm4.jpg', contentPath: '/nextjs/thumbnails/clair-obscur-4.jpg', source: ''},
            {id: 'm5', type: MediaType.IMAGE, time: new Date('2025-04-25T13:00:00Z'), uiRelativePath: 'm5.jpg', contentPath: '/nextjs/thumbnails/astro-bot-01.jpg', source: ''},
            {id: 'm6', type: MediaType.IMAGE, time: new Date('2025-04-25T14:00:00Z'), uiRelativePath: 'm6.jpg', contentPath: '/nextjs/thumbnails/astro-bot-02.jpg', source: ''},
        ],
    },
    {
        day: new Date(2025, 4, 10),
        medias: [
            {id: 'm7', type: MediaType.IMAGE, time: new Date('2025-05-10T09:00:00Z'), uiRelativePath: 'm7.jpg', contentPath: '/nextjs/thumbnails/the-witcher-3-01.jpg', source: ''},
            {id: 'm8', type: MediaType.IMAGE, time: new Date('2025-05-10T10:00:00Z'), uiRelativePath: 'm8.jpg', contentPath: '/nextjs/thumbnails/the-witcher-3-02.jpg', source: ''},
            {id: 'm9', type: MediaType.VIDEO, time: new Date('2025-05-10T11:00:00Z'), uiRelativePath: 'm9.mp4', contentPath: '/nextjs/thumbnails/the-witcher-3-03.jpg', source: ''},
            {id: 'm10', type: MediaType.IMAGE, time: new Date('2025-05-10T12:00:00Z'), uiRelativePath: 'm10.jpg', contentPath: '/nextjs/thumbnails/death-stranding-1-01.jpg', source: ''},
            {id: 'm11', type: MediaType.IMAGE, time: new Date('2025-05-10T13:00:00Z'), uiRelativePath: 'm11.jpg', contentPath: '/nextjs/thumbnails/death-stranding-1-02.jpg', source: ''},
        ],
    },
    {
        day: new Date(2025, 4, 28),
        medias: [
            {id: 'm12', type: MediaType.IMAGE, time: new Date('2025-05-28T10:00:00Z'), uiRelativePath: 'm12.jpg', contentPath: '/nextjs/thumbnails/death-stranding-2-01.jpg', source: ''},
            {id: 'm13', type: MediaType.IMAGE, time: new Date('2025-05-28T11:00:00Z'), uiRelativePath: 'm13.jpg', contentPath: '/nextjs/thumbnails/death-stranding-1-03.jpg', source: ''},
            {id: 'm14', type: MediaType.IMAGE, time: new Date('2025-05-28T14:00:00Z'), uiRelativePath: 'm14.jpg', contentPath: '/nextjs/thumbnails/death-stranding-1-04.jpg', source: ''},
        ],
    },
];

// sampleAlbums is reverse-chron: index 0 = newest (Clair Obscur 2025-04)
// Displaying sampleAlbums[1] (Astro Bot): next=sampleAlbums[0] (top banner), previous=sampleAlbums[2] (bottom)
const baseState = {
    ...loadedStateWithTwoAlbums,
    allAlbums: sampleAlbums,
    albums: sampleAlbums,
    mediasLoadedFromAlbumId: sampleAlbums[1].albumId,
    medias: multiDayMedias,
};

const fakeUser = {
    name: 'Tony Stark',
    email: 'tony@stark-industries.com',
    picture: '/static/tonystark-profile.jpg',
    isOwner: true,
};

const meta = {
    title: 'Catalog/AlbumPageContent',
    component: AlbumPageContent,
    parameters: {layout: 'fullscreen'},
    decorators: [
        (Story: () => React.ReactNode) => (
            <AppBackground>
                <Box sx={{position: 'fixed', top: 0, left: 0, right: 0, zIndex: 1100}}>
                    <AppHeader user={fakeUser} logoutUrl="/auth/logout" isScrolled={false} basePath=""/>
                </Box>
                <Story/>
            </AppBackground>
        ),
    ],
    args: {
        initialState: baseState,
    },
} satisfies Meta<typeof AlbumPageContent>;

export default meta;
type Story = StoryObj<typeof meta>;

// Middle album: next=Clair Obscur (banner on scroll-up), previous=Witcher (bottom card)
export const WithMedias: Story = {};

// Newest album: no next banner, previous=Astro Bot card at bottom
export const NextAlbum: Story = {
    args: {
        initialState: {
            ...baseState,
            mediasLoadedFromAlbumId: sampleAlbums[0].albumId,
        },
    },
};

// Oldest displayed: next=DS2 banner on scroll-up, no previous card; few medias so bottom is visible
export const PreviousAlbum: Story = {
    args: {
        initialState: {
            ...baseState,
            mediasLoadedFromAlbumId: sampleAlbums[2].albumId,
            medias: someMediasByDays.map(d => ({
                ...d,
                medias: d.medias.map(m => ({...m, contentPath: '/nextjs/thumbnails/the-witcher-3-01.jpg'})),
            })),
        },
    },
};

// Only album: no next banner, no previous card
export const SingleAlbum: Story = {
    args: {
        initialState: {
            ...loadedStateWithTwoAlbums,
            allAlbums: [sampleAlbums[0]],
            albums: [sampleAlbums[0]],
            mediasLoadedFromAlbumId: sampleAlbums[0].albumId,
            medias: multiDayMedias,
        },
    },
};

// Empty album — both neighbours
export const NoMedias: Story = {
    args: {
        initialState: {
            ...baseState,
            mediasLoadedFromAlbumId: sampleAlbums[1].albumId,
            medias: [],
            mediasLoaded: true,
        },
    },
};

// Empty album — next album only (no previous)
export const NoMediasNextOnly: Story = {
    args: {
        initialState: {
            ...baseState,
            mediasLoadedFromAlbumId: sampleAlbums[0].albumId,
            medias: [],
            mediasLoaded: true,
        },
    },
};

// Empty album — previous album only (no next)
export const NoMediasPreviousOnly: Story = {
    args: {
        initialState: {
            ...baseState,
            mediasLoadedFromAlbumId: sampleAlbums[sampleAlbums.length - 1].albumId,
            medias: [],
            mediasLoaded: true,
        },
    },
};

// Empty album — alone, no neighbours
export const NoMediasAlone: Story = {
    args: {
        initialState: {
            ...loadedStateWithTwoAlbums,
            allAlbums: [sampleAlbums[0]],
            albums: [sampleAlbums[0]],
            mediasLoadedFromAlbumId: sampleAlbums[0].albumId,
            medias: [],
            mediasLoaded: true,
        },
    },
};

export const WithError: Story = {
    args: {
        initialState: {...loadedStateWithTwoAlbums, error: new Error('Failed to fetch medias from the server')},
    },
};
