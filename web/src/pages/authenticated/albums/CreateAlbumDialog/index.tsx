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


interface CreateAlbum {
    name: string
    start: Dayjs | null
    end: Dayjs | null
    forceFolderNameEnabled: boolean
    forceFolderName?: string
}

const saturdayTwoWeeksAgo = dayjs().startOf("week").subtract(8, "days")

const emptyCreateAlbum = (defaultDate: Dayjs): CreateAlbum => ({
    name: "",
    start: defaultDate,
    end: defaultDate.add(7, "days").endOf("day"),
    forceFolderNameEnabled: false
})

export default function CreateAlbumDialog({open, error, onClose, onSubmit, defaultDate = saturdayTwoWeeksAgo}: {
    open: boolean,
    onClose: () => void,
    error?: string,
    onSubmit: (album: CreateAlbum) => void,
    defaultDate?: Dayjs,
}) {
    const [startsAtStartOfTheDay, setStartsAtStartOfTheDay] = useState<boolean>(true)
    const [endsAtEndOfTheDay, setEndsAtEndOfTheDay] = useState<boolean>(true)
    const [withCustomFolderName, setWithCustomFolderName] = useState<boolean>(false)
    const [createForm, setCreateForm] = useState<CreateAlbum>(emptyCreateAlbum(defaultDate))

    const theme = useTheme();
    const isMobile = useMediaQuery(theme.breakpoints.down('sm'));

    const handleClose = useCallback(() => {
        onClose()
        setCreateForm(emptyCreateAlbum(defaultDate))
    }, [setCreateForm, onClose])

    const handleSubmit = useCallback(() => {
        onSubmit(createForm)
        setCreateForm(emptyCreateAlbum(defaultDate))
    }, [createForm, onSubmit])

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
                            variant={isMobile ? 'standard' : 'outlined'}
                            margin="dense"
                            size='medium'
                            id="email"
                            label="Name"
                            type="string"
                            onChange={(event: React.ChangeEvent<HTMLInputElement>) => setCreateForm(album => ({...album, name: event.target.value}))}
                            value={createForm.name}
                        />
                    </Grid>
                    <Grid xs={6}>
                        {startsAtStartOfTheDay && (
                            <DatePicker
                                label="First day"
                                value={createForm.start}
                                onChange={start => setCreateForm(form => ({...form, start}))}
                                renderInput={(params: any) => <TextField {...params} sx={{width: "100%"}}/>}
                            />) || (
                            <DateTimePicker
                                label="First day"
                                value={createForm.start}
                                onChange={start => setCreateForm(form => ({...form, start}))}
                                renderInput={(params: any) => <TextField {...params} sx={{width: "100%"}}/>}
                            />
                        )}
                    </Grid>
                    <Grid xs={6}>
                        <FormControlLabel control={<Checkbox checked={startsAtStartOfTheDay}
                                                             onChange={(event: React.ChangeEvent<HTMLInputElement>) => setStartsAtStartOfTheDay(event.target.checked)}/>}
                                          label="at the start of the day"/>
                    </Grid>
                    <Grid xs={6}>
                        {endsAtEndOfTheDay && (
                            <DatePicker
                                label="Last day"
                                value={createForm.end}
                                onChange={end => setCreateForm(form => ({...form, end}))}
                                renderInput={(params: any) => <TextField {...params} sx={{width: "100%"}}/>}
                            />) || (
                            <DateTimePicker
                                label="Last day"
                                value={createForm.end}
                                onChange={end => setCreateForm(form => ({...form, end}))}
                                renderInput={(params: any) => <TextField {...params} sx={{width: "100%"}}/>}
                            />
                        )}
                    </Grid>
                    <Grid xs={6}>
                        <FormControlLabel control={<Checkbox checked={endsAtEndOfTheDay}
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
                                <Checkbox checked={withCustomFolderName}
                                          onChange={(event: React.ChangeEvent<HTMLInputElement>) => setWithCustomFolderName(event.target.checked)}
                                />
                                <Divider sx={{height: 28, m: 0.5}} orientation="vertical"/>
                                <InputBase
                                    sx={{ml: 1, flex: 1}}
                                    placeholder="Custom folder name (ex: '/2025-08_Summer')"
                                    disabled={!withCustomFolderName}
                                />
                            </Paper>
                        </Tooltip>
                    </Grid>
                </Grid>
            </DialogContent>
            <DialogActions>
                <Button onClick={handleClose} color='info'>Cancel</Button>
                <Button onClick={handleSubmit} color='primary' variant='contained'>Save</Button>
            </DialogActions>
        </Dialog>
    );
}