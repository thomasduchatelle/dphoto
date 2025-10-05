'use client';

import React, {useCallback, useEffect, useState} from "react";
import {
    Alert,
    Box,
    Button,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    FormControl,
    InputLabel,
    LinearProgress,
    MenuItem,
    Select,
    Typography
} from "@mui/material";
import {Album, AlbumId, albumIdEquals} from "../../../core/catalog";
import {toLocaleDateWithDay} from "../../../core/utils/date-utils";

interface DeleteAlbumDialogProps {
    albums: Album[];
    initialSelectedAlbumId?: AlbumId;
    isOpen: boolean;
    isLoading: boolean;
    error?: string;
    onDelete: (albumId: AlbumId) => void;
    onClose: () => void;
}

function formatAlbumLabel(album: Album) {
    let datePart = "";
    if (album.start && album.end) {
        const start = toLocaleDateWithDay(album.start);
        const end = toLocaleDateWithDay(album.end);
        datePart = `${start} → ${end}`;
    }
    return `${album.name} (${datePart} ; ${album.totalCount} media${album.totalCount <= 1 ? "" : "s"})`;
}

interface DeleteAlbumDialogState {
    wasOpen: boolean;
    selectedAlbumIndex: number;
    confirmMode: boolean;
}

export const DeleteAlbumDialog: React.FC<DeleteAlbumDialogProps> = ({
                                                                        albums,
                                                                        initialSelectedAlbumId,
                                                                        isOpen,
                                                                        isLoading,
                                                                        error,
                                                                        onDelete,
                                                                        onClose,
                                                                    }) => {
    const [{selectedAlbumIndex, confirmMode}, setState] = useState<DeleteAlbumDialogState>({
        wasOpen: false,
        selectedAlbumIndex: 0,
        confirmMode: false,
    });

    useEffect(() => {
        setState(prevState => prevState.wasOpen !== isOpen ? ({
            wasOpen: isOpen,
            selectedAlbumIndex: initialSelectedAlbumId ? albums.findIndex(album => albumIdEquals(album.albumId, initialSelectedAlbumId)) ?? 0 : 0,
            confirmMode: false,
        }) : prevState);
    }, [setState, isOpen, initialSelectedAlbumId, albums]);

    const onSelectionChange = useCallback((event: any) => {
        const newSelectedAlbumId = event.target.value as number;
        setState(prevState => ({
            ...prevState,
            selectedAlbumIndex: newSelectedAlbumId,
        }));
    }, [setState]);

    const onSelectionConfirmed = useCallback(() => {
        setState(prevState => ({
            ...prevState,
            confirmMode: true,
        }));
    }, [setState]);

    const onConfirmDelete = useCallback(() => {
        if (selectedAlbumIndex < albums.length) {
            onDelete(albums[selectedAlbumIndex].albumId);
        }
    }, [onDelete, selectedAlbumIndex, albums]);

    const selectedAlbum = selectedAlbumIndex < albums.length ? albums[selectedAlbumIndex] : undefined;

    return (
        <Dialog open={isOpen} onClose={onClose} maxWidth="sm" fullWidth>
            <Box sx={{
                height: '4px',
                marginTop: '0px !important',
            }}>
                {isLoading && <LinearProgress sx={{
                    borderRadius: {
                        sm: '4px 4px 0px 0px'
                    },
                }}/>}
            </Box>
            <DialogTitle>Delete an album</DialogTitle>
            <DialogContent>
                {error && (
                    <Alert severity="error" sx={{mb: 2}}>
                        {error}
                    </Alert>
                )}
                {!confirmMode ? (
                    <>
                        <Typography sx={{mb: 2}}>
                            What album do you want to delete? The medias will be re-assigned to the appropriate album.
                        </Typography>
                        <FormControl fullWidth sx={{mb: 2}}>
                            <InputLabel id="delete-album-select-label">Album</InputLabel>
                            <Select
                                labelId="delete-album-select-label"
                                value={selectedAlbumIndex}
                                label="Album"
                                onChange={onSelectionChange}
                                disabled={isLoading || albums.length === 0}
                                data-testid="album-select"
                            >
                                {albums.length === 0 ? (
                                    <MenuItem value="0" disabled>
                                        No albums available
                                    </MenuItem>
                                ) : (
                                    albums.map((album, index) => (
                                        <MenuItem key={index} value={index}>
                                            {formatAlbumLabel(album)}
                                        </MenuItem>
                                    ))
                                )}
                            </Select>
                        </FormControl>
                    </>
                ) : (
                    <>
                        <Typography sx={{mb: 2}}>
                            Are you sure you want to delete{" "}
                            <b>
                                {selectedAlbum ? formatAlbumLabel(selectedAlbum) : ""}
                            </b>
                            ?
                        </Typography>
                    </>
                )}
            </DialogContent>
            <DialogActions>
                <Button
                    onClick={onClose}
                    color="info"
                    disabled={isLoading}
                >
                    Cancel
                </Button>
                {!confirmMode ? (
                    <Button
                        onClick={onSelectionConfirmed}
                        color="error"
                        variant="contained"
                        disabled={!selectedAlbum || isLoading || albums.length === 0}
                    >
                        Delete
                    </Button>
                ) : (
                    <Button
                        onClick={onConfirmDelete}
                        color="error"
                        variant="contained"
                        disabled={isLoading}
                    >
                        Confirm
                    </Button>
                )}
            </DialogActions>
        </Dialog>
    );
};
