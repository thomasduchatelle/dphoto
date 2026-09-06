import type {Meta, StoryObj} from '@storybook/nextjs-vite';
import {AlbumPageContent} from './index';
import {AppBackground} from '@/components/AppLayout/AppBackground';
import {loadedStateWithTwoAlbums} from '@/domains/catalog/tests/test-helper-state';

const meta = {
    title: 'Catalog/AlbumPageContent',
    component: AlbumPageContent,
    parameters: {layout: 'fullscreen'},
    decorators: [(Story) => <AppBackground><Story/></AppBackground>],
    args: {
        initialState: loadedStateWithTwoAlbums,
    },
} satisfies Meta<typeof AlbumPageContent>;

export default meta;
type Story = StoryObj<typeof meta>;

export const WithMedias: Story = {};

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
