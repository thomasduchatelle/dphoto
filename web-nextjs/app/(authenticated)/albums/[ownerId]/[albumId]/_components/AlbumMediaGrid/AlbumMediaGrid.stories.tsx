import type {Meta, StoryObj} from '@storybook/nextjs-vite';
import {fn} from 'storybook/test';
import {Box} from '@mui/material';
import {AlbumMediaGrid} from './index';
import {AppBackground} from '@/components/AppLayout/AppBackground';
import {FAVORITE_TAG, MediaType, MediaWithinADay} from '@/domains/catalog/language';

const albumId = {owner: 'sandfall', folderName: 'clair-obscur'};

const oneDayMedias: MediaWithinADay[] = [{
    day: new Date(2025, 3, 25),
    medias: [
        {
            id: 'm1',
            type: MediaType.IMAGE,
            time: new Date('2025-04-25T09:00:00Z'),
            uiRelativePath: 'm1.jpg',
            contentPath: '/thumbnails/clair-obscur-1.jpg',
            thumbnailUrl: '/thumbnails/clair-obscur-1.jpg',
            source: '',
            tags: [FAVORITE_TAG],
        },
        {
            id: 'm2',
            type: MediaType.IMAGE,
            time: new Date('2025-04-25T10:00:00Z'),
            uiRelativePath: 'm2.jpg',
            contentPath: '/thumbnails/clair-obscur-2.jpg',
            thumbnailUrl: '/thumbnails/clair-obscur-2.jpg',
            source: '',
            tags: [],
        },
        {
            id: 'm3',
            type: MediaType.VIDEO,
            time: new Date('2025-04-25T11:00:00Z'),
            uiRelativePath: 'm3.mp4',
            contentPath: '/thumbnails/clair-obscur-3.jpg',
            thumbnailUrl: '/video-placeholder.png',
            source: '',
            tags: ['people:luke'],
        },
        {
            id: 'm4',
            type: MediaType.IMAGE,
            time: new Date('2025-04-25T12:00:00Z'),
            uiRelativePath: 'm4.jpg',
            contentPath: '/thumbnails/clair-obscur-4.jpg',
            thumbnailUrl: '/thumbnails/clair-obscur-4.jpg',
            source: '',
            tags: [FAVORITE_TAG, 'people:luke'],
        },
    ],
}];

const meta = {
    title: 'Catalog/AlbumMediaGrid',
    component: AlbumMediaGrid,
    parameters: {layout: 'fullscreen'},
    decorators: [(Story) => <AppBackground><Story/></AppBackground>],
    args: {
        albumId,
        medias: oneDayMedias,
        isOwner: true,
        onToggleFavorite: fn(),
    },
} satisfies Meta<typeof AlbumMediaGrid>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const StarsRevealed: Story = {
    decorators: [
        (Story) => (
            <Box sx={{'& .MediaThumbnail-star': {display: 'flex !important'}}}>
                <Story/>
            </Box>
        ),
    ],
};

export const Viewer: Story = {
    args: {
        isOwner: false,
    },
};
