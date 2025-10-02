import React from 'react';
import {action, Story} from '@ladle/react';
import {Button} from '@mui/material';
import {EditNameDialog} from "../pages/authenticated/albums/EditNameDialog";

export default {
    title: 'Albums / EditNameDialog',
};

type Props = React.ComponentProps<typeof EditNameDialog>;

const EditNameDialogWrapper: Story<Partial<Props>> = (props) => {
    const [isOpen, setIsOpen] = React.useState(true);
    const [albumName, setAlbumName] = React.useState(props.albumName || '');
    const [folderName, setFolderName] = React.useState(props.customFolderName || '');
    const [isFolderNameEnabled, setIsFolderNameEnabled] = React.useState(props.isCustomFolderNameEnabled || false);

    React.useEffect(() => {
        setAlbumName(props.albumName || '');
        setFolderName(props.customFolderName || '');
        setIsFolderNameEnabled(props.isCustomFolderNameEnabled || false);
    }, [props.albumName, props.customFolderName, props.isCustomFolderNameEnabled]);

    const albumNameError = albumName.trim() === "" ? "Album name cannot be blank" : undefined;
    const folderNameError = isFolderNameEnabled && folderName.trim() === "" ? "Folder name cannot be blank" : undefined;
    const isSaveEnabled = !albumNameError && !folderNameError;

    return (
        <>
            <Button variant='contained' onClick={() => setIsOpen(true)}>
                Reopen Dialog
            </Button>
            <EditNameDialog
                {...props as Props}
                isOpen={isOpen}
                albumName={albumName}
                customFolderName={folderName}
                isCustomFolderNameEnabled={isFolderNameEnabled}
                nameError={albumNameError}
                folderNameError={folderNameError}
                isSaveEnabled={isSaveEnabled}
                onAlbumNameChange={setAlbumName}
                onFolderNameChange={setFolderName}
                onFolderNameEnabledChange={(enabled) => {
                    setIsFolderNameEnabled(enabled);
                    if (enabled) {
                        setFolderName("/vacation-photos");
                    } else {
                        setFolderName("");
                    }
                }}
                onClose={() => setIsOpen(false)}
                onSave={action('onSave')}
            />
        </>
    );
};

export const Default = (args: Props) => <EditNameDialogWrapper {...args} />
Default.args = {
    albumName: "January 2025",
    customFolderName: "",
    isCustomFolderNameEnabled: false,
    isLoading: false,
};

export const WithFolderNameEnabled = (args: Props) => <EditNameDialogWrapper {...args} />
WithFolderNameEnabled.args = {
    albumName: "Summer Vacation",
    customFolderName: "summer-vacation-2024",
    isCustomFolderNameEnabled: true,
    isLoading: false,
};

export const WithValidationErrors = (args: Props) => <EditNameDialogWrapper {...args} />
WithValidationErrors.args = {
    albumName: "",
    customFolderName: "",
    isCustomFolderNameEnabled: true,
    isLoading: false,
};

export const Loading = (args: Props) => <EditNameDialogWrapper {...args} />
Loading.args = {
    albumName: "January 2025",
    customFolderName: "/january-2025",
    isCustomFolderNameEnabled: true,
    isLoading: true,
};
