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
import React, {useCallback, useState} from "react";
import Grid from '@mui/material/Unstable_Grid2';
import {Close} from "@mui/icons-material";
import dayjs, {Dayjs} from 'dayjs';
import {DatePicker, DateTimePicker} from '@mui/x-date-pickers'

export const albumFolderNameAlreadyTakenErr = "AlbumFolderNameAlreadyTakenErr";
export const albumStartAndEndDateMandatoryErr = "AlbumStartAndEndDateMandatoryErr";

interface CreateAlbumDialogState {
    name: string
    start: Dayjs | null
    end: Dayjs | null
    forceFolderName: string
    startsAtStartOfTheDay: boolean
    endsAtEndOfTheDay: boolean
    withCustomFolderName: boolean
    errorCode?: string
}

const saturdayTwoWeeksAgo = dayjs().startOf("week").subtract(8, "days")

const emptyCreateAlbum = (defaultDate: Dayjs): CreateAlbumDialogState => ({
    name: "",
    start: defaultDate,
    end: defaultDate.add(7, "days").endOf("day"),
    forceFolderName: "",
    startsAtStartOfTheDay: true,
    endsAtEndOfTheDay: true,
    withCustomFolderName: false,
})

export default function CreateAlbumDialog({open, onClose, onSubmit, defaultDate = saturdayTwoWeeksAgo, defaultErrorCode}: {
    open: boolean,
    onClose: () => void,
    onSubmit: (album: CreateAlbumDialogState) => void,
    defaultDate?: Dayjs,
    defaultErrorCode?: string
}) {
    const [state, setState] = useState<CreateAlbumDialogState>({...emptyCreateAlbum(defaultDate), errorCode: defaultErrorCode})
    const setStartsAtStartOfTheDay = (startsAtStartOfTheDay: boolean) => setState(prev => ({...prev, startsAtStartOfTheDay}))
    const setEndsAtEndOfTheDay = (endsAtEndOfTheDay: boolean) => setState(prev => ({...prev, endsAtEndOfTheDay}))
    const setWithCustomFolderName = (withCustomFolderName: boolean) => setState(prev => ({...prev, withCustomFolderName}))

    const theme = useTheme();
    const isMobile = useMediaQuery(theme.breakpoints.down('sm'));

    const canBeSubmitted = state.name.length > 0
    const dateErrorArgs = state.errorCode === albumStartAndEndDateMandatoryErr ? {
        error: true,
        helperText: "Start and end dates are mandatory, and end date must be after the start date.",
    } : {}

    const handleOnNameChange = useCallback((event: React.ChangeEvent<HTMLInputElement>) => {
        setState(prev => ({...prev, name: event.target.value, errorCode: prev.errorCode !== albumFolderNameAlreadyTakenErr ? prev.errorCode : undefined}))
    }, [setState])
    const handleStartDateChange = useCallback((start: Dayjs | null) => setState(prev => ({
        ...prev,
        start,
        errorCode: prev.errorCode === albumStartAndEndDateMandatoryErr ? undefined : prev.errorCode
    })), [setState])
    const handleEndDateChange = useCallback((end: Dayjs | null) => setState(prev => ({
        ...prev,
        end,
        errorCode: prev.errorCode === albumStartAndEndDateMandatoryErr ? undefined : prev.errorCode
    })), [setState])
    const handleOnFolderNameChange = useCallback((event: React.ChangeEvent<HTMLInputElement>) => {
        setState(prev => ({
            ...prev,
            forceFolderName: event.target.value,
            errorCode: prev.errorCode !== albumFolderNameAlreadyTakenErr ? prev.errorCode : undefined
        }))
    }, [setState])

    const handleClose = useCallback(() => {
        onClose()
        setState(emptyCreateAlbum(defaultDate))
    }, [setState, onClose, defaultDate])

    const handleSubmit = useCallback(() => {
        onSubmit(state)
        setState(emptyCreateAlbum(defaultDate))
    }, [state, setState, onSubmit, defaultDate])

    return (
        <Dialog
            open={open}
            onClose={handleClose}
            fullWidth
            fullScreen={isMobile}
            maxWidth='md'
        >
            <DialogTitle>Creates an album</DialogTitle>
            <IconButton
                aria-label="close"
                onClick={handleClose}
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
                {/*{error?.type === "general" && (*/}
                {/*    <Alert severity='error' sx={theme => ({mb: theme.spacing(2)})}>{error.message}</Alert>*/}
                {/*)}*/}
                <Grid container spacing={2} alignItems='center'>
                    <Grid sm={12} xs={12}>
                        <TextField
                            autoFocus
                            fullWidth
                            label="Name"
                            type="string"
                            onChange={handleOnNameChange}
                            value={state.name}
                            helperText={state.errorCode === albumFolderNameAlreadyTakenErr && "The name must be unique (or the folder name must be explicitely set)"}
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
                                renderInput={(params: any) => <TextField {...params} sx={{width: "100%"}} {...dateErrorArgs}/>}
                            />) : (
                            <DateTimePicker
                                label="Last day"
                                value={state.end}
                                onChange={handleEndDateChange}
                                renderInput={(params: any) => <TextField {...params} sx={{width: "100%"}} {...dateErrorArgs}/>}
                            />
                        )}
                    </Grid>
                    <Grid xs={6}>
                        <FormControlLabel control={<Checkbox checked={state.endsAtEndOfTheDay}
                                                             onChange={(event: React.ChangeEvent<HTMLInputElement>) => setEndsAtEndOfTheDay(event.target.checked)}/>}
                                          label="at the end of the day"/>
                    </Grid>
                    <Grid xs={12}>
                        <Tooltip title="The name of the physical foldername is generated from the date and the name ; but can be overriden.">
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
                <Button onClick={handleClose} color='info'>Cancel</Button>
                <Button onClick={handleSubmit} color='primary' variant='contained' disabled={!canBeSubmitted}>Save</Button>
            </DialogActions>
        </Dialog>
    );
}