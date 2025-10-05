'use client';

import {Box, List, ListItem, ListItemAvatar, ListItemSecondaryAction, ListItemText, Skeleton} from "@mui/material";
import {Album, AlbumId, albumIdEquals} from "../../../core/catalog";
import {AlbumListEntry} from "./AlbumListEntry";
import React from "react";

const AlbumsList = ({albums, loaded, selectedAlbumId, openSharingModal}: {
    albums: Album[]
    loaded: boolean
    selectedAlbumId?: AlbumId
    openSharingModal: (albumId: AlbumId) => void
}) => {
    const isSelected = (album: Album) => albumIdEquals(selectedAlbumId, album.albumId)

    return (
        <Box sx={{overflow: 'auto'}}>
            <List>
                {!loaded ? (Array.from(Array(4).keys()).map((v, index) => (
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
                ))) : (
                    albums.map((album, index) => (
                        <AlbumListEntry
                            key={album.albumId.owner + album.albumId.folderName}
                            album={album}
                            selected={isSelected(album)}
                            onClickOnSharedWith={openSharingModal}
                        />
                    ))
                )}
            </List>
        </Box>
    )
}

export default AlbumsList
