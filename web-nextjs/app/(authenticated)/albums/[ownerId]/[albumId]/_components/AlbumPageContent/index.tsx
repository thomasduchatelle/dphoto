'use client';

import {useReducer} from 'react';
import {notFound} from 'next/navigation';
import {Box} from '@mui/material';
import {catalogReducer, catalogThunks, CatalogViewerState} from '@/domains/catalog';
import {catalogViewerPageSelector} from '@/domains/catalog/navigation/selector-catalog-viewer-page';
import {useThunks} from '@/libs/dthunks/react';
import {ErrorMessage} from '@/components/ErrorMessage';
import {AlbumHeader} from '../AlbumHeader';
import {AlbumRail} from '../AlbumRail';
import {AlbumMediaGrid} from '../AlbumMediaGrid';
import {NeighbourAlbumLink} from '../NeighbourAlbumLink';
import {PreviousAlbumBanner} from '../PreviousAlbumBanner';
import {AlbumActionsFab} from '../AlbumActionsFab';
import {NoMedia} from '../NoMedia';

export interface AlbumPageContentProps {
    initialState: CatalogViewerState;
}

export function AlbumPageContent({initialState}: AlbumPageContentProps) {
    const [state, dispatch] = useReducer(catalogReducer, initialState);

    const {onPageRefresh, loadAlbumPage, deleteAlbum, updateAlbumDates, submitCreateAlbum, saveAlbumName, grantAlbumAccess, revokeAlbumAccess, ...dispatchOnlyThunks} = catalogThunks;
    useThunks(dispatchOnlyThunks, {dispatch}, state);

    const {displayedAlbum, medias, albums, previousAlbum, nextAlbum, albumNotFound, error} = catalogViewerPageSelector(state);

    if (error) {
        return <ErrorMessage error={error} title="Failed to load the album"/>;
    }

    if (albumNotFound) {
        notFound();
    }

    return (
        <Box sx={{display: 'flex', alignItems: 'flex-start', mx: {xs: -2, sm: -3, md: -4}, mt: {xs: -2, sm: -3, md: -4}}}>
            <AlbumRail albums={albums} displayedAlbumId={displayedAlbum?.albumId}/>

            <Box sx={{flex: 1, minWidth: 0}}>
                <AlbumHeader album={displayedAlbum}/>

                {medias.length === 0 ? (
                    <NoMedia nextAlbum={nextAlbum} previousAlbum={previousAlbum}/>
                ) : (
                    <>
                        <AlbumMediaGrid medias={medias}/>
                        {nextAlbum && (
                            <Box
                                sx={{
                                    display: {xs: 'block', lg: 'none'},
                                    px: 1.5,
                                    pt: 2,
                                    pb: 10,
                                    borderTop: '1px solid rgba(255,255,255,0.07)',
                                }}
                            >
                                <NeighbourAlbumLink album={nextAlbum} direction="next"/>
                            </Box>
                        )}
                    </>
                )}
            </Box>

            <PreviousAlbumBanner album={previousAlbum}/>
            <AlbumActionsFab/>
        </Box>
    );
}
