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
import {NextAlbumBanner} from '../NextAlbumBanner';
import {AlbumActionsFab} from '../AlbumActionsFab';
import {NoMedia} from '../NoMedia';

export interface AlbumPageContentProps {
    initialState: CatalogViewerState;
}

export function AlbumPageContent({initialState}: AlbumPageContentProps) {
    const [state, dispatch] = useReducer(catalogReducer, initialState);

    const {onPageRefresh, loadAlbumPage, deleteAlbum, updateAlbumDates, submitCreateAlbum, saveAlbumName, grantAlbumAccess, revokeAlbumAccess, ...dispatchOnlyThunks} = catalogThunks;
    useThunks(dispatchOnlyThunks, {dispatch}, state);

    const {displayedAlbum, medias, mediasLoaded, albums, previousAlbum, nextAlbum, albumNotFound, error} = catalogViewerPageSelector(state);

    if (error) {
        return <ErrorMessage error={error} title="Failed to load the album"/>;
    }

    if (albumNotFound) {
        notFound();
    }

    if (!mediasLoaded) {
        return null;
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
                        {(nextAlbum || previousAlbum) && (
                            <Box
                                sx={{
                                    display: {xs: 'flex', lg: 'none'},
                                    flexDirection: 'column',
                                    gap: 1.5,
                                    px: 1.5,
                                    pt: 2,
                                    pb: 10,
                                    borderTop: '1px solid rgba(255,255,255,0.07)',
                                }}
                            >
                                {nextAlbum && <NeighbourAlbumLink album={nextAlbum}/>}
                                {previousAlbum && <NeighbourAlbumLink album={previousAlbum}/>}
                            </Box>
                        )}
                    </>
                )}
            </Box>

            <NextAlbumBanner album={nextAlbum}/>
            <AlbumActionsFab/>
        </Box>
    );
}
