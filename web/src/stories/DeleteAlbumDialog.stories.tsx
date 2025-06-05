import React from "react";
import {ComponentMeta, ComponentStory} from "@storybook/react";
import {DeleteAlbumDialog} from "../pages/authenticated/albums/DeleteAlbumDialog";
import {Button} from "@mui/material";
import {Album, AlbumId} from "../core/catalog";

// Helper to create AlbumId
const makeAlbumId = (owner: string, folderName: string): AlbumId => ({
    owner,
    folderName,
});

// Helper to create Album
const makeAlbum = (
    owner: string,
    folderName: string,
    name: string,
    start: Date,
    end: Date,
    totalCount: number
): Album => ({
    albumId: makeAlbumId(owner, folderName),
    name,
    start,
    end,
    totalCount,
    temperature: 0,
    relativeTemperature: 0,
    sharedWith: [],
});

const albums: Album[] = [
    makeAlbum("user1", "summer-2023", "Summer 2023", new Date("2023-06-01"), new Date("2023-08-31"), 42),
    makeAlbum("user1", "winter-2022", "Winter 2022", new Date("2022-12-01"), new Date("2023-02-28"), 15),
    makeAlbum("user1", "spring-2024", "Spring 2024", new Date("2024-03-01"), new Date("2024-05-31"), 0),
];

export default {
    title: "Albums/DeleteAlbumDialog",
    component: DeleteAlbumDialog,
} as ComponentMeta<typeof DeleteAlbumDialog>;

const Template: ComponentStory<typeof DeleteAlbumDialog> = (args) => {
    const [open, setOpen] = React.useState(true);
    return (
        <>
            <Button onClick={() => setOpen(true)} variant="contained">
                Open Delete Dialog
            </Button>
            <DeleteAlbumDialog
                {...args}
                isOpen={open}
                onClose={() => setOpen(false)}
            />
        </>
    );
};

export const Default = Template.bind({});
Default.args = {
    albums,
    initialSelectedAlbumId: albums[0].albumId,
    isLoading: false,
    error: undefined,
    onDelete: (albumId: AlbumId) => alert(`Deleted album ${albumId.owner}/${albumId.folderName}`),
};
Default.parameters = {
    delay: 300,
};

export const WithPreselectedAlbum = Template.bind({});
WithPreselectedAlbum.args = {
    albums,
    initialSelectedAlbumId: albums[1].albumId,
    isLoading: false,
    error: undefined,
    onDelete: (albumId: AlbumId) => alert(`Deleted album ${albumId.owner}/${albumId.folderName}`),
};
WithPreselectedAlbum.parameters = {
    delay: 300,
};

export const Loading = Template.bind({});
Loading.args = {
    albums,
    initialSelectedAlbumId: albums[0].albumId,
    isLoading: true,
    error: undefined,
    onDelete: (albumId: AlbumId) => alert(`Deleted album ${albumId.owner}/${albumId.folderName}`),
};
Loading.parameters = {
    storyshots: {disable: true},
    delay: 300,
};

export const Error = Template.bind({});
Error.args = {
    albums,
    initialSelectedAlbumId: albums[0].albumId,
    isLoading: false,
    error: "Something went wrong while deleting the album.",
    onDelete: (albumId: AlbumId) => alert(`Deleted album ${albumId.owner}/${albumId.folderName}`),
};
Error.parameters = {
    delay: 300,
};

export const NoAlbumsAvailable = Template.bind({});
NoAlbumsAvailable.args = {
    albums: [],
    initialSelectedAlbumId: undefined,
    isLoading: false,
    error: undefined,
    onDelete: (albumId: AlbumId) => alert(`Deleted album ${albumId.owner}/${albumId.folderName}`),
};
NoAlbumsAvailable.parameters = {
    delay: 300,
};

// Note - the Confirm Mode cannot be tested because it's not possible to interact with the dialog in Storybook.