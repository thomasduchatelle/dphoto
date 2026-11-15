import type {Meta, StoryObj} from '@storybook/nextjs-vite';
import {TagsPageContentView} from './index';
import AppLayout from '@/components/AppLayout';
import {FAVORITE_TAG, TagMedia, tagPageReducer, tagPageSelector, initialTagPageState} from '@/domains/catalog/tags';
import {tagMediasLoaded} from '@/domains/catalog/tags';
import {MediaType} from '@/domains/catalog/language';

const favouriteMedias: TagMedia[] = [
    {
        id: 'co-1',
        type: MediaType.IMAGE,
        time: new Date('2025-04-25T09:00:00Z'),
        albumId: {owner: 'sandfall', folderName: 'clair-obscur'},
        tags: [FAVORITE_TAG],
        uiRelativePath: '/albums/sandfall/clair-obscur/co-1/clair-obscur-1.jpg',
        contentPath: '/thumbnails/clair-obscur-1.jpg',
        thumbnailUrl: '/thumbnails/clair-obscur-1.jpg',
        source: '',
    },
    {
        id: 'co-3',
        type: MediaType.IMAGE,
        time: new Date('2025-04-25T11:30:00Z'),
        albumId: {owner: 'sandfall', folderName: 'clair-obscur'},
        tags: [FAVORITE_TAG, 'people:luke'],
        uiRelativePath: '/albums/sandfall/clair-obscur/co-3/clair-obscur-3.jpg',
        contentPath: '/thumbnails/clair-obscur-3.jpg',
        thumbnailUrl: '/thumbnails/clair-obscur-3.jpg',
        source: 'Canon EOS R6',
    },
    {
        id: 'ab-2',
        type: MediaType.IMAGE,
        time: new Date('2024-09-06T14:12:00Z'),
        albumId: {owner: 'sony', folderName: 'astro-bot'},
        tags: [FAVORITE_TAG],
        uiRelativePath: '/albums/sony/astro-bot/ab-2/astro-bot-02.jpg',
        contentPath: '/thumbnails/astro-bot-02.jpg',
        thumbnailUrl: '/thumbnails/astro-bot-02.jpg',
        source: '',
    },
    {
        id: 'ds1-1',
        type: MediaType.IMAGE,
        time: new Date('2023-11-08T18:40:00Z'),
        albumId: {owner: 'kojima', folderName: 'death-stranding-1'},
        tags: [FAVORITE_TAG, 'place:school'],
        uiRelativePath: '/albums/kojima/death-stranding-1/ds1-1/death-stranding-1-01.jpg',
        contentPath: '/thumbnails/death-stranding-1-01.jpg',
        thumbnailUrl: '/thumbnails/death-stranding-1-01.jpg',
        source: 'iPhone 15 Pro',
    },
];

const selectionFor = (medias: TagMedia[]) =>
    tagPageSelector(
        tagPageReducer(initialTagPageState(FAVORITE_TAG), tagMediasLoaded({tag: FAVORITE_TAG, medias})),
    );

const meta = {
    title: 'catalog/TagsPageContent',
    component: TagsPageContentView,
    parameters: {layout: 'fullscreen'},
    decorators: [
        (Story: () => React.ReactNode) => (
            <AppLayout
                user={{
                    name: 'Tony Stark',
                    email: 'tony@stark-industries.com',
                    picture: '/tonystark-profile.jpg',
                    isOwner: true,
                }}
                logoutUrl="/auth/logout">
                <Story/>
            </AppLayout>
        ),
    ],
    args: selectionFor(favouriteMedias),
} satisfies Meta<typeof TagsPageContentView>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const Empty: Story = {
    args: selectionFor([]),
};

export const WithError: Story = {
    args: {
        ...selectionFor(favouriteMedias),
        error: new Error('Failed to reach the server'),
    },
};
