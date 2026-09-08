import type {Meta, StoryObj} from '@storybook/nextjs-vite';
import {AlbumPageContent} from './index';
import {AppBackground} from '@/components/AppLayout/AppBackground';
import {
    loadedStateWithTwoAlbums,
    march2025,
    someMediasByDays,
    twoAlbums,
} from '@/domains/catalog/tests/test-helper-state';
import {MediaType} from '@/domains/catalog/language';

const multiDayMedias = [
    {
        day: new Date(2025, 0, 5),
        medias: [
            {id: 'm1', type: MediaType.IMAGE, time: new Date('2025-01-05T09:00:00Z'), uiRelativePath: 'm1.jpg', contentPath: '/tonystark-profile.jpg', source: ''},
            {id: 'm2', type: MediaType.IMAGE, time: new Date('2025-01-05T10:00:00Z'), uiRelativePath: 'm2.jpg', contentPath: '/tonystark-profile.jpg', source: ''},
            {id: 'm3', type: MediaType.VIDEO, time: new Date('2025-01-05T11:00:00Z'), uiRelativePath: 'm3.mp4', contentPath: '/tonystark-profile.jpg', source: ''},
            {id: 'm4', type: MediaType.IMAGE, time: new Date('2025-01-05T12:00:00Z'), uiRelativePath: 'm4.jpg', contentPath: '/tonystark-profile.jpg', source: ''},
            {id: 'm5', type: MediaType.IMAGE, time: new Date('2025-01-05T13:00:00Z'), uiRelativePath: 'm5.jpg', contentPath: '/tonystark-profile.jpg', source: ''},
            {id: 'm6', type: MediaType.IMAGE, time: new Date('2025-01-05T14:00:00Z'), uiRelativePath: 'm6.jpg', contentPath: '/tonystark-profile.jpg', source: ''},
        ],
    },
    {
        day: new Date(2025, 0, 8),
        medias: [
            {id: 'm7', type: MediaType.IMAGE, time: new Date('2025-01-08T09:00:00Z'), uiRelativePath: 'm7.jpg', contentPath: '/tonystark-profile.jpg', source: ''},
            {id: 'm8', type: MediaType.IMAGE, time: new Date('2025-01-08T10:00:00Z'), uiRelativePath: 'm8.jpg', contentPath: '/tonystark-profile.jpg', source: ''},
            {id: 'm9', type: MediaType.VIDEO, time: new Date('2025-01-08T11:00:00Z'), uiRelativePath: 'm9.mp4', contentPath: '/tonystark-profile.jpg', source: ''},
            {id: 'm10', type: MediaType.IMAGE, time: new Date('2025-01-08T12:00:00Z'), uiRelativePath: 'm10.jpg', contentPath: '/tonystark-profile.jpg', source: ''},
            {id: 'm11', type: MediaType.IMAGE, time: new Date('2025-01-08T13:00:00Z'), uiRelativePath: 'm11.jpg', contentPath: '/tonystark-profile.jpg', source: ''},
        ],
    },
    {
        day: new Date(2025, 0, 12),
        medias: [
            {id: 'm12', type: MediaType.IMAGE, time: new Date('2025-01-12T10:00:00Z'), uiRelativePath: 'm12.jpg', contentPath: '/tonystark-profile.jpg', source: ''},
            {id: 'm13', type: MediaType.IMAGE, time: new Date('2025-01-12T11:00:00Z'), uiRelativePath: 'm13.jpg', contentPath: '/tonystark-profile.jpg', source: ''},
            {id: 'm14', type: MediaType.IMAGE, time: new Date('2025-01-12T14:00:00Z'), uiRelativePath: 'm14.jpg', contentPath: '/tonystark-profile.jpg', source: ''},
        ],
    },
];

const threeAlbums = [...twoAlbums, march2025];

const stateWithThreeAlbums = {
    ...loadedStateWithTwoAlbums,
    allAlbums: threeAlbums,
    albums: threeAlbums,
    mediasLoadedFromAlbumId: twoAlbums[1].albumId,
    medias: multiDayMedias,
};

const meta = {
    title: 'Catalog/AlbumPageContent',
    component: AlbumPageContent,
    parameters: {layout: 'fullscreen'},
    decorators: [(Story: () => React.ReactNode) => <AppBackground><Story/></AppBackground>],
    args: {
        initialState: stateWithThreeAlbums,
    },
} satisfies Meta<typeof AlbumPageContent>;

export default meta;
type Story = StoryObj<typeof meta>;

export const WithMedias: Story = {};

export const NextAlbum: Story = {
    args: {
        initialState: {
            ...stateWithThreeAlbums,
            mediasLoadedFromAlbumId: twoAlbums[0].albumId,
            medias: multiDayMedias,
        },
    },
};

export const PreviousAlbum: Story = {
    args: {
        initialState: {
            ...stateWithThreeAlbums,
            mediasLoadedFromAlbumId: twoAlbums[1].albumId,
            medias: someMediasByDays,
        },
    },
};

export const SingleAlbum: Story = {
    args: {
        initialState: {
            ...loadedStateWithTwoAlbums,
            allAlbums: [twoAlbums[0]],
            albums: [twoAlbums[0]],
            mediasLoadedFromAlbumId: twoAlbums[0].albumId,
            medias: multiDayMedias,
        },
    },
};

export const NoMedias: Story = {
    args: {
        initialState: {...loadedStateWithTwoAlbums, medias: [], mediasLoaded: true},
    },
};

export const WithError: Story = {
    args: {
        initialState: {...loadedStateWithTwoAlbums, error: new Error('Failed to fetch medias from the server')},
    },
};
