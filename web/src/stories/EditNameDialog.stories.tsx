import React from 'react';
import {ComponentMeta, ComponentStory} from '@storybook/react';
import {EditNameDialog} from "../pages/authenticated/albums/EditNameDialog";

export default {
    title: 'Albums/EditNameDialog',
    component: EditNameDialog,
} as ComponentMeta<typeof EditNameDialog>;

const Template: ComponentStory<typeof EditNameDialog> = (args) => {
    const [albumName, setAlbumName] = React.useState(args.albumName);
    const [folderName, setFolderName] = React.useState(args.customFolderName);
    const [isFolderNameEnabled, setIsFolderNameEnabled] = React.useState(args.isCustomFolderNameEnabled);

    React.useEffect(() => {
        setAlbumName(args.albumName);
        setFolderName(args.customFolderName);
        setIsFolderNameEnabled(args.isCustomFolderNameEnabled);
    }, [args.albumName, args.customFolderName, args.isCustomFolderNameEnabled]);

    const albumNameError = albumName.trim() === "" ? "Album name cannot be blank" : undefined;
    const folderNameError = isFolderNameEnabled && folderName.trim() === "" ? "Folder name cannot be blank" : undefined;
    const isSaveEnabled = !albumNameError && !folderNameError;

    return (
        <EditNameDialog
            {...args}
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
        />
    );
};

export const Default = Template.bind({});
Default.args = {
    isOpen: true,
    albumName: "January 2025",
    customFolderName: "",
    isCustomFolderNameEnabled: false,
    isLoading: false,
    onClose: () => {},
    onSave: () => {},
};
Default.parameters = {
    delay: 300,
};

export const WithFolderNameEnabled = Template.bind({});
WithFolderNameEnabled.args = {
    isOpen: true,
    albumName: "Summer Vacation",
    customFolderName: "summer-vacation-2024",
    isCustomFolderNameEnabled: true,
    isLoading: false,
    onClose: () => {},
    onSave: () => {},
};
WithFolderNameEnabled.parameters = {
    delay: 300,
};

export const WithValidationErrors = Template.bind({});
WithValidationErrors.args = {
    isOpen: true,
    albumName: "",
    customFolderName: "",
    isCustomFolderNameEnabled: true,
    isLoading: false,
    onClose: () => {},
    onSave: () => {},
};
WithValidationErrors.parameters = {
    delay: 300,
};

export const Loading = Template.bind({});
Loading.args = {
    isOpen: true,
    albumName: "January 2025",
    customFolderName: "/january-2025",
    isCustomFolderNameEnabled: true,
    isLoading: true,
    onClose: () => {},
    onSave: () => {},
};
Loading.parameters = {
    storyshots: {disable: true},
};
