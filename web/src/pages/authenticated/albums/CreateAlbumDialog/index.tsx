import React, {useCallback, useState} from 'react';
import dayjs, {Dayjs} from 'dayjs';
import CreateAlbumDialog, {albumFolderNameAlreadyTakenErr, albumStartAndEndDateMandatoryErr} from './CreateAlbumDialogPure';
import {OnCreateNewAlbumRequestType} from '../../../../core/catalog-react';
import {isCatalogError} from '../../../../core/catalog/domain/errors';
import {CreateAlbumState} from "../../../../core/catalog/domain/CreateAlbumController";

const saturdayTwoWeeksAgo = dayjs().startOf("week").subtract(8, "days");

const emptyCreateAlbum = (defaultDate: Dayjs): CreateAlbumState => ({
    open: false,
    name: "",
    start: defaultDate,
    end: defaultDate.add(7, "days").endOf("day"),
    forceFolderName: "",
    startsAtStartOfTheDay: true,
    endsAtEndOfTheDay: true,
    withCustomFolderName: false,
});

export default function CreateAlbumDialogContainer({open, onClose, onSubmit, defaultDate = saturdayTwoWeeksAgo, defaultErrorCode}: {
    open: boolean,
    onClose: () => void,
    onSubmit: OnCreateNewAlbumRequestType,
    defaultDate?: Dayjs,
    defaultErrorCode?: string
}) {
    const [state, setState] = useState<CreateAlbumState>({...emptyCreateAlbum(defaultDate), errorCode: defaultErrorCode});
    const setStartsAtStartOfTheDay = (startsAtStartOfTheDay: boolean) => setState(prev => ({...prev, startsAtStartOfTheDay}));
    const setEndsAtEndOfTheDay = (endsAtEndOfTheDay: boolean) => setState(prev => ({...prev, endsAtEndOfTheDay}));
    const setWithCustomFolderName = (withCustomFolderName: boolean) => setState(prev => ({...prev, withCustomFolderName}));

    const handleOnNameChange = useCallback((event: React.ChangeEvent<HTMLInputElement>) => {
        setState(prev => ({...prev, name: event.target.value, errorCode: prev.errorCode !== albumFolderNameAlreadyTakenErr ? prev.errorCode : undefined}));
    }, [setState]);
    const handleStartDateChange = useCallback((start: Dayjs | null) => setState(prev => ({
        ...prev,
        start,
        errorCode: prev.errorCode === albumStartAndEndDateMandatoryErr ? undefined : prev.errorCode
    })), [setState]);
    const handleEndDateChange = useCallback((end: Dayjs | null) => setState(prev => ({
        ...prev,
        end,
        errorCode: prev.errorCode === albumStartAndEndDateMandatoryErr ? undefined : prev.errorCode
    })), [setState]);
    const handleOnFolderNameChange = useCallback((event: React.ChangeEvent<HTMLInputElement>) => {
        setState(prev => ({
            ...prev,
            forceFolderName: event.target.value,
            errorCode: prev.errorCode !== albumFolderNameAlreadyTakenErr ? prev.errorCode : undefined
        }));
    }, [setState]);

    const handleClose = useCallback(() => {
        onClose();
        setState(emptyCreateAlbum(defaultDate));
    }, [setState, onClose, defaultDate]);

    const handleSubmit = useCallback(() => {
        if (state.start && state.end) {
            onSubmit({
                name: state.name,
                start: state.startsAtStartOfTheDay ? state.start.startOf("day").toDate() : state.start.toDate(),
                end: state.endsAtEndOfTheDay ? state.end.endOf("day").toDate() : state.end.toDate(),
                forcedFolderName: state.withCustomFolderName ? state.forceFolderName : ""
            })
                .then(() => setState(emptyCreateAlbum(defaultDate)))
                .catch(err => {
                    console.log("Failed to create the album", err);
                    setState(prev => ({
                        ...prev,
                        errorCode: isCatalogError(err) ? err.errorCode : err,
                    }));
                });
        } else {
            setState(prev => ({
                ...prev,
                errorCode: albumStartAndEndDateMandatoryErr,
            }));
        }
    }, [state, setState, onSubmit, defaultDate]);

    return (
        <CreateAlbumDialog
            open={open}
            onClose={handleClose}
            onSubmit={handleSubmit}
            state={state}
            setStartsAtStartOfTheDay={setStartsAtStartOfTheDay}
            setEndsAtEndOfTheDay={setEndsAtEndOfTheDay}
            setWithCustomFolderName={setWithCustomFolderName}
            handleOnNameChange={handleOnNameChange}
            handleStartDateChange={handleStartDateChange}
            handleEndDateChange={handleEndDateChange}
            handleOnFolderNameChange={handleOnFolderNameChange}
        />
    );
}