import React from 'react';
import {
    Alert,
    Box,
    Button,
    Checkbox,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    Divider,
    FormControlLabel,
    IconButton,
    InputBase,
    LinearProgress,
    Paper,
    TextField,
    Tooltip,
    useMediaQuery,
    useTheme
} from "@mui/material";
import Grid from '@mui/material/Unstable_Grid2';
import {Close} from "@mui/icons-material";
import {DatePicker, DateTimePicker} from '@mui/x-date-pickers';
import {albumFolderNameAlreadyTakenErr, albumStartAndEndDateMandatoryErr, CreateAlbumHandlers, CreateAlbumState} from "../../../../core/catalog";

export function CreateAlbumDialog({
                                      state,
                                      onCloseCreateAlbumDialog,
                                      onSubmitCreateAlbum,
                                      onNameChange,
                                      onFolderNameChange,
                                      onWithCustomFolderNameChange,
                                      onStartsAtStartOfTheDayChange,
                                      onEndsAtEndOfTheDayChange,
                                      onStartDateChange,
                                      onEndDateChange,
                                  }: {
    state: CreateAlbumState,
} & CreateAlbumHandlers) {
    const theme = useTheme();
    const isMobile = useMediaQuery(theme.breakpoints.down('sm'));

    const canBeSubmitted = state.name.length > 0 && !state.creationInProgress;
    const dateErrorArgs = state.errorCode === albumStartAndEndDateMandatoryErr ? {
        error: true,
        helperText: "Start and end dates are mandatory, and end date must be after the start date.",
    } : {};

    const errorMessage = getErrorMessage(state.errorCode);

    return (
        <Dialog
            open={state.open}
            onClose={onCloseCreateAlbumDialog}
            fullWidth
            fullScreen={isMobile}
            maxWidth='md'
        >
            <Box sx={{
                height: '4px',
                marginTop: '0px !important',
            }}>
                {state.creationInProgress && <LinearProgress sx={{
                    borderRadius: {
                        sm: '4px 4px 0px 0px'
                    },
                }}/>}
            </Box>
            <DialogTitle>Creates an album</DialogTitle>
            <IconButton
                aria-label="close"
                onClick={onCloseCreateAlbumDialog}
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
                            disabled={state.creationInProgress}
                            onChange={(event) => onNameChange(event.target.value)}
                            value={state.name}
                            helperText={state.errorCode === albumFolderNameAlreadyTakenErr && "The name must be unique (or the folder name must be explicitly set)"}
                            error={state.errorCode === albumFolderNameAlreadyTakenErr}
                        />
                    </Grid>
                    <Grid xs={6}>
                        {state.startsAtStartOfTheDay ? (
                            <DatePicker
                                label="First day"
                                disabled={state.creationInProgress}
                                value={state.start}
                                onChange={onStartDateChange}
                                renderInput={(params: any) => <TextField {...params} sx={{width: "100%"}} {...dateErrorArgs} helperText=''/>}
                            />) : (
                            <DateTimePicker
                                label="First day"
                                disabled={state.creationInProgress}
                                value={state.start}
                                onChange={onStartDateChange}
                                renderInput={(params: any) => <TextField {...params} sx={{width: "100%"}} {...dateErrorArgs} helperText=''/>}
                            />
                        )}
                    </Grid>
                    <Grid xs={6}>
                        <FormControlLabel control={<Checkbox checked={state.startsAtStartOfTheDay}
                                                             disabled={state.creationInProgress}
                                                             onChange={(event: React.ChangeEvent<HTMLInputElement>) => onStartsAtStartOfTheDayChange(event.target.checked)}/>}
                                          label="at the start of the day"/>
                    </Grid>
                    <Grid xs={6}>
                        {state.endsAtEndOfTheDay ? (
                            <DatePicker
                                label="Last day"
                                disabled={state.creationInProgress}
                                value={state.end}
                                onChange={onEndDateChange}
                                renderInput={(params: any) => <TextField {...params} sx={{width: "100%"}} {...dateErrorArgs} />}
                            />) : (
                            <DateTimePicker
                                label="Last day"
                                disabled={state.creationInProgress}
                                value={state.end}
                                onChange={onEndDateChange}
                                renderInput={(params: any) => <TextField {...params} sx={{width: "100%"}} {...dateErrorArgs} />}
                            />
                        )}
                    </Grid>
                    <Grid xs={6}>
                        <FormControlLabel control={<Checkbox checked={state.endsAtEndOfTheDay}
                                                             disabled={state.creationInProgress}
                                                             onChange={(event: React.ChangeEvent<HTMLInputElement>) => onEndsAtEndOfTheDayChange(event.target.checked)}/>}
                                          label="at the end of the day"/>
                    </Grid>
                    <Grid xs={12}>
                        <Tooltip title="The name of the physical folder name is generated from the date and the name; but can be overridden.">
                            <Paper
                                component="form"
                                sx={theme => ({
                                    p: '2px 4px',
                                    display: 'flex',
                                    alignItems: 'center',
                                    width: "99%",
                                    border: `solid 1px ${theme.palette.grey.A400}`
                                })}
                                elevation={0}
                            >
                                <Checkbox checked={state.withCustomFolderName}
                                          disabled={state.creationInProgress}
                                          onChange={(event: React.ChangeEvent<HTMLInputElement>) => onWithCustomFolderNameChange(event.target.checked)}
                                />
                                <Divider sx={{height: 28, m: 0.5}} orientation="vertical"/>
                                <InputBase
                                    sx={{ml: 1, flex: 1}}
                                    placeholder="Custom folder name (ex: '/2025-08_Summer')"
                                    disabled={!state.withCustomFolderName || state.creationInProgress}
                                    value={state.forceFolderName}
                                    onChange={(event) => onFolderNameChange(event.target.value)}
                                />
                            </Paper>
                        </Tooltip>
                    </Grid>
                </Grid>
            </DialogContent>
            <DialogActions>
                <Button onClick={onCloseCreateAlbumDialog} color='info'>Cancel</Button>
                <Button onClick={() => onSubmitCreateAlbum(state)} color='primary' variant='contained' disabled={!canBeSubmitted}>Save</Button>
            </DialogActions>
        </Dialog>
    );
}

function getErrorMessage(errorCode: string | undefined): string {
    switch (errorCode) {
        case undefined:
        case "":
        case albumFolderNameAlreadyTakenErr:
        case albumStartAndEndDateMandatoryErr:
            return "";

        default:
            return "Album couldn't be saved. Refresh your page and retry, or let the developer known.";
    }
}
