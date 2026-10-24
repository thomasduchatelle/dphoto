import type {Meta, StoryObj} from '@storybook/nextjs-vite';
import {fn} from 'storybook/test';
import {CreateAlbumModal} from './index';

const meta = {
    title: 'Catalog/CreateAlbumModal',
    component: CreateAlbumModal,
    parameters: {
        layout: 'fullscreen',
    },
    args: {
        open: true,
        albumName: '',
        customFolderName: '',
        isCustomFolderNameEnabled: false,
        start: new Date('2025-04-26T00:00:00.000Z'),
        end: new Date('2025-05-01T00:00:00.000Z'),
        startsAtStartOfTheDay: true,
        endsAtEndOfTheDay: true,
        isLoading: false,
        canSubmit: false,
        onClose: fn(),
        onSubmit: fn(),
        onNameChange: fn(),
        onFolderNameChange: fn(),
        onWithCustomFolderNameChange: fn(),
        onStartsAtStartOfTheDayChange: fn(),
        onEndsAtEndOfTheDayChange: fn(),
        onStartDateChange: fn(),
        onEndDateChange: fn(),
    },
} satisfies Meta<typeof CreateAlbumModal>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const FilledValid: Story = {
    args: {
        albumName: 'Clair Obscur',
        canSubmit: true,
    },
};

export const InvalidDatesError: Story = {
    args: {
        albumName: 'Clair Obscur',
        start: new Date('2025-05-10T00:00:00.000Z'),
        end: new Date('2025-05-01T00:00:00.000Z'),
        dateRangeError: 'The end date cannot be before the start date',
        canSubmit: false,
    },
};

export const CustomFolderName: Story = {
    args: {
        albumName: 'Clair Obscur',
        isCustomFolderNameEnabled: true,
        customFolderName: '2025-04_clair-obscur',
        canSubmit: true,
    },
};

export const SubmissionError: Story = {
    args: {
        albumName: 'Clair Obscur',
        error: "The album couldn't be created, try again later.",
        canSubmit: true,
    },
};

export const NameAlreadyExists: Story = {
    args: {
        albumName: 'Clair Obscur',
        nameError: 'An album with this folder name already exists.',
        canSubmit: false,
    },
};

export const FolderNameAlreadyExists: Story = {
    args: {
        albumName: 'Clair Obscur',
        isCustomFolderNameEnabled: true,
        customFolderName: '2025-04_clair-obscur',
        folderNameError: 'An album with this folder name already exists.',
        canSubmit: false,
    },
};

export const Loading: Story = {
    args: {
        albumName: 'Clair Obscur',
        isLoading: true,
        canSubmit: false,
    },
};
