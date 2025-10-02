import React from "react";
import {action, Story} from "@ladle/react";
import {DeleteAlbumDialog} from "../pages/authenticated/albums/DeleteAlbumDialog";
import {Button} from "@mui/material";
import {Album, AlbumId} from "../core/catalog";

export default {
    title: "Albums / DeleteAlbumDialog",
};

const makeAlbumId = (owner: string, folderName: string): AlbumId => ({
    owner,
    folderName,
});

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

type Props = React.ComponentProps<typeof DeleteAlbumDialog>;

const DeleteAlbumDialogWrapper: Story<Partial<Props>> = (props) => {
    const [open, setOpen] = React.useState(true);
    return (
        <>
            <Button onClick={() => setOpen(true)} variant="contained">
                Open Delete Dialog
            </Button>
            <DeleteAlbumDialog
                {...props as Props}
                isOpen={open}
                onClose={() => setOpen(false)}
                onDelete={action('onDelete')}
            />
        </>
    );
};

export const Default = (args: Props) => <DeleteAlbumDialogWrapper {...args} />
Default.args = {
    albums,
    initialSelectedAlbumId: albums[0].albumId,
    isLoading: false,
    error: undefined,
};

export const WithPreselectedAlbum = (args: Props) => <DeleteAlbumDialogWrapper {...args} />
WithPreselectedAlbum.args = {
    albums,
    initialSelectedAlbumId: albums[1].albumId,
    isLoading: false,
    error: undefined,
};

export const Loading = (args: Props) => <DeleteAlbumDialogWrapper {...args} />
Loading.args = {
    albums,
    initialSelectedAlbumId: albums[0].albumId,
    isLoading: true,
    error: undefined,
};

export const Error = (args: Props) => <DeleteAlbumDialogWrapper {...args} />
Error.args = {
    albums,
    initialSelectedAlbumId: albums[0].albumId,
    isLoading: false,
    error: "Something went wrong while deleting the album.",
};

export const NoAlbumsAvailable = (args: Props) => <DeleteAlbumDialogWrapper {...args} />
NoAlbumsAvailable.args = {
    albums: [],
    initialSelectedAlbumId: undefined,
    isLoading: false,
    error: undefined,
};