'use client';

import {AlbumGrid} from "../AlbumGrid";
import {CreateAlbumModal} from "../CreateAlbumModal";
import {catalogReducer, catalogThunks, CatalogViewerState} from "@/domains/catalog";
import {createDialogSelector} from "@/domains/catalog/album-create";
import {newClientAdapterFactory} from "@/domains/catalog/adapters/client-adapter-factory";
import {useMemo, useReducer} from "react";
import {useThunks} from "@/libs/dthunks/react";
import {ErrorMessage} from "@/components/ErrorMessage";

export const HomeContent = ({initialState}: { initialState: CatalogViewerState }) => {
    const [state, dispatch] = useReducer(catalogReducer, initialState);

    const adapter = useMemo(() => newClientAdapterFactory(), []);

    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const {onPageRefresh, loadAlbumPage, deleteAlbum, updateAlbumDates, saveAlbumName, grantAlbumAccess, revokeAlbumAccess, ...clientThunks} = catalogThunks;
    const thunks = useThunks(clientThunks, {dispatch, adapter}, state);

    const createDialog = createDialogSelector(state);
    const canCreateAlbum = state.currentUser.isOwner;

    if (state.error) {
        return (
            <ErrorMessage error={state.error} title="Failed to load the albums"/>
        )
    }

    return (
        <>
            <AlbumGrid
                albums={state.albums}
                filterOptions={state.albumFilterOptions}
                activeFilter={state.albumFilter}
                onFilterChange={thunks.onAlbumFilterChange}
                onShare={thunks.openSharingModal}
                onCreateAlbum={canCreateAlbum ? thunks.openCreateDialog : undefined}
            />
            <CreateAlbumModal
                {...createDialog}
                onClose={thunks.closeCreateDialog}
                onSubmit={thunks.submitCreateAlbum}
                onNameChange={thunks.changeAlbumName}
                onFolderNameChange={thunks.changeFolderName}
                onWithCustomFolderNameChange={thunks.changeFolderNameEnabled}
                onStartsAtStartOfTheDayChange={thunks.updateCreateDialogStartsAtStartOfTheDay}
                onEndsAtEndOfTheDayChange={thunks.updateCreateDialogEndsAtEndOfTheDay}
                onStartDateChange={thunks.updateCreateDialogStartDate}
                onEndDateChange={thunks.updateCreateDialogEndDate}
            />
        </>
    )
}
