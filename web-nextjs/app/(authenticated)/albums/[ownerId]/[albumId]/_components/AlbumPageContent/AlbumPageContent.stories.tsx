import type {Meta, StoryObj} from '@storybook/nextjs-vite';
import {AlbumPageContent} from './index';
import AppLayout from '@/components/AppLayout';
import {MediaType, MediaWithinADay} from '@/domains/catalog/language';
import {loadedStateWithTwoAlbums, sampleAlbums, someMediasByDays,} from '@/domains/catalog/tests/test-helper-state';

const multiDayMedias: MediaWithinADay[] = [
    {
        day: new Date(2025, 3, 25),
        medias: [
            {
                id: 'm1',
                type: MediaType.IMAGE,
                time: new Date('2025-04-25T09:00:00Z'),
                uiRelativePath: 'm1.jpg',
                contentPath: '/thumbnails/clair-obscur-1.jpg',
                source: ''
            },
            {
                id: 'm2',
                type: MediaType.IMAGE,
                time: new Date('2025-04-25T10:00:00Z'),
                uiRelativePath: 'm2.jpg',
                contentPath: '/thumbnails/clair-obscur-2.jpg',
                source: ''
            },
            {
                id: 'm3',
                type: MediaType.VIDEO,
                time: new Date('2025-04-25T11:00:00Z'),
                uiRelativePath: 'm3.mp4',
                contentPath: '/thumbnails/clair-obscur-3.jpg',
                source: ''
            },
            {
                id: 'm4',
                type: MediaType.IMAGE,
                time: new Date('2025-04-25T12:00:00Z'),
                uiRelativePath: 'm4.jpg',
                contentPath: '/thumbnails/clair-obscur-4.jpg',
                source: ''
            },
            {
                id: 'm5',
                type: MediaType.IMAGE,
                time: new Date('2025-04-25T13:00:00Z'),
                uiRelativePath: 'm5.jpg',
                contentPath: '/thumbnails/astro-bot-01.jpg',
                source: ''
            },
            {
                id: 'm6',
                type: MediaType.IMAGE,
                time: new Date('2025-04-25T14:00:00Z'),
                uiRelativePath: 'm6.jpg',
                contentPath: '/thumbnails/astro-bot-02.jpg',
                source: ''
            },
        ],
    },
    {
        day: new Date(2025, 4, 10),
        medias: [
            {
                id: 'm7',
                type: MediaType.IMAGE,
                time: new Date('2025-05-10T09:00:00Z'),
                uiRelativePath: 'm7.jpg',
                contentPath: '/thumbnails/the-witcher-3-01.jpg',
                source: ''
            },
            {
                id: 'm8',
                type: MediaType.IMAGE,
                time: new Date('2025-05-10T10:00:00Z'),
                uiRelativePath: 'm8.jpg',
                contentPath: '/thumbnails/the-witcher-3-02.jpg',
                source: ''
            },
            {
                id: 'm9',
                type: MediaType.VIDEO,
                time: new Date('2025-05-10T11:00:00Z'),
                uiRelativePath: 'm9.mp4',
                contentPath: '/thumbnails/the-witcher-3-03.jpg',
                source: ''
            },
            {
                id: 'm10',
                type: MediaType.IMAGE,
                time: new Date('2025-05-10T12:00:00Z'),
                uiRelativePath: 'm10.jpg',
                contentPath: '/thumbnails/death-stranding-1-01.jpg',
                source: ''
            },
            {
                id: 'm11',
                type: MediaType.IMAGE,
                time: new Date('2025-05-10T13:00:00Z'),
                uiRelativePath: 'm11.jpg',
                contentPath: '/thumbnails/death-stranding-1-02.jpg',
                source: ''
            },
        ],
    },
    {
        day: new Date(2025, 4, 28),
        medias: [
            {
                id: 'm12',
                type: MediaType.IMAGE,
                time: new Date('2025-05-28T10:00:00Z'),
                uiRelativePath: 'm12.jpg',
                contentPath: '/thumbnails/death-stranding-2-01.jpg',
                source: ''
            },
            {
                id: 'm13',
                type: MediaType.IMAGE,
                time: new Date('2025-05-28T11:00:00Z'),
                uiRelativePath: 'm13.jpg',
                contentPath: '/thumbnails/death-stranding-1-03.jpg',
                source: ''
            },
            {
                id: 'm14',
                type: MediaType.IMAGE,
                time: new Date('2025-05-28T14:00:00Z'),
                uiRelativePath: 'm14.jpg',
                contentPath: '/thumbnails/death-stranding-1-04.jpg',
                source: ''
            },
        ],
    },
];

const baseState = {
    ...loadedStateWithTwoAlbums,
    allAlbums: sampleAlbums,
    albums: sampleAlbums,
    mediasLoadedFromAlbumId: sampleAlbums[1].albumId,
    medias: multiDayMedias,
};

const meta = {
    title: 'Catalog/AlbumPageContent',
    component: AlbumPageContent,
    parameters: {layout: 'fullscreen'},
    decorators: [
        (Story: () => React.ReactNode) => (
            <AppLayout user={{
                name: 'Tony Stark',
                email: 'tony@stark-industries.com',
                picture: '/tonystark-profile.jpg',
                isOwner: true,
            }} logoutUrl="/auth/logout">
                <Story/>
            </AppLayout>
        ),
    ],
    args: {
        initialState: baseState,
    },
} satisfies Meta<typeof AlbumPageContent>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const NextAlbum: Story = {
    args: {
        initialState: {
            ...baseState,
            mediasLoadedFromAlbumId: sampleAlbums[0].albumId,
        },
    },
};

export const PreviousAlbum: Story = {
    args: {
        initialState: {
            ...baseState,
            mediasLoadedFromAlbumId: sampleAlbums[2].albumId,
            medias: someMediasByDays.map(d => ({
                ...d,
                medias: d.medias.map(m => ({...m, contentPath: '/thumbnails/the-witcher-3-01.jpg'})),
            })),
        },
    },
};

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

export const NoMedias: Story = {
    args: {
        initialState: {
            ...baseState,
            medias: [],
        },
    },
};

export const NoMediasNextOnly: Story = {
    args: {
        initialState: {
            ...baseState,
            mediasLoadedFromAlbumId: sampleAlbums[0].albumId,
            medias: [],
        },
    },
};

export const NoMediasPreviousOnly: Story = {
    args: {
        initialState: {
            ...baseState,
            mediasLoadedFromAlbumId: sampleAlbums[sampleAlbums.length - 1].albumId,
            medias: [],
        },
    },
};

export const NoMediasAlone: Story = {
    args: {
        initialState: {
            ...loadedStateWithTwoAlbums,
            allAlbums: [sampleAlbums[0]],
            albums: [sampleAlbums[0]],
            mediasLoadedFromAlbumId: sampleAlbums[0].albumId,
            medias: [],
        },
    },
};

export const WithError: Story = {
    args: {
        initialState: {...loadedStateWithTwoAlbums, error: new Error('Failed to fetch medias from the server')},
    },
};
