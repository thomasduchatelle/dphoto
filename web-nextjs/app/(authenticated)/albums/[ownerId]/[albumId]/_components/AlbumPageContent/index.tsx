'use client';

import {useReducer} from 'react';
import {notFound} from 'next/navigation';
import {Box} from '@mui/material';
import {albumIsOwnedByCurrentUser, catalogReducer, catalogThunks, CatalogViewerState} from '@/domains/catalog';
import {catalogViewerPageSelector} from '@/domains/catalog/navigation/selector-catalog-viewer-page';
import {sharingDialogSelector} from '@/domains/catalog/sharing/selector-sharingDialogSelector';
import {deleteDialogSelector} from '@/domains/catalog/album-delete/selector-deleteDialogSelector';
import {editDatesDialogSelector} from '@/domains/catalog/album-edit-dates/selector-editDatesDialogSelector';
import {editNameDialogSelector} from '@/domains/catalog/album-edit-name/selector-editNameDialogSelector';
import {useThunks} from '@/libs/dthunks/react';
import {ErrorMessage} from '@/components/ErrorMessage';
import {AlbumHeader} from '../AlbumHeader';
import {AlbumRail} from '../AlbumRail';
import {AlbumMediaGrid} from '../AlbumMediaGrid';
import {NeighbourAlbumLink} from '../NeighbourAlbumLink';
import {PreviousAlbumBanner} from '../PreviousAlbumBanner';
import {AlbumActionsFab} from '../AlbumActionsFab';
import {NoMedia} from '../NoMedia';
import {ShareAlbumModalStub} from '../ShareAlbumModalStub';
import {RenameAlbumModalStub} from '../RenameAlbumModalStub';
import {EditDatesModalStub} from '../EditDatesModalStub';
import {DeleteAlbumModalStub} from '../DeleteAlbumModalStub';

export interface AlbumPageContentProps {
    initialState: CatalogViewerState;
}

export function AlbumPageContent({initialState}: AlbumPageContentProps) {
    const [state, dispatch] = useReducer(catalogReducer, initialState);

    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const {onPageRefresh, loadAlbumPage, deleteAlbum, updateAlbumDates, submitCreateAlbum, saveAlbumName, grantAlbumAccess, revokeAlbumAccess, ...dispatchOnlyThunks} = catalogThunks;
    const thunks = useThunks(dispatchOnlyThunks, {dispatch}, state);

    const {displayedAlbum, medias, albums, previousAlbum, nextAlbum, albumNotFound, error} = catalogViewerPageSelector(state);
    const isOwner = displayedAlbum ? albumIsOwnedByCurrentUser(displayedAlbum) : false;
    const shareDialog = sharingDialogSelector(state);
    const deleteDialog = deleteDialogSelector(state);
    const editDatesDialog = editDatesDialogSelector(state);
    const editNameDialog = editNameDialogSelector(state);

    if (error) {
        return <ErrorMessage error={error} title="Failed to load the album"/>;
    }

    if (albumNotFound) {
        notFound();
    }

    const onShare = () => displayedAlbum && thunks.openSharingModal(displayedAlbum.albumId);

    return (
        <Box sx={{display: 'flex', alignItems: 'flex-start', mx: {xs: -2, sm: -3, md: -4}, mt: {xs: -2, sm: -3, md: -4}}}>
            <AlbumRail albums={albums} displayedAlbumId={displayedAlbum?.albumId}/>

            <Box sx={{flex: 1, minWidth: 0}}>
                <AlbumHeader
                    album={displayedAlbum}
                    isOwner={isOwner}
                    onShare={onShare}
                    onRename={thunks.openEditNameDialog}
                    onEditDates={thunks.openEditDatesDialog}
                    onDelete={thunks.openDeleteAlbumDialog}
                />

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
            <AlbumActionsFab
                isOwner={isOwner}
                onShare={onShare}
                onRename={thunks.openEditNameDialog}
                onEditDates={thunks.openEditDatesDialog}
                onDelete={thunks.openDeleteAlbumDialog}
            />

            <ShareAlbumModalStub open={shareDialog.open} onClose={thunks.closeSharingModal}/>
            <RenameAlbumModalStub open={editNameDialog.isOpen} onClose={thunks.closeEditNameDialog}/>
            <EditDatesModalStub open={editDatesDialog.isOpen} onClose={thunks.closeEditDatesDialog}/>
            <DeleteAlbumModalStub open={deleteDialog.isOpen} onClose={thunks.closeDeleteAlbumDialog}/>
        </Box>
    );
}
