'use client';

import React from 'react';
import {
    Alert,
    Box,
    Button,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    Grid,
    IconButton,
    LinearProgress,
    TextField,
    useMediaQuery,
    useTheme,
} from "@mui/material";
import {Close} from "@mui/icons-material";
import {LocalizationProvider} from "@mui/x-date-pickers";
import {AdapterDayjs} from "@mui/x-date-pickers/AdapterDayjs";
import {CreateDialogSelection} from "@/domains/catalog/album-create";
import {DateRangePicker} from "./DateRangePicker";
import {FolderNameInput} from "./FolderNameInput";

export interface CreateAlbumModalHandlers {
    onClose: () => void;
    onSubmit: () => Promise<void> | void;
    onNameChange: (name: string) => void;
    onFolderNameChange: (folderName: string) => void;
    onWithCustomFolderNameChange: (withCustom: boolean) => void;
    onStartsAtStartOfTheDayChange: (startsAtStart: boolean) => void;
    onEndsAtEndOfTheDayChange: (endsAtEnd: boolean) => void;
    onStartDateChange: (date: Date | null) => void;
    onEndDateChange: (date: Date | null) => void;
}

export type CreateAlbumModalProps = CreateDialogSelection & CreateAlbumModalHandlers;

export function CreateAlbumModal({
                                     open,
                                     albumName,
                                     start,
                                     end,
                                     customFolderName,
                                     startsAtStartOfTheDay,
                                     endsAtEndOfTheDay,
                                     isCustomFolderNameEnabled,
                                     isLoading,
                                     error,
                                     canSubmit,
                                     folderNameError,
                                     dateRangeError,
                                     nameError,
                                     onClose,
                                     onSubmit,
                                     onNameChange,
                                     onFolderNameChange,
                                     onWithCustomFolderNameChange,
                                     onStartsAtStartOfTheDayChange,
                                     onEndsAtEndOfTheDayChange,
                                     onStartDateChange,
                                     onEndDateChange,
                                 }: CreateAlbumModalProps) {
    const theme = useTheme();
    const isMobile = useMediaQuery(theme.breakpoints.down('sm'));

    const handleSubmit = (event: React.FormEvent) => {
        event.preventDefault();
        if (canSubmit) {
            void onSubmit();
        }
    };

    return (
        <LocalizationProvider dateAdapter={AdapterDayjs}>
            <Dialog
                open={open}
                onClose={onClose}
                fullWidth
                fullScreen={isMobile}
                maxWidth="md"
                aria-labelledby="create-album-dialog-title"
            >
                <Box sx={{height: '4px', marginTop: '0px !important'}}>
                    {isLoading && (
                        <LinearProgress sx={{borderRadius: {sm: '4px 4px 0px 0px'}}}/>
                    )}
                </Box>
                <DialogTitle id="create-album-dialog-title">Create an album</DialogTitle>
                <IconButton
                    aria-label="close"
                    onClick={onClose}
                    color="primary"
                    sx={{
                        position: 'absolute',
                        right: 8,
                        top: 8,
                        color: (theme) => theme.palette.grey[500],
                    }}
                >
                    <Close/>
                </IconButton>
                <form onSubmit={handleSubmit}>
                    <DialogContent>
                        <Grid container spacing={2} alignItems="center">
                            {error && (
                                <Grid size={{sm: 12, xs: 12}}>
                                    <Alert severity="error">{error}</Alert>
                                </Grid>
                            )}
                            <Grid size={{sm: 12, xs: 12}}>
                                <TextField
                                    autoFocus
                                    fullWidth
                                    label="Name"
                                    type="text"
                                    disabled={isLoading}
                                    onChange={(event) => onNameChange(event.target.value)}
                                    value={albumName}
                                    helperText={nameError}
                                    error={!!nameError}
                                />
                            </Grid>
                            <DateRangePicker
                                startDate={start}
                                endDate={end}
                                startAtDayStart={startsAtStartOfTheDay}
                                endAtDayEnd={endsAtEndOfTheDay}
                                onStartDateChange={onStartDateChange}
                                onEndDateChange={onEndDateChange}
                                onStartsAtStartOfTheDayChange={onStartsAtStartOfTheDayChange}
                                onEndsAtEndOfTheDayChange={onEndsAtEndOfTheDayChange}
                                disabled={isLoading}
                                dateError={!!dateRangeError}
                                dateHelperText={dateRangeError}
                            />
                            <Grid size={12}>
                                <FolderNameInput
                                    useCustomFolderName={isCustomFolderNameEnabled}
                                    value={customFolderName}
                                    placeholder="Custom folder name (ex: '/2025-08_Summer')"
                                    disabled={isLoading}
                                    onEnabledChange={onWithCustomFolderNameChange}
                                    onValueChange={onFolderNameChange}
                                    error={folderNameError}
                                />
                            </Grid>
                        </Grid>
                    </DialogContent>
                    <DialogActions>
                        <Button onClick={onClose} color="info">Cancel</Button>
                        <Button
                            type="submit"
                            color="primary"
                            variant="contained"
                            disabled={!canSubmit}
                        >
                            Create
                        </Button>
                    </DialogActions>
                </form>
            </Dialog>
        </LocalizationProvider>
    );
}
