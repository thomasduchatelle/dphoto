import React from "react";
import {action, Story} from "@ladle/react";
import {EditDatesDialog} from "../pages/authenticated/albums/EditDatesDialog";
import {Button} from "@mui/material";
import {AdapterDayjs} from "@mui/x-date-pickers/AdapterDayjs";
import {LocalizationProvider} from "@mui/x-date-pickers";
import dayjs from "dayjs";
import fr from "dayjs/locale/fr";

export default {
    title: "Albums / EditDatesDialog",
};

dayjs.locale(fr)

type Props = React.ComponentProps<typeof EditDatesDialog>;

const EditDatesDialogWrapper: Story<Partial<Props>> = (props) => {
    const [open, setOpen] = React.useState(true);
    const [startDate, setStartDate] = React.useState<Date | null>(props.startDate || null);
    const [endDate, setEndDate] = React.useState<Date | null>(props.endDate || null);
    const [startAtDayStart, setStartAtDayStart] = React.useState(props.startAtDayStart || false);
    const [endAtDayEnd, setEndAtDayEnd] = React.useState(props.endAtDayEnd || false);

    React.useEffect(() => {
        setStartDate(props.startDate || null);
        setEndDate(props.endDate || null);
        setStartAtDayStart(props.startAtDayStart || false);
        setEndAtDayEnd(props.endAtDayEnd || false);
    }, [props.startDate, props.endDate, props.startAtDayStart, props.endAtDayEnd]);

    return (
        <LocalizationProvider dateAdapter={AdapterDayjs} adapterLocale='fr'>
            <Button onClick={() => setOpen(true)} variant="contained">
                Open Edit Dates Dialog
            </Button>
            <EditDatesDialog
                {...props as Props}
                isOpen={open}
                onClose={() => setOpen(false)}
                startDate={startDate!}
                endDate={endDate!}
                startAtDayStart={startAtDayStart}
                endAtDayEnd={endAtDayEnd}
                onStartDateChange={setStartDate}
                onEndDateChange={setEndDate}
                onStartAtDayStartChange={setStartAtDayStart}
                onEndAtDayEndChange={setEndAtDayEnd}
                onSave={action('onSave')}
            />
        </LocalizationProvider>
    );
};

export const Default = (args: Props) => <EditDatesDialogWrapper {...args} />
Default.args = {
    albumName: "First Contact",
    startDate: new Date("2063-04-05T00:00:00Z"),
    endDate: new Date("2063-04-05T23:59:00Z"),
    startAtDayStart: true,
    endAtDayEnd: true,
    isLoading: false,
    isSaveEnabled: true,
};

export const WithSpecificTimes = (args: Props) => <EditDatesDialogWrapper {...args} />
WithSpecificTimes.args = {
    albumName: "Mission to Mars",
    startDate: new Date("2063-04-05T10:30:00Z"),
    endDate: new Date("2063-04-07T15:45:00Z"),
    startAtDayStart: false,
    endAtDayEnd: false,
    isLoading: false,
    isSaveEnabled: true,
};

export const Loading = (args: Props) => <EditDatesDialogWrapper {...args} />
Loading.args = {
    albumName: "First Contact",
    startDate: new Date("2063-04-05T00:00:00Z"),
    endDate: new Date("2063-04-05T23:59:00Z"),
    startAtDayStart: true,
    endAtDayEnd: true,
    isLoading: true,
    isSaveEnabled: false,
};
Loading.meta = {skipSnapshot: true}

export const WithError = (args: Props) => <EditDatesDialogWrapper {...args} />
WithError.args = {
    albumName: "First Contact",
    startDate: new Date("2063-04-05T00:00:00Z"),
    endDate: new Date("2063-04-05T23:59:00Z"),
    startAtDayStart: true,
    endAtDayEnd: true,
    isLoading: false,
    errorCode: "This is a user friendly error message (or technical).",
    isSaveEnabled: true,
};

export const WithDateRangeError = (args: Props) => <EditDatesDialogWrapper {...args} />
WithDateRangeError.args = {
    albumName: "First Contact",
    startDate: new Date("2063-04-07T00:00:00Z"),
    endDate: new Date("2063-04-05T00:00:00Z"),
    startAtDayStart: true,
    endAtDayEnd: true,
    isLoading: false,
    dateRangeError: "The end date cannot be before the start date",
    isSaveEnabled: false,
};
