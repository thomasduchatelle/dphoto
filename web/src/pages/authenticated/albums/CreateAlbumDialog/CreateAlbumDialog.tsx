import React from 'react';
import {
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
    Paper,
    TextField,
    Tooltip,
    useMediaQuery,
    useTheme
} from "@mui/material";
import Grid from '@mui/material/Unstable_Grid2';
import {Close} from "@mui/icons-material";
import {DatePicker, DateTimePicker} from '@mui/x-date-pickers';
import {Dayjs} from 'dayjs';
import {CreateAlbumState} from "../../../../core/catalog/domain/CreateAlbumController";

export const albumFolderNameAlreadyTakenErr = "AlbumFolderNameAlreadyTakenErr";
export const albumStartAndEndDateMandatoryErr = "AlbumStartAndEndDateMandatoryErr";

export function CreateAlbumDialog({
                                              state,
                                              onClose,
                                              onSubmit,
                                              setStartsAtStartOfTheDay,
                                              setEndsAtEndOfTheDay,
                                              setWithCustomFolderName,
                                              handleOnNameChange,
                                              handleStartDateChange,
                                              handleEndDateChange,
                                              handleOnFolderNameChange
                                          }: {
    state: CreateAlbumState,
    onClose: () => void,
    onSubmit: (state: CreateAlbumState) => void,
    setStartsAtStartOfTheDay: (startsAtStartOfTheDay: boolean) => void,
    setEndsAtEndOfTheDay: (endsAtEndOfTheDay: boolean) => void,
    setWithCustomFolderName: (withCustomFolderName: boolean) => void,
    handleOnNameChange: (event: React.ChangeEvent<HTMLInputElement>) => void,
    handleStartDateChange: (start: Dayjs | null) => void,
    handleEndDateChange: (end: Dayjs | null) => void,
    handleOnFolderNameChange: (event: React.ChangeEvent<HTMLInputElement>) => void
}) {
    const theme = useTheme();
    const isMobile = useMediaQuery(theme.breakpoints.down('sm'));

    const canBeSubmitted = state.name.length > 0;
    const dateErrorArgs = state.errorCode === albumStartAndEndDateMandatoryErr ? {
        error: true,
        helperText: "Start and end dates are mandatory, and end date must be after the start date.",
    } : {};

    return (
        <Dialog
            open={state.open}
            onClose={onClose}
            fullWidth
            fullScreen={isMobile}
            maxWidth='md'
        >
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
                        <TextField
                            autoFocus
                            fullWidth
                            label="Name"
                            type="string"
                            onChange={handleOnNameChange}
                            value={state.name}
                            helperText={state.errorCode === albumFolderNameAlreadyTakenErr && "The name must be unique (or the folder name must be explicitly set)"}
                            error={state.errorCode === albumFolderNameAlreadyTakenErr}
                        />
                    </Grid>
                    <Grid xs={6}>
                        {state.startsAtStartOfTheDay ? (
                            <DatePicker
                                label="First day"
                                value={state.start}
                                onChange={handleStartDateChange}
                                renderInput={(params: any) => <TextField {...params} sx={{width: "100%"}} {...dateErrorArgs} helperText=''/>}
                            />) : (
                            <DateTimePicker
                                label="First day"
                                value={state.start}
                                onChange={handleStartDateChange}
                                renderInput={(params: any) => <TextField {...params} sx={{width: "100%"}} {...dateErrorArgs} helperText=''/>}
                            />
                        )}
                    </Grid>
                    <Grid xs={6}>
                        <FormControlLabel control={<Checkbox checked={state.startsAtStartOfTheDay}
                                                             onChange={(event: React.ChangeEvent<HTMLInputElement>) => setStartsAtStartOfTheDay(event.target.checked)}/>}
                                          label="at the start of the day"/>
                    </Grid>
                    <Grid xs={6}>
                        {state.endsAtEndOfTheDay ? (
                            <DatePicker
                                label="Last day"
                                value={state.end}
                                onChange={handleEndDateChange}
                                renderInput={(params: any) => <TextField {...params} sx={{width: "100%"}} {...dateErrorArgs} />}
                            />) : (
                            <DateTimePicker
                                label="Last day"
                                value={state.end}
                                onChange={handleEndDateChange}
                                renderInput={(params: any) => <TextField {...params} sx={{width: "100%"}} {...dateErrorArgs} />}
                            />
                        )}
                    </Grid>
                    <Grid xs={6}>
                        <FormControlLabel control={<Checkbox checked={state.endsAtEndOfTheDay}
                                                             onChange={(event: React.ChangeEvent<HTMLInputElement>) => setEndsAtEndOfTheDay(event.target.checked)}/>}
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
                                          onChange={(event: React.ChangeEvent<HTMLInputElement>) => setWithCustomFolderName(event.target.checked)}
                                />
                                <Divider sx={{height: 28, m: 0.5}} orientation="vertical"/>
                                <InputBase
                                    sx={{ml: 1, flex: 1}}
                                    placeholder="Custom folder name (ex: '/2025-08_Summer')"
                                    disabled={!state.withCustomFolderName}
                                    value={state.forceFolderName}
                                    onChange={handleOnFolderNameChange}
                                />
                            </Paper>
                        </Tooltip>
                    </Grid>
                </Grid>
            </DialogContent>
            <DialogActions>
                <Button onClick={onClose} color='info'>Cancel</Button>
                <Button onClick={() => onSubmit(state)} color='primary' variant='contained' disabled={!canBeSubmitted}>Save</Button>
            </DialogActions>
        </Dialog>
    );
}