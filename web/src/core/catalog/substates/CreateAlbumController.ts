import dayjs, {Dayjs} from "dayjs";
import {isCatalogError} from "../domain";
import {CreateAlbumThunk} from "../thunks";

export const albumFolderNameAlreadyTakenErr = "AlbumFolderNameAlreadyTakenErr";
export const albumStartAndEndDateMandatoryErr = "AlbumStartAndEndDateMandatoryErr";

export interface CreateAlbumState {
    open: boolean
    creationInProgress: boolean
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
    onCloseCreateAlbumDialog: () => void
    onSubmitCreateAlbum: (state: CreateAlbumState) => void
    onNameChange: (value: string) => void
    onFolderNameChange: (value: string) => void
    onWithCustomFolderNameChange: (withCustomFolderName: boolean) => void
    onStartsAtStartOfTheDayChange: (startsAtStartOfTheDay: boolean) => void
    onEndsAtEndOfTheDayChange: (endsAtEndOfTheDay: boolean) => void
    onStartDateChange: (start: Dayjs | null) => void
    onEndDateChange: (end: Dayjs | null) => void
}

export interface CreateAlbumControls {
    openDialogForCreateAlbum: () => void
    // openEdit: (album: Album) => void // TODO manage EDIT mode
}

export class CreateAlbumController implements CreateAlbumHandlers, CreateAlbumControls {
    constructor(
        private readonly setState: (stateUpdater: (prev: CreateAlbumState) => CreateAlbumState) => void,
        private readonly createAlbumPort: CreateAlbumThunk,
        private readonly firstDay: Dayjs = dayjs().startOf("week").subtract(9, "days"),
    ) {
    }

    openDialogForCreateAlbum = (): void => {
        this.setState(prev => ({
            ...emptyCreateAlbum(this.firstDay),
            open: true
        }));
    }

    onCloseCreateAlbumDialog = (): void => {
        this.setState(prev => ({...prev, open: false}));
    }

    onSubmitCreateAlbum = async (state: CreateAlbumState): Promise<void> => {
        const actualStart = state.startsAtStartOfTheDay ? state.start?.startOf("day") : state.start;
        const actualEnd = state.endsAtEndOfTheDay ? state.end?.add(1, "day").startOf("day") : state.end;

        if (actualStart && actualEnd && actualStart.isBefore(actualEnd)) {
            this.setState(prev => ({...prev, creationInProgress: true}));

            await this.createAlbumPort({
                name: state.name,
                start: actualStart.toDate(),
                end: actualEnd.toDate(),
                forcedFolderName: state.withCustomFolderName ? state.forceFolderName : ""
            })
                .then((albumId) => {
                    this.setState(prev => ({...prev, open: false}))
                    return albumId
                })
                .catch((err: Error) => {
                    console.log("Failed to create the album", err);
                    this.setState(prev => ({
                        ...prev,
                        errorCode: isCatalogError(err) ? err.errorCode : err.message,
                        creationInProgress: false,
                    }))
                    return Promise.reject(err);
                })
        } else {
            this.setState(prev => ({
                ...prev,
                creationInProgress: false,
                errorCode: albumStartAndEndDateMandatoryErr,
            }));
        }
    }

    onWithCustomFolderNameChange = (withCustomFolderName: boolean): void => {
        this.setState(prev => ({...prev, withCustomFolderName}));
    }

    onNameChange = (value: string): void => {
        this.setState(prev => ({
            ...prev,
            name: value,
            errorCode: prev.errorCode !== albumFolderNameAlreadyTakenErr ? prev.errorCode : undefined
        }));
    }

    onFolderNameChange = (value: string): void => {
        this.setState(prev => ({
            ...prev,
            forceFolderName: value,
            errorCode: prev.errorCode !== albumFolderNameAlreadyTakenErr ? prev.errorCode : undefined
        }));
    }

    onStartsAtStartOfTheDayChange = (startsAtStartOfTheDay: boolean): void => {
        this.setState(prev => ({...prev, startsAtStartOfTheDay}));
    }

    onEndsAtEndOfTheDayChange = (endsAtEndOfTheDay: boolean): void => {
        this.setState(prev => ({...prev, endsAtEndOfTheDay}));
    }

    onStartDateChange = (start: Dayjs | null): void => {
        this.setState(prev => ({
            ...prev,
            start,
            errorCode: prev.errorCode === albumStartAndEndDateMandatoryErr ? undefined : prev.errorCode
        }));
    }

    onEndDateChange = (end: Dayjs | null): void => {
        this.setState(prev => ({
            ...prev,
            end,
            errorCode: prev.errorCode === albumStartAndEndDateMandatoryErr ? undefined : prev.errorCode
        }));
    }
}

export const emptyCreateAlbum = (defaultDate: Dayjs): CreateAlbumState => ({
    open: false,
    creationInProgress: false,
    name: "",
    start: defaultDate,
    end: defaultDate.add(8, "days").endOf("day"),
    forceFolderName: "",
    startsAtStartOfTheDay: true,
    endsAtEndOfTheDay: true,
    withCustomFolderName: false,
});