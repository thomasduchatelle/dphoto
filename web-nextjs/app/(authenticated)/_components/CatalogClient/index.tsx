'use client';

import {useReducer, useCallback} from 'react';
import {CatalogViewerState, AlbumFilterEntry} from '@/domains/catalog/language/catalog-state';
import {catalogReducer} from '@/domains/catalog/actions';
import {albumsFiltered} from '@/domains/catalog/navigation/action-albumsFiltered';
import {HomePageContent} from '../HomePageContent';

export interface CatalogClientProps {
    initialState: CatalogViewerState;
}

export function CatalogClient({initialState}: CatalogClientProps) {
    const [state, dispatch] = useReducer(catalogReducer, initialState);

    const handleFilterChange = useCallback((newFilter: AlbumFilterEntry) => {
        dispatch(albumsFiltered({criterion: newFilter.criterion}));
    }, []);

    return (
        <HomePageContent
            albums={state.albums}
            filterOptions={state.albumFilterOptions}
            activeFilter={state.albumFilter}
            onFilterChange={handleFilterChange}
            error={state.error}
        />
    );
}
