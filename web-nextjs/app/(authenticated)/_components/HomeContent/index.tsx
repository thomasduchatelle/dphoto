'use client';

import {AlbumGrid} from "../AlbumGrid";
import {AlbumFilterEntry, albumsFiltered, catalogReducer, CatalogViewerState} from "@/domains/catalog";
import {useCallback, useReducer} from "react";
import {ErrorMessage} from "@/components/ErrorMessage";

export const HomeContent = ({initialState}: { initialState: CatalogViewerState }) => {
    const [state, dispatch] = useReducer(catalogReducer, initialState);

    const handleFilterChange = useCallback((newFilter: AlbumFilterEntry) => {
        dispatch(albumsFiltered({criterion: newFilter.criterion}));
    }, []);

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
            onFilterChange={handleFilterChange}
            onShare={() => {
            }}
        />
    );
}