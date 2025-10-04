import React from 'react';
import {Alert, Box, Button, Dialog, DialogActions, DialogContent, DialogTitle, IconButton, LinearProgress, TextField, useMediaQuery, useTheme} from "@mui/material";
import {Close} from "@mui/icons-material";
import {EditNameDialogSelection} from "../../../core/catalog";
import {FolderNameInput} from "../FolderNameInput";

export interface EditNameDialogHandlers {
    onClose: () => void;
    onAlbumNameChange: (albumName: string) => void;
    onFolderNameEnabledChange: (enabled: boolean) => void;
    onFolderNameChange: (folderName: string) => void;
    onSave: () => void;
}

export function EditNameDialog({
    isOpen,
    onClose,
    albumName,
    customFolderName,
    isCustomFolderNameEnabled,
    isLoading,
    technicalError,
    nameError,
    folderNameError,
    isSaveEnabled,
    onAlbumNameChange,
    onFolderNameEnabledChange,
    onFolderNameChange,
    onSave,
}: EditNameDialogSelection & EditNameDialogHandlers) {
    const theme = useTheme();
    const isMobile = useMediaQuery(theme.breakpoints.down('sm'));

    return (
        <Dialog
            open={isOpen}
            onClose={onClose}
            fullWidth
            fullScreen={isMobile}
            maxWidth='md'
        >
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
            <DialogTitle>Edit Name</DialogTitle>
            <IconButton
                aria-label="close"
                onClick={onClose}
                color='primary'
                sx={{
                    position: 'absolute',
                    right: 8,
                    top: 8,
                    color: (theme) => theme.palette.grey[500],
                }}
            >
                <Close/>
            </IconButton>
            <DialogContent>
                <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, mt: 1 }}>
                    {technicalError && (
                        <Alert severity="error">
                            {technicalError}
                        </Alert>
                    )}
                    
                    <TextField
                        label="Name"
                        value={albumName}
                        onChange={(e) => onAlbumNameChange(e.target.value)}
                        error={!!nameError}
                        helperText={nameError}
                        disabled={isLoading}
                        fullWidth
                    />
                    
                    <FolderNameInput
                        useCustomFolderName={isCustomFolderNameEnabled}
                        value={customFolderName}
                        placeholder="Enter folder name"
                        disabled={isLoading}
                        onEnabledChange={onFolderNameEnabledChange}
                        onValueChange={onFolderNameChange}
                        tooltip="Override the automatically generated folder name"
                        error={folderNameError}
                    />
                </Box>
            </DialogContent>
            <DialogActions>
                <Button onClick={onClose} color='info' disabled={isLoading}>Cancel</Button>
                <Button onClick={onSave} variant="contained" disabled={!isSaveEnabled || isLoading}>Save</Button>
            </DialogActions>
        </Dialog>
    );
}
