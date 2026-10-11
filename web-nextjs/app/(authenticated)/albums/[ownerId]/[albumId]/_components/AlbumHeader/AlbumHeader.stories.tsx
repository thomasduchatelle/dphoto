import type {Meta, StoryObj} from '@storybook/nextjs-vite';
import {fn} from 'storybook/test';
import {AlbumHeader} from './index';
import {AppBackground} from '@/components/AppLayout/AppBackground';
import {sampleAlbums} from '@/domains/catalog/tests/test-helper-state';

const meta = {
    title: 'Catalog/AlbumHeader',
    component: AlbumHeader,
    parameters: {layout: 'fullscreen'},
    decorators: [
        (Story) => (
            <AppBackground>
                <Story/>
            </AppBackground>
        ),
    ],
    args: {
        album: sampleAlbums[0],
        isOwner: true,
        onShare: fn(),
        onRename: fn(),
        onEditDates: fn(),
        onDelete: fn(),
    },
} satisfies Meta<typeof AlbumHeader>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const Viewer: Story = {
    args: {
        isOwner: false,
    },
};

export const NoAlbum: Story = {
    args: {
        album: undefined,
    },
};
