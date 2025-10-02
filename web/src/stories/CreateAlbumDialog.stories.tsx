import React from 'react';
import {AdapterDayjs} from "@mui/x-date-pickers/AdapterDayjs";
import {LocalizationProvider} from '@mui/x-date-pickers';
import dayjs from "dayjs";
import fr from "dayjs/locale/fr";
import {CreateAlbumDialog} from "../pages/authenticated/albums/CreateAlbumDialog";
import {Story} from "@ladle/react";
import {CreateDialogSelection} from "../core/catalog";

dayjs.locale(fr)

const defaultStartDate = new Date("2024-12-21T00:00:00Z")
const endDate = new Date("2024-12-29T23:59:00Z")

export default {
    title: 'Albums/CreateAlbumDialog',
};

type Props = Exclude<CreateDialogSelection, "open">

export const Default: Story<Props> = (props) => (
    <LocalizationProvider dateAdapter={AdapterDayjs} adapterLocale='fr'>
        <CreateAlbumDialog onClose={function (): void {
            throw new Error("Function not implemented.");
        }} onSubmit={function (): Promise<void> {
            throw new Error("Function not implemented.");
        }} onNameChange={function (name: string): void {
            throw new Error("Function not implemented.");
        }} onFolderNameChange={function (folderName: string): void {
            throw new Error("Function not implemented.");
        }} onWithCustomFolderNameChange={function (withCustom: boolean): void {
            throw new Error("Function not implemented.");
        }} onStartsAtStartOfTheDayChange={function (startsAtStart: boolean): void {
            throw new Error("Function not implemented.");
        }} onEndsAtEndOfTheDayChange={function (endsAtEnd: boolean): void {
            throw new Error("Function not implemented.");
        }} onStartDateChange={function (date: Date | null): void {
            throw new Error("Function not implemented.");
        }} onEndDateChange={function (date: Date | null): void {
            throw new Error("Function not implemented.");
        }} {...props}/>
    </LocalizationProvider>
);

Default.args = {
    label: 'Hello world',
    disabled: false,
    count: 2,
    colors: ['Red', 'Blue'],
};


export const Default2 = () => {
    const args = {
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

    return (
        <LocalizationProvider dateAdapter={AdapterDayjs} adapterLocale='fr'>
            <CreateAlbumDialog {...args}/>
        </LocalizationProvider>
    )
}
// InFrame.meta = {iframed: true}
