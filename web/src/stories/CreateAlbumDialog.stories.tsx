import React from 'react';
import {AdapterDayjs} from "@mui/x-date-pickers/AdapterDayjs";
import {LocalizationProvider} from '@mui/x-date-pickers';
import dayjs from "dayjs";
import fr from "dayjs/locale/fr";
import {CreateAlbumDialog} from "../pages/authenticated/albums/CreateAlbumDialog";
import {Story, action} from "@ladle/react";
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

    return (
        <LocalizationProvider dateAdapter={AdapterDayjs} adapterLocale='fr'>
            <Button variant='contained' onClick={() => setOpen(true)}>
                Reopen Dialog
            </Button>
            <CreateAlbumDialog
                {...props}
                open={open}
                onClose={() => {
                    setOpen(false);
                    action("onClose")();
                }}
                onSubmit={async () => action("onSubmit")()}
                onNameChange={action("onNameChange")}
                onFolderNameChange={action("onFolderNameChange")}
                onWithCustomFolderNameChange={action("onWithCustomFolderNameChange")}
                onStartsAtStartOfTheDayChange={action("onStartsAtStartOfTheDayChange")}
                onEndsAtEndOfTheDayChange={action("onEndsAtEndOfTheDayChange")}
                onStartDateChange={action("onStartDateChange")}
                onEndDateChange={action("onEndDateChange")}
            />
        </LocalizationProvider>
    );
};

export const Default: Story<Props> = (props) => <CreateAlbumDialogLadle {...props} />

Default.args = {
    albumName: "",
    start: defaultStartDate,
    end: endDate,
    customFolderName: "",
    startsAtStartOfTheDay: true,
    endsAtEndOfTheDay: true,
    isCustomFolderNameEnabled: false,
    isLoading: false,
    canSubmit: true,
};
