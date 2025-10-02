import React from 'react';
import {ComponentMeta, ComponentStory} from '@storybook/react';
import {AdapterDayjs} from "@mui/x-date-pickers/AdapterDayjs";
import {LocalizationProvider} from '@mui/x-date-pickers';
import dayjs from "dayjs";
import fr from "dayjs/locale/fr";
import {CreateAlbumDialog} from "../pages/authenticated/albums/CreateAlbumDialog";

dayjs.locale(fr)

// More on default export: https://storybook.js.org/docs/react/writing-stories/introduction#default-export
export default {
    title: 'Albums/CreateAlbumDialogOld',
    component: CreateAlbumDialog,
} as ComponentMeta<typeof CreateAlbumDialog>;

const defaultStartDate = new Date("2024-12-21T00:00:00Z")
const endDate = new Date("2024-12-29T23:59:00Z")

// More on component templates: https://storybook.js.org/docs/react/writing-stories/introduction#using-args
const Template: ComponentStory<typeof CreateAlbumDialog> = (args) => (
    <LocalizationProvider dateAdapter={AdapterDayjs} adapterLocale='fr'>
        <CreateAlbumDialog {...args}/>
    </LocalizationProvider>
);

// it should display the model with no name, defaulted start and end date (1 week apart), no folder name, "create" button disabled
export const Empty = Template.bind({});
Empty.args = {
    open: true,
    albumName: "",
    start: defaultStartDate,
    end: endDate,
    customFolderName: "",
    startsAtStartOfTheDay: true,
    endsAtEndOfTheDay: true,
    isCustomFolderNameEnabled: false,
    isLoading: false,
    canSubmit: false,
    onClose: () => {
    },
    onSubmit: () => Promise.resolve(),
    onNameChange: () => {
    },
    onStartDateChange: () => {
    },
    onEndDateChange: () => {
    },
    onFolderNameChange: () => {
    },
    onWithCustomFolderNameChange: () => {
    },
    onStartsAtStartOfTheDayChange: () => {
    },
    onEndsAtEndOfTheDayChange: () => {
    },
};
Empty.parameters = {
    delay: 300,
};

// it should have the "save" button enabled when the name is not empty
export const WithAName = Template.bind({});
WithAName.args = {
    ...Empty.args,
    albumName: 'Avenger 3',
    canSubmit: true,
};
WithAName.parameters = {
    delay: 300,
};

// it should render an error on the Name field when the error albumFolderNameAlreadyTakenErr is raised ; and the error should clear when the name of folder name are updated
export const NameAlreadyExists = Template.bind({});
NameAlreadyExists.args = {
    ...Empty.args,
    nameError: "Name must be unique",
};
NameAlreadyExists.parameters = {
    delay: 300,
};

// it should render an error on the Name field when the error albumFolderNameAlreadyTakenErr is raised ; and the error should clear when the name of folder name are updated
export const FolderNameAlreadyExists = Template.bind({});
FolderNameAlreadyExists.args = {
    ...Empty.args,
    folderNameError: "Name must be unique",
    customFolderName: "/avenger-3",
    isCustomFolderNameEnabled: true,
};
FolderNameAlreadyExists.parameters = {
    delay: 300,
};

// it should render an error on the Name field when the error albumFolderNameAlreadyTakenErr is raised ; and the error should clear when the name of folder name are updated
export const StartAndEndDateAreMandatory = Template.bind({});
StartAndEndDateAreMandatory.args = {
    ...Empty.args,
    dateRangeError: "Start date must be before end date",
};
StartAndEndDateAreMandatory.parameters = {
    delay: 300,
};

// it should render an error on the Name field when the error albumFolderNameAlreadyTakenErr is raised ; and the error should clear when the name of folder name are updated
export const Loading = Template.bind({});
Loading.args = {
    ...WithAName.args,
    isLoading: true,
    canSubmit: false,
};
Loading.parameters = {
    storyshots: {disable: true},
};

// it should render an error on the Name field when the error albumFolderNameAlreadyTakenErr is raised ; and the error should clear when the name of folder name are updated
export const GenericError = Template.bind({});
GenericError.args = {
    ...WithAName.args,
    error: 'Something weird and different than the known errors.'
};
GenericError.parameters = {
    delay: 300,
};

