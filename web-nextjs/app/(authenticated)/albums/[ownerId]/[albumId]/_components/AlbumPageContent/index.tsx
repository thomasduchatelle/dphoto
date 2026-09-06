'use client';

import {useReducer} from 'react';
import {notFound} from 'next/navigation';
import {catalogReducer, catalogThunks, CatalogViewerState} from '@/domains/catalog';
import {useThunks} from '@/libs/dthunks/react';
import {ErrorMessage} from '@/components/ErrorMessage';
import {NoMedia} from '../NoMedia';
import Link from '@/components/Link';
import {Button} from '@mui/material';

export function AlbumPageContent({initialState}: { initialState: CatalogViewerState }) {
    const [state, dispatch] = useReducer(catalogReducer, initialState);

    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const {onPageRefresh, loadAlbumPage, deleteAlbum, updateAlbumDates, submitCreateAlbum, saveAlbumName, grantAlbumAccess, revokeAlbumAccess, ...dispatchOnlyThunks} = catalogThunks;
    useThunks(dispatchOnlyThunks, {dispatch}, state);

    if (state.error) {
        return <ErrorMessage error={state.error} title="Failed to load the album"/>;
    }

    if (state.albumNotFound) {
        notFound();
    }

    if (!state.mediasLoaded || state.medias.length === 0) {
        return <NoMedia/>;
    }

    return (
        <div>
            <Button component={Link} href="/" prefetch={false} variant="text">
                Back to Albums
            </Button>
            <ul>
                {state.medias.flatMap(({medias}) =>
                    medias.map(media => (
                        <li key={media.id}>
                            {media.uiRelativePath} — {media.time.toISOString()}
                        </li>
                    ))
                )}
            </ul>
        </div>
    );
}
