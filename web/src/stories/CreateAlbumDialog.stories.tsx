import React from 'react';
import {ComponentMeta, ComponentStory} from '@storybook/react';
import CreateAlbumDialog from "../pages/authenticated/albums/CreateAlbumDialog";
import {AdapterDayjs} from "@mui/x-date-pickers/AdapterDayjs";
import {LocalizationProvider} from '@mui/x-date-pickers';
import dayjs from "dayjs";
import fr from "dayjs/locale/fr";
import {screen, userEvent, within} from "@storybook/testing-library";

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

export const WithAName = Template.bind({});
WithAName.args = {
    open: true,
};
WithAName.parameters = {
    delay: 300,
};
WithAName.play = async ({canvasElement}) => {
    const canvas = within(screen.getByRole('dialog'));
    const nameInput = canvas.getByLabelText(/Name/, {
        selector: 'input',
    })

    userEvent.type(nameInput, 'Avenger 3');
};

export const AlreadyExists = Template.bind({});
AlreadyExists.args = {
    open: true,
};
AlreadyExists.parameters = {
    delay: 300,
};
// AlreadyExists.play = async ({canvasElement}) => {
//     const canvas = within(screen.getByRole('dialog'));
//     const nameInput = canvas.getByLabelText(/Name/, {
//         selector: 'input',
//     })
//
//     userEvent.type(nameInput, 'Avenger 3');
// };

