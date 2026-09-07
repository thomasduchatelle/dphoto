'use client';

import {AlbumGrid} from "../AlbumGrid";
import {catalogReducer, catalogThunks, CatalogViewerState} from "@/domains/catalog";
import {useReducer} from "react";
import {useThunks} from "@/libs/dthunks/react";
import {ErrorMessage} from "@/components/ErrorMessage";

export const HomeContent = ({initialState}: { initialState: CatalogViewerState }) => {
    const [state, dispatch] = useReducer(catalogReducer, initialState);

    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const {onPageRefresh, loadAlbumPage, deleteAlbum, updateAlbumDates, submitCreateAlbum, saveAlbumName, grantAlbumAccess, revokeAlbumAccess, ...dispatchOnlyThunks} = catalogThunks;
    const thunks = useThunks(dispatchOnlyThunks, {dispatch}, state);

    if (state.error) {
        return (
            <ErrorMessage error={state.error} title="Failed to load the albums"/>
        )
    }

    return (
        <AlbumGrid
            albums={state.albums}
            filterOptions={state.albumFilterOptions}
            activeFilter={state.albumFilter}
            onFilterChange={thunks.onAlbumFilterChange}
            onShare={thunks.openSharingModal}
        />
    );
}