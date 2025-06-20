import React from "react";
import {ComponentMeta, ComponentStory} from "@storybook/react";
import {EditDatesDialog} from "../pages/authenticated/albums/EditDatesDialog";
import {Button} from "@mui/material";
import {AdapterDayjs} from "@mui/x-date-pickers/AdapterDayjs";
import {LocalizationProvider} from "@mui/x-date-pickers";

export default {
    title: "Albums/EditDatesDialog",
    component: EditDatesDialog,
} as ComponentMeta<typeof EditDatesDialog>;

const Template: ComponentStory<typeof EditDatesDialog> = (args) => {
    const [open, setOpen] = React.useState(true);
    const [startDate, setStartDate] = React.useState<Date | null>(args.startDate);
    const [endDate, setEndDate] = React.useState<Date | null>(args.endDate);
    const [startAtDayStart, setStartAtDayStart] = React.useState(args.startAtDayStart);
    const [endAtDayEnd, setEndAtDayEnd] = React.useState(args.endAtDayEnd);

    React.useEffect(() => {
        setStartDate(args.startDate);
        setEndDate(args.endDate);
        setStartAtDayStart(args.startAtDayStart);
        setEndAtDayEnd(args.endAtDayEnd);
    }, [args.startDate, args.endDate, args.startAtDayStart, args.endAtDayEnd]);

    return (
        <LocalizationProvider dateAdapter={AdapterDayjs} adapterLocale='fr'>
            <Button onClick={() => setOpen(true)} variant="contained">
                Open Edit Dates Dialog
            </Button>
            <EditDatesDialog
                {...args}
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
            />
        </LocalizationProvider>
    );
};

export const Default = Template.bind({});
Default.args = {
    albumName: "First Contact",
    startDate: new Date("2063-04-05T00:00:00Z"),
    endDate: new Date("2063-04-05T23:59:00Z"),
    startAtDayStart: true,
    endAtDayEnd: true,
};
Default.parameters = {
    delay: 300,
};
