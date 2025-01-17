import React from 'react';
import {ComponentMeta, ComponentStory} from '@storybook/react';
import CreateAlbumDialog from "../pages/authenticated/albums/CreateAlbumDialog";
import {AdapterDayjs} from "@mui/x-date-pickers/AdapterDayjs";
import {LocalizationProvider} from '@mui/x-date-pickers';
import dayjs from "dayjs";
import fr from "dayjs/locale/fr";

dayjs.locale(fr)

// More on default export: https://storybook.js.org/docs/react/writing-stories/introduction#default-export
export default {
    title: 'Albums/CreateAlbumDialog',
    component: CreateAlbumDialog,
} as ComponentMeta<typeof CreateAlbumDialog>;

// overridesDefaultStartDate(dayjs("2025-01-01"))

// More on component templates: https://storybook.js.org/docs/react/writing-stories/introduction#using-args
const Template: ComponentStory<typeof CreateAlbumDialog> = (args) => (
    <LocalizationProvider dateAdapter={AdapterDayjs} adapterLocale='fr'>
        <CreateAlbumDialog defaultDate={dayjs("2024-12-21")} {...args}/>
    </LocalizationProvider>
);

export const Empty = Template.bind({});
Empty.args = {
    open: true,
};
Empty.parameters = {
    delay: 300,
};

