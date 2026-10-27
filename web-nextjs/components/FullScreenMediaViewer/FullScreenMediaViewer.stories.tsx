import type {Meta, StoryObj} from '@storybook/nextjs-vite';
import {fn} from 'storybook/test';
import {FullScreenMediaViewer} from './index';
import {Media, MediaType} from '@/domains/catalog/language';

const sampleMedias: Media[] = [
    {
        id: 'media-1',
        type: MediaType.IMAGE,
        time: new Date('2025-04-24T10:00:00Z'),
        uiRelativePath: 'media-1/clair-obscur-1.jpg',
        contentPath: '/thumbnails/clair-obscur-1.jpg',
        thumbnailUrl: '/thumbnails/clair-obscur-1.jpg',
        source: '',
    },
    {
        id: 'media-2',
        type: MediaType.IMAGE,
        time: new Date('2025-04-24T11:00:00Z'),
        uiRelativePath: 'media-2/clair-obscur-2.jpg',
        contentPath: '/thumbnails/clair-obscur-2.jpg',
        thumbnailUrl: '/thumbnails/clair-obscur-2.jpg',
        source: '',
    },
    {
        id: 'media-3',
        type: MediaType.IMAGE,
        time: new Date('2025-04-24T12:00:00Z'),
        uiRelativePath: 'media-3/clair-obscur-3.jpg',
        contentPath: '/thumbnails/clair-obscur-3.jpg',
        thumbnailUrl: '/thumbnails/clair-obscur-3.jpg',
        source: '',
    },
];

const meta = {
    title: 'catalog/FullScreenMediaViewer',
    component: FullScreenMediaViewer,
    parameters: {
        layout: 'fullscreen',
    },
    args: {
        medias: sampleMedias,
        displayedIndex: 0,
        onClose: fn(),
        onNavigate: fn(),
    },
} satisfies Meta<typeof FullScreenMediaViewer>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const Slideshow: Story = {
    args: {
        slideshow: true,
    },
};

export const SecondMedia: Story = {
    args: {
        displayedIndex: 1,
    },
};

export const LastMedia: Story = {
    args: {
        displayedIndex: 2,
    },
};

export const BrokenImage: Story = {
    args: {
        medias: [
            {
                id: 'broken-media',
                type: MediaType.IMAGE,
                time: new Date('2025-04-24T10:00:00Z'),
                uiRelativePath: 'broken-media/missing.jpg',
                contentPath: '/thumbnails/this-media-does-not-exist.jpg',
                thumbnailUrl: '/thumbnails/this-media-does-not-exist.jpg',
                source: '',
            },
        ],
        displayedIndex: 0,
    },
};
