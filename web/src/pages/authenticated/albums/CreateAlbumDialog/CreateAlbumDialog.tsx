import React from 'react';
import {
    Alert,
    Box,
    Button,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    IconButton,
    LinearProgress,
    TextField,
    useMediaQuery,
    useTheme
} from "@mui/material";
import Grid from '@mui/material/Unstable_Grid2';
import {Close} from "@mui/icons-material";
import {CreateDialogSelection} from "../../../../core/catalog";
import {DateRangePicker} from "../DateRangePicker";
import {FolderNameInput} from "../FolderNameInput";

export interface CreateAlbumDialogHandlers {
    onClose: () => void;
    onSubmit: () => Promise<void>;
    onNameChange: (name: string) => void;
    onFolderNameChange: (folderName: string) => void;
    onWithCustomFolderNameChange: (withCustom: boolean) => void;
    onStartsAtStartOfTheDayChange: (startsAtStart: boolean) => void;
    onEndsAtEndOfTheDayChange: (endsAtEnd: boolean) => void;
    onStartDateChange: (date: Date | null) => void;
    onEndDateChange: (date: Date | null) => void;
}

export function CreateAlbumDialog({
                                      open,
                                      name,
                                      start,
                                      end,
                                      forceFolderName,
                                      startsAtStartOfTheDay,
                                      endsAtEndOfTheDay,
                                      withCustomFolderName,
                                      isLoading,
                                      error,
                                      canSubmit,
                                      onClose,
                                      onSubmit,
                                      onNameChange,
                                      onFolderNameChange,
                                      onWithCustomFolderNameChange,
                                      onStartsAtStartOfTheDayChange,
                                      onEndsAtEndOfTheDayChange,
                                      onStartDateChange,
                                      onEndDateChange,
                                  }: CreateDialogSelection & CreateAlbumDialogHandlers) {
    const theme = useTheme();
    const isMobile = useMediaQuery(theme.breakpoints.down('sm'));

    const dateError = error === "AlbumStartAndEndDateMandatoryErr";
    const dateHelperText = dateError ? "Start and end dates are mandatory, and end date must be after the start date." : "";

    const errorMessage = getErrorMessage(error);

    return (
        <Dialog
            open={open}
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
            <DialogTitle>Creates an album</DialogTitle>
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
                <Grid container spacing={2} alignItems='center'>
                    <Grid sm={12} xs={12}>
                        {errorMessage && <Alert severity="error">
                            {errorMessage}
                        </Alert>}
                    </Grid>
                    <Grid sm={12} xs={12}>
                        <TextField
                            autoFocus
                            fullWidth
                            label="Name"
                            type="string"
                            disabled={isLoading}
                            onChange={(event) => onNameChange(event.target.value)}
                            value={name}
                        />
                    </Grid>
                    <DateRangePicker
                        startDate={start || new Date()}
                        endDate={end || new Date()}
                        startAtDayStart={startsAtStartOfTheDay}
                        endAtDayEnd={endsAtEndOfTheDay}
                        onStartDateChange={onStartDateChange}
                        onEndDateChange={onEndDateChange}
                        onStartsAtStartOfTheDayChange={onStartsAtStartOfTheDayChange}
                        onEndsAtEndOfTheDayChange={onEndsAtEndOfTheDayChange}
                        disabled={isLoading}
                        dateError={dateError}
                        dateHelperText={dateHelperText}
                    />
                    <Grid xs={12}>
                        <FolderNameInput
                            useCustomFolderName={withCustomFolderName}
                            value={forceFolderName}
                            placeholder="Custom folder name (ex: '/2025-08_Summer')"
                            disabled={isLoading}
                            onEnabledChange={onWithCustomFolderNameChange}
                            onValueChange={onFolderNameChange}
                            error={error === "AlbumFolderNameAlreadyTakenErr" ? "The folder name must be unique" : undefined}
                        />
                    </Grid>
                </Grid>
            </DialogContent>
            <DialogActions>
                <Button onClick={onClose} color='info'>Cancel</Button>
                <Button onClick={onSubmit} color='primary' variant='contained' disabled={!canSubmit}>Save</Button>
            </DialogActions>
        </Dialog>
    );
}

function getErrorMessage(error: string | undefined): string {
    switch (error) {
        case undefined:
        case "":
        case "AlbumFolderNameAlreadyTakenErr":
        case "AlbumStartAndEndDateMandatoryErr":
            return "";

        default:
            return "Album couldn't be saved. Refresh your page and retry, or let the developer known.";
    }
}
