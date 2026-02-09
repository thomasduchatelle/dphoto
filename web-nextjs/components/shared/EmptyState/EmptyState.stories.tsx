import type {Meta, StoryObj} from '@storybook/nextjs-vite';
import {fn} from 'storybook/test';
import {EmptyState} from './index';
import {CloudUpload as CloudUploadIcon, PhotoLibrary as PhotoLibraryIcon, SearchOff as SearchOffIcon} from '@mui/icons-material';

const meta = {
    title: 'Shared/EmptyState',
    component: EmptyState,
    parameters: {
        layout: 'centered',
    },
    tags: ['autodocs'],
} satisfies Meta<typeof EmptyState>;

export default meta;
type Story = StoryObj<typeof meta>;

export const NoAlbums: Story = {
    render: () => (
        <EmptyState
            icon={<PhotoLibraryIcon/>}
            title="No albums found"
            message="Create your first album to get started."
            action={{
                label: 'Create Album',
                onClick: fn(),
            }}
        />
    ),
};

export const NoMedia: Story = {
    render: () => (
        <EmptyState
            icon={<CloudUploadIcon/>}
            title="No photos in this album"
            message="Upload photos to see them here."
            action={{
                label: 'Upload Photos',
                onClick: fn(),
            }}
        />
    ),
};

export const NoAction: Story = {
    render: () => (
        <EmptyState
            icon={<SearchOffIcon/>}
            title="No results found"
            message="Try adjusting your search or filter criteria."
        />
    ),
};

export const WithIcon: Story = {
    render: () => (
        <EmptyState
            icon={<PhotoLibraryIcon/>}
            title="No albums"
            message="You don't have any albums yet."
        />
    ),
};

export const WithoutIcon: Story = {
    render: () => (
        <EmptyState
            title="Nothing to display"
            message="There is no content available at this time."
        />
    ),
};
