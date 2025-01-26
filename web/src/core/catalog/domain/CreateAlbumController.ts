import dayjs, {Dayjs} from "dayjs";
import React from "react";
import {isCatalogError} from "./errors";
import {OnCreateNewAlbumRequestType} from "../../catalog-react";
import {albumFolderNameAlreadyTakenErr, albumStartAndEndDateMandatoryErr} from "../../../pages/authenticated/albums/CreateAlbumDialog/CreateAlbumDialog";

export interface CreateAlbumState {
    open: boolean
    name: string
    start: Dayjs | null
    end: Dayjs | null
    forceFolderName: string
    startsAtStartOfTheDay: boolean
    endsAtEndOfTheDay: boolean
    withCustomFolderName: boolean
    errorCode?: string
}

export interface CreateAlbumHandlers {
    onClose: () => void
    onSubmit: (state: CreateAlbumState) => void
    setStartsAtStartOfTheDay: (startsAtStartOfTheDay: boolean) => void
    setEndsAtEndOfTheDay: (endsAtEndOfTheDay: boolean) => void
    setWithCustomFolderName: (withCustomFolderName: boolean) => void
    handleOnNameChange: (event: React.ChangeEvent<HTMLInputElement>) => void
    handleStartDateChange: (start: Dayjs | null) => void
    handleEndDateChange: (end: Dayjs | null) => void
    handleOnFolderNameChange: (event: React.ChangeEvent<HTMLInputElement>) => void
}

export interface CreateAlbumControls {
    openNew: () => void
    // openEdit: (album: Album) => void // TODO manage EDIT mode
}

export interface CreateAlbumPort {
    createAlbum: OnCreateNewAlbumRequestType
}

export class CreateAlbumController implements CreateAlbumHandlers, CreateAlbumControls {
    constructor(
        private readonly setState: (stateUpdater: (prev: CreateAlbumState) => CreateAlbumState) => void,
        private readonly createAlbumPort: CreateAlbumPort,
        private readonly firstDay: Dayjs = dayjs(),
    ) {
    }

    openNew = (): void => {
        this.setState(prev => emptyCreateAlbum(this.firstDay));
    }

    onClose = (): void => {
        this.setState(prev => ({...prev, open: false}));
    }

    onSubmit = (state: CreateAlbumState): void => {
        if (state.start && state.end) {
            this.createAlbumPort.createAlbum({
                name: state.name,
                start: state.startsAtStartOfTheDay ? state.start.startOf("day").toDate() : state.start.toDate(),
                end: state.endsAtEndOfTheDay ? state.end.endOf("day").toDate() : state.end.toDate(),
                forcedFolderName: state.withCustomFolderName ? state.forceFolderName : ""
            })
                .catch(err => {
                    console.log("Failed to create the album", err);
                    this.setState(prev => ({
                        ...prev,
                        errorCode: isCatalogError(err) ? err.errorCode : err,
                    }));
                });
        } else {
            this.setState(prev => ({
                ...prev,
                errorCode: albumStartAndEndDateMandatoryErr,
            }));
        }
    }

    setWithCustomFolderName = (withCustomFolderName: boolean): void => {
        this.setState(prev => ({...prev, withCustomFolderName}));
    }

    handleOnNameChange = (event: React.ChangeEvent<HTMLInputElement>): void => {
        this.setState(prev => ({
            ...prev,
            name: event.target.value,
            errorCode: prev.errorCode !== albumFolderNameAlreadyTakenErr ? prev.errorCode : undefined
        }));
    }

    handleOnFolderNameChange = (event: React.ChangeEvent<HTMLInputElement>): void => {
        this.setState(prev => ({
            ...prev,
            forceFolderName: event.target.value,
            errorCode: prev.errorCode !== albumFolderNameAlreadyTakenErr ? prev.errorCode : undefined
        }));
    }

    setStartsAtStartOfTheDay = (startsAtStartOfTheDay: boolean): void => {
        this.setState(prev => ({...prev, startsAtStartOfTheDay}));
    }

    setEndsAtEndOfTheDay = (endsAtEndOfTheDay: boolean): void => {
        this.setState(prev => ({...prev, endsAtEndOfTheDay}));
    }

    handleStartDateChange = (start: Dayjs | null): void => {
        this.setState(prev => ({
            ...prev,
            start,
            errorCode: prev.errorCode === albumStartAndEndDateMandatoryErr ? undefined : prev.errorCode
        }));
    }

    handleEndDateChange = (end: Dayjs | null): void => {
        this.setState(prev => ({
            ...prev,
            end,
            errorCode: prev.errorCode === albumStartAndEndDateMandatoryErr ? undefined : prev.errorCode
        }));
    }
}

export const emptyCreateAlbum = (defaultDate: Dayjs): CreateAlbumState => ({
    open: false,
    name: "",
    start: defaultDate,
    end: defaultDate.add(7, "days").endOf("day"),
    forceFolderName: "",
    startsAtStartOfTheDay: true,
    endsAtEndOfTheDay: true,
    withCustomFolderName: false,
});