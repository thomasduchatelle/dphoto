'use client';

import {AlbumGrid} from "../AlbumGrid";
import {Album, AlbumFilterCriterion, AlbumFilterEntry, AlbumId, catalogReducer, catalogThunks, CatalogViewerState} from "@/domains/catalog";
import {useReducer} from "react";
import {useSessionAwareThunks} from "@/libs/security/client";
import {ErrorMessage} from "@/components/ErrorMessage";

export const HomeContent = ({initialState}: { initialState: CatalogViewerState }) => {
    const [state, dispatch] = useReducer(catalogReducer, initialState);

    const {onPageRefresh, loadAlbumPage, deleteAlbum, updateAlbumDates, submitCreateAlbum, saveAlbumName, grantAlbumAccess, revokeAlbumAccess, ...dispatchOnlyThunks} = catalogThunks;
    const thunks = useSessionAwareThunks(dispatchOnlyThunks, {dispatch}, state);

    return (
        <HomeContentView
            error={state.error}
            albums={state.albums}
            filterOptions={state.albumFilterOptions}
            activeFilter={state.albumFilter}
            onFilterChange={thunks.onAlbumFilterChange}
            onShare={thunks.openSharingModal}
        />
    );
}

export interface HomeContentViewProps {
    error?: Error;
    albums: Album[];
    filterOptions: AlbumFilterEntry[];
    activeFilter: AlbumFilterEntry;
    onFilterChange: (filter: AlbumFilterCriterion) => void;
    onShare: (albumId: AlbumId) => void;
}

export const HomeContentView = ({error, albums, filterOptions, activeFilter, onFilterChange, onShare}: HomeContentViewProps) => {
    if (error) {
        return (
            <ErrorMessage error={error} title="Failed to load the albums"/>
        );
    }

    return (
        <AlbumGrid
            albums={albums}
            filterOptions={filterOptions}
            activeFilter={activeFilter}
            onFilterChange={onFilterChange}
            onShare={onShare}
        />
    );
}
