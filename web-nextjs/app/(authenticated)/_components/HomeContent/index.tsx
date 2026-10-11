'use client';

import {AlbumGrid} from "../AlbumGrid";
import {CreateAlbumModal} from "../CreateAlbumModal";
import {DeleteAlbumModal} from "@/components/DeleteAlbumModal";
import {catalogReducer, catalogThunks, CatalogViewerState} from "@/domains/catalog";
import {createDialogSelector} from "@/domains/catalog/album-create";
import {deleteDialogSelector} from "@/domains/catalog/album-delete/selector-deleteDialogSelector";
import {newClientAdapterFactory} from "@/domains/catalog/adapters/client-adapter-factory";
import {useMemo, useReducer} from "react";
import {useSessionAwareThunks} from "@/libs/security/client";
import {ErrorMessage} from "@/components/ErrorMessage";

export const HomeContent = ({initialState}: { initialState: CatalogViewerState }) => {
    const [state, dispatch] = useReducer(catalogReducer, initialState);

    const adapter = useMemo(() => newClientAdapterFactory(), []);

    const thunks = useSessionAwareThunks(catalogThunks, {dispatch, adapter}, state);

    const createDialog = createDialogSelector(state);
    const deleteDialog = deleteDialogSelector(state);
    const canCreateAlbum = state.currentUser.isOwner;

    if (state.error) {
        return (
            <ErrorMessage error={state.error} title="Failed to load the albums"/>
        );
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
                onDeleteAlbum={canCreateAlbum ? thunks.openDeleteAlbumDialog : undefined}
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
            <DeleteAlbumModal
                {...deleteDialog}
                onClose={thunks.closeDeleteAlbumDialog}
                onDelete={thunks.deleteAlbum}
            />
        </>
    );
}
