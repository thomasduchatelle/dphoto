import {useState} from 'react';
import type {Meta, StoryObj} from '@storybook/nextjs-vite';
import {fn} from 'storybook/test';
import {Button} from '@mui/material';
import {EditNameDialog, EditNameDialogHandlers} from './index';
import {EditNameDialogSelection} from '@/domains/catalog/album-edit-name/selector-editNameDialogSelector';

type Props = EditNameDialogSelection & EditNameDialogHandlers;

const EditNameDialogWrapper = (props: Props) => {
    const [open, setOpen] = useState(true);
    const [albumName, setAlbumName] = useState(props.albumName);
    const [customFolderName, setCustomFolderName] = useState(props.customFolderName);
    const [isCustomFolderNameEnabled, setIsCustomFolderNameEnabled] = useState(props.isCustomFolderNameEnabled);

    return (
        <>
            <Button variant='contained' onClick={() => setOpen(true)}>
                Reopen Dialog
            </Button>
            <EditNameDialog
                {...props}
                isOpen={open}
                albumName={albumName}
                customFolderName={customFolderName}
                isCustomFolderNameEnabled={isCustomFolderNameEnabled}
                onClose={() => {
                    props.onClose();
                    setOpen(false);
                }}
                onAlbumNameChange={(name) => {
                    props.onAlbumNameChange(name);
                    setAlbumName(name);
                }}
                onFolderNameEnabledChange={(enabled) => {
                    props.onFolderNameEnabledChange(enabled);
                    setIsCustomFolderNameEnabled(enabled);
                }}
                onFolderNameChange={(name) => {
                    props.onFolderNameChange(name);
                    setCustomFolderName(name);
                }}
            />
        </>
    );
};

const meta = {
    title: 'Catalog/EditNameDialog',
    component: EditNameDialogWrapper,
    parameters: {layout: 'fullscreen'},
    args: {
        isOpen: true,
        albumName: 'Clair Obscur',
        originalName: 'Clair Obscur',
        customFolderName: '',
        isCustomFolderNameEnabled: false,
        isLoading: false,
        isSaveEnabled: true,
        onClose: fn(),
        onAlbumNameChange: fn(),
        onFolderNameEnabledChange: fn(),
        onFolderNameChange: fn(),
        onSave: fn(),
    },
} satisfies Meta<typeof EditNameDialogWrapper>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const LongName: Story = {
    args: {
        albumName: 'Voyage au bout du monde et retour en passant par les étoiles',
        originalName: 'Voyage au bout du monde et retour en passant par les étoiles',
    },
};

export const ValidationError: Story = {
    args: {
        albumName: '',
        originalName: 'Clair Obscur',
        nameError: 'Name is required',
        isSaveEnabled: false,
    },
};

export const SubmissionError: Story = {
    args: {
        technicalError: 'Failed to rename the album: the server is unreachable.',
    },
};

export const Loading: Story = {
    args: {
        isLoading: true,
        isSaveEnabled: false,
    },
};

export const CustomFolderName: Story = {
    args: {
        isCustomFolderNameEnabled: true,
        customFolderName: 'clair-obscur-2025',
    },
};
