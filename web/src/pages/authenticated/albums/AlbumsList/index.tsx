import {Box, List, ListItem, ListItemAvatar, ListItemSecondaryAction, ListItemText, Skeleton} from "@mui/material";
import {Album} from "../../../../core/catalog";
import {AlbumListEntry} from "./AlbumListEntry";
import ShareDialog from "../ShareDialog";
import React from "react";
import {useSharingModalController} from "../share-controller";

const AlbumsList = ({albums, loaded, selected}: {
    albums: Album[]
    loaded: boolean
    selected?: Album
}) => {
    // const {openSharingModal, open, error, onRevoke, sharedWith,  onClose, onGrant} = useSharingModalController()
    const {openSharingModal, ...shareDialogProps} = useSharingModalController()

    const isSelected = (a: Album, b?: Album) => a.albumId.owner === b?.albumId?.owner && a.albumId.folderName === b?.albumId?.folderName

    return (
        <Box sx={{overflow: 'auto'}}>
            <List>
                {!loaded && (Array.from(Array(4).keys()).map((v, index) => (
                    <ListItem
                        key={index}
                        divider={index < 2}>
                        <ListItemAvatar>
                            <Skeleton variant="circular" width={35} height={35} animation='wave'/>
                        </ListItemAvatar>
                        <ListItemText
                            primary={<Skeleton variant="text" width={200} animation='wave'/>}
                            secondary={<Skeleton variant="text" width={300} animation='wave'/>}
                        />
                        <ListItemSecondaryAction>
                            <Skeleton variant="circular" width={35} height={35} animation='wave'/>
                        </ListItemSecondaryAction>
                    </ListItem>
                )))}

                {albums.map((album, index) => (
                    <AlbumListEntry
                        key={album.albumId.owner + album.albumId.folderName}
                        album={album}
                        selected={isSelected(album, selected)}
                        onClickOnSharedWith={openSharingModal}
                    />
                ))}
            </List>

            <ShareDialog {...shareDialogProps} />
        </Box>
    )
}

export default AlbumsList
