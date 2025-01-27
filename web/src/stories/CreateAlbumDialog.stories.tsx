import React from 'react';
import {ComponentMeta, ComponentStory} from '@storybook/react';
import {AdapterDayjs} from "@mui/x-date-pickers/AdapterDayjs";
import {LocalizationProvider} from '@mui/x-date-pickers';
import dayjs from "dayjs";
import fr from "dayjs/locale/fr";
import {screen, userEvent, within} from "@storybook/testing-library";
import {
    albumFolderNameAlreadyTakenErr,
    albumStartAndEndDateMandatoryErr,
    CreateAlbumDialog,
    CreateAlbumDialogContainer
} from "../pages/authenticated/albums/CreateAlbumDialog";
import {emptyCreateAlbum} from "../core/catalog/domain/CreateAlbumController";
import {Button} from "@mui/material";

dayjs.locale(fr)

// More on default export: https://storybook.js.org/docs/react/writing-stories/introduction#default-export
export default {
    title: 'Albums/CreateAlbumDialog',
    component: CreateAlbumDialog,
} as ComponentMeta<typeof CreateAlbumDialog>;

const defaultStartDate = dayjs("2024-12-21")

// More on component templates: https://storybook.js.org/docs/react/writing-stories/introduction#using-args
const Template: ComponentStory<typeof CreateAlbumDialog> = (args) => (
    <LocalizationProvider dateAdapter={AdapterDayjs} adapterLocale='fr'>
        <CreateAlbumDialog {...args}/>
    </LocalizationProvider>
);
const TemplateContainer: ComponentStory<typeof CreateAlbumDialog> = (args) => (
    <LocalizationProvider dateAdapter={AdapterDayjs} adapterLocale='fr'>
        <CreateAlbumDialogContainer firstDay={defaultStartDate}>
            {({openNew}) => {
                return <Button onClick={openNew} variant='contained'>Click to open</Button>
            }}
        </CreateAlbumDialogContainer>
    </LocalizationProvider>
);

// it should open the dialog and set a name when used with the container
export const WithContainer = TemplateContainer.bind({});
WithContainer.args = {};
WithContainer.parameters = {
    delay: 300,
};
WithContainer.play = async ({canvasElement}) => {
    const fullCanvas = within(canvasElement);
    await userEvent.click(fullCanvas.getAllByRole("button")[0]);

    const canvas = within(screen.getByRole('dialog'));
    const nameInput = canvas.getByLabelText(/Name/, {
        selector: 'input',
    })

    userEvent.type(nameInput, 'Avenger 3');
};

// it should display the model with no name, defaulted start and end date (1 week apart), no folder name, "create" button disabled
export const Empty = Template.bind({});
Empty.args = {
    state: {
        ...emptyCreateAlbum(defaultStartDate),
        open: true,
    }
};
Empty.parameters = {
    delay: 300,
};

// it should have the "save" button enabled when the name is not empty
export const WithAName = Template.bind({});
WithAName.args = {
    state: {
        ...emptyCreateAlbum(defaultStartDate),
        open: true,
    }
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

// it should render an error on the Name field when the error albumFolderNameAlreadyTakenErr is raised ; and the error should clear when the name of folder name are updated
export const AlreadyExists = Template.bind({});
AlreadyExists.args = {
    state: {
        ...emptyCreateAlbum(defaultStartDate),
        open: true,
        errorCode: albumFolderNameAlreadyTakenErr
    },
};
AlreadyExists.parameters = {
    delay: 300,
};

// it should render an error on the Name field when the error albumFolderNameAlreadyTakenErr is raised ; and the error should clear when the name of folder name are updated
export const StartAndEndDateAreMandatory = Template.bind({});
StartAndEndDateAreMandatory.args = {
    state: {
        ...emptyCreateAlbum(defaultStartDate),
        open: true,
        errorCode: albumStartAndEndDateMandatoryErr
    },
};
StartAndEndDateAreMandatory.parameters = {
    delay: 300,
};

