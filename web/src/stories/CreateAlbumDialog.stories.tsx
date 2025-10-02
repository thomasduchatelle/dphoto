import React from 'react';
import {AdapterDayjs} from "@mui/x-date-pickers/AdapterDayjs";
import {LocalizationProvider} from '@mui/x-date-pickers';
import dayjs from "dayjs";
import fr from "dayjs/locale/fr";
import {CreateAlbumDialog} from "../pages/authenticated/albums/CreateAlbumDialog";
import {action, Story} from "@ladle/react";
import {CreateDialogSelection} from "../core/catalog";
import {Button} from "@mui/material";

dayjs.locale(fr)

const defaultStartDate = new Date("2024-12-21T00:00:00Z")
const endDate = new Date("2024-12-29T23:59:00Z")

export default {
    title: 'Albums/CreateAlbumDialog',
};

type Props = Omit<CreateDialogSelection, "open">

const CreateAlbumDialogLadle = (props: Props) => {
    const [open, setOpen] = React.useState(true);
    const [albumName, setAlbumName] = React.useState(props.albumName);
    const [customFolderName, setCustomFolderName] = React.useState(props.customFolderName);
    const [isCustomFolderNameEnabled, setIsCustomFolderNameEnabled] = React.useState(props.isCustomFolderNameEnabled);
    const [startsAtStartOfTheDay, setStartsAtStartOfTheDay] = React.useState(props.startsAtStartOfTheDay);
    const [endsAtEndOfTheDay, setEndsAtEndOfTheDay] = React.useState(props.endsAtEndOfTheDay);
    const [start, setStart] = React.useState(props.start);
    const [end, setEnd] = React.useState(props.end);

    return (
        <LocalizationProvider dateAdapter={AdapterDayjs} adapterLocale='fr'>
            <Button variant='contained' onClick={() => setOpen(true)}>
                Reopen Dialog
            </Button>
            <CreateAlbumDialog
                {...props}
                albumName={albumName}
                customFolderName={customFolderName}
                isCustomFolderNameEnabled={isCustomFolderNameEnabled}
                startsAtStartOfTheDay={startsAtStartOfTheDay}
                endsAtEndOfTheDay={endsAtEndOfTheDay}
                start={start}
                end={end}
                open={open}
                onClose={() => setOpen(false)}
                onSubmit={async () => action("onSubmit")()}
                onNameChange={setAlbumName}
                onFolderNameChange={setCustomFolderName}
                onWithCustomFolderNameChange={setIsCustomFolderNameEnabled}
                onStartsAtStartOfTheDayChange={setStartsAtStartOfTheDay}
                onEndsAtEndOfTheDayChange={setEndsAtEndOfTheDay}
                onStartDateChange={setStart}
                onEndDateChange={setEnd}
            />
        </LocalizationProvider>
    );
};

export const Default: Story<Props> = (props) => <CreateAlbumDialogLadle {...props} />

Default.args = {
    albumName: "Return of the Jedi",
    start: defaultStartDate,
    end: endDate,
    customFolderName: "",
    startsAtStartOfTheDay: true,
    endsAtEndOfTheDay: true,
    isCustomFolderNameEnabled: false,
    isLoading: false,
    canSubmit: true,
};
