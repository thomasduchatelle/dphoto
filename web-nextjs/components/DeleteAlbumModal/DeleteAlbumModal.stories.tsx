import type {Meta, StoryObj} from '@storybook/nextjs-vite';
import {fn, userEvent, within} from 'storybook/test';
import {DeleteAlbumModal} from './index';
import {Album, AlbumId} from '@/domains/catalog/language';

const createAlbumId = (owner: string, folderName: string): AlbumId => ({owner, folderName});

const clairObscurId = createAlbumId('sandfall', 'clair-obscur');

const sampleAlbums: Album[] = [
    {
        albumId: clairObscurId,
        name: 'Clair Obscur',
        start: new Date('2025-04-24'),
        end: new Date('2025-06-01'),
        totalCount: 47,
        temperature: 6.7,
        relativeTemperature: 1,
        sharedWith: [],
        covers: [],
    },
    {
        albumId: createAlbumId('sony', 'astro-bot'),
        name: 'Astro Bot',
        start: new Date('2024-09-06'),
        end: new Date('2024-09-12'),
        totalCount: 23,
        temperature: 12.1,
        relativeTemperature: 0.72,
        sharedWith: [],
        covers: [],
    },
    {
        albumId: createAlbumId('cdprojekt', 'the-witcher-3'),
        name: 'The Witcher 3: Wild Hunt — Blood and Wine',
        start: new Date('2024-05-30'),
        end: new Date('2024-06-14'),
        totalCount: 189,
        temperature: 9.4,
        relativeTemperature: 0.55,
        sharedWith: [],
        covers: [],
    },
];

const meta = {
    title: 'Catalog/DeleteAlbumModal',
    component: DeleteAlbumModal,
    parameters: {
        layout: 'fullscreen',
    },
    args: {
        albums: sampleAlbums,
        initialSelectedAlbumId: clairObscurId,
        isOpen: true,
        isLoading: false,
        error: undefined,
        onClose: fn(),
        onDelete: fn(),
    },
} satisfies Meta<typeof DeleteAlbumModal>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const NoPreselection: Story = {
    args: {
        initialSelectedAlbumId: undefined,
    },
};

export const ConfirmationStep: Story = {
    play: async () => {
        const dialog = within(document.body);
        const deleteButton = await dialog.findByRole('button', {name: 'Delete'});
        await userEvent.click(deleteButton);
    },
};

export const Loading: Story = {
    args: {
        isLoading: true,
    },
};

export const SubmissionError: Story = {
    args: {
        error: 'Unable to delete the album, please try again.',
    },
};
