'use client';

import {useCallback, useState} from 'react';
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
    SelectChangeEvent,
    Typography,
} from '@mui/material';
import {Album, AlbumId, albumIdEquals} from '@/domains/catalog/language';
import {DeleteDialogFrag} from '@/domains/catalog/album-delete/selector-deleteDialogSelector';
import {toLocaleDateWithDay} from '@/libs/dates';

export interface DeleteAlbumModalHandlers {
    onClose(): void;

    onDelete(albumId: AlbumId): void;
}

export type DeleteAlbumModalProps = DeleteDialogFrag & DeleteAlbumModalHandlers;

interface DeleteAlbumModalState {
    openedSignature: string;
    selectedAlbumIndex: number;
    confirmMode: boolean;
}

function initialSelectedIndex(albums: Album[], initialSelectedAlbumId?: AlbumId): number {
    if (!initialSelectedAlbumId) {
        return 0;
    }
    return Math.max(0, albums.findIndex(album => albumIdEquals(album.albumId, initialSelectedAlbumId)));
}

function openedSignatureFor(isOpen: boolean, initialSelectedAlbumId?: AlbumId): string {
    return `${isOpen ? '1' : '0'}:${initialSelectedAlbumId ? `${initialSelectedAlbumId.owner}/${initialSelectedAlbumId.folderName}` : ''}`;
}

function formatAlbumLabel(album: Album): string {
    let datePart = '';
    if (album.start && album.end) {
        const start = toLocaleDateWithDay(album.start);
        const end = toLocaleDateWithDay(album.end);
        datePart = `${start} → ${end}`;
    }
    return `${album.name} (${datePart} ; ${album.totalCount} media${album.totalCount <= 1 ? '' : 's'})`;
}

export function DeleteAlbumModal({
                                     albums,
                                     initialSelectedAlbumId,
                                     isOpen,
                                     isLoading,
                                     error,
                                     onClose,
                                     onDelete,
                                 }: DeleteAlbumModalProps) {
    const currentSignature = openedSignatureFor(isOpen, initialSelectedAlbumId);
    const [state, setState] = useState<DeleteAlbumModalState>(() => ({
        openedSignature: currentSignature,
        selectedAlbumIndex: initialSelectedIndex(albums, initialSelectedAlbumId),
        confirmMode: false,
    }));

    let {selectedAlbumIndex, confirmMode} = state;
    if (state.openedSignature !== currentSignature) {
        selectedAlbumIndex = initialSelectedIndex(albums, initialSelectedAlbumId);
        confirmMode = false;
        setState({openedSignature: currentSignature, selectedAlbumIndex, confirmMode});
    }

    const onSelectionChange = useCallback((event: SelectChangeEvent<number>) => {
        const newSelectedAlbumIndex = event.target.value as number;
        setState(prevState => ({
            ...prevState,
            selectedAlbumIndex: newSelectedAlbumIndex,
        }));
    }, []);

    const onSelectionConfirmed = useCallback(() => {
        setState(prevState => ({
            ...prevState,
            confirmMode: true,
        }));
    }, []);

    const onConfirmDelete = useCallback(() => {
        if (selectedAlbumIndex < albums.length) {
            onDelete(albums[selectedAlbumIndex].albumId);
        }
    }, [onDelete, selectedAlbumIndex, albums]);

    const selectedAlbum = selectedAlbumIndex < albums.length ? albums[selectedAlbumIndex] : undefined;

    return (
        <Dialog open={isOpen} onClose={onClose} maxWidth="sm" fullWidth>
            <Box sx={{height: '4px', marginTop: '0px !important'}}>
                {isLoading && (
                    <LinearProgress sx={{borderRadius: {sm: '4px 4px 0px 0px'}}}/>
                )}
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
                                    <MenuItem value={0} disabled>
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
                    <Typography sx={{mb: 2}}>
                        Are you sure you want to delete{' '}
                        <b>{selectedAlbum ? formatAlbumLabel(selectedAlbum) : ''}</b>
                        ?
                    </Typography>
                )}
            </DialogContent>
            <DialogActions>
                <Button
                    onClick={onClose}
                    color="info"
                    disabled={isLoading}
                    autoFocus
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
}
