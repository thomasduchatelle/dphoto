import React from 'react';
import {ComponentMeta, ComponentStory} from '@storybook/react';
import {AdapterDayjs} from "@mui/x-date-pickers/AdapterDayjs";
import {LocalizationProvider} from '@mui/x-date-pickers';
import dayjs from "dayjs";
import fr from "dayjs/locale/fr";
import {screen, userEvent, within} from "@storybook/testing-library";
import {CreateAlbumDialog, CreateAlbumDialogContainer} from "../pages/authenticated/albums/CreateAlbumDialog";
import {albumFolderNameAlreadyTakenErr, albumStartAndEndDateMandatoryErr, emptyCreateAlbum} from "../core/catalog";
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
const TemplateInContainer: ComponentStory<typeof CreateAlbumDialog> = (args) => (
    <LocalizationProvider dateAdapter={AdapterDayjs} adapterLocale='fr'>
        <CreateAlbumDialogContainer firstDay={defaultStartDate} createAlbum={() => Promise.resolve({owner: "tony", folderName: "/ironman-1"})} {...args}>
            {({openDialogForCreateAlbum}) => {
                return <Button onClick={openDialogForCreateAlbum} variant='contained'>Click to open</Button>
            }}
        </CreateAlbumDialogContainer>
    </LocalizationProvider>
);

// it should open the dialog and set a name when used with the container
export const InContainer = TemplateInContainer.bind({});
InContainer.args = {};
InContainer.parameters = {
    delay: 300,
};
InContainer.play = async ({canvasElement}) => {
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
        name: 'Avenger 3',
    }
};
WithAName.parameters = {
    delay: 300,
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

// it should render an error on the Name field when the error albumFolderNameAlreadyTakenErr is raised ; and the error should clear when the name of folder name are updated
export const Loading = Template.bind({});
Loading.args = {
    state: {
        ...emptyCreateAlbum(defaultStartDate),
        name: 'Avenger 3',
        creationInProgress: true,
        open: true,
    },
};
Loading.parameters = {
    storyshots: {disable: true},
};

// it should render an error on the Name field when the error albumFolderNameAlreadyTakenErr is raised ; and the error should clear when the name of folder name are updated
export const GenericError = Template.bind({});
GenericError.args = {
    state: {
        ...emptyCreateAlbum(defaultStartDate),
        name: 'Avenger 3',
        open: true,
        errorCode: 'Something weird and different than the known errors.'
    },
};
GenericError.parameters = {
    delay: 300,
};

