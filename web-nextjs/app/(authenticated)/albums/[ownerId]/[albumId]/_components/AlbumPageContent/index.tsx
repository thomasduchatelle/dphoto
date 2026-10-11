'use client';

import {useMemo, useReducer} from 'react';
import {notFound} from 'next/navigation';
import {Box} from '@mui/material';
import {catalogReducer, catalogThunks, CatalogViewerState} from '@/domains/catalog';
import {CatalogViewerPageSelection, catalogViewerPageSelector} from '@/domains/catalog/navigation/selector-catalog-viewer-page';
import {displayedAlbumSelector} from '@/domains/catalog/language/selector-displayedAlbum';
import {deleteDialogSelector} from '@/domains/catalog/album-delete/selector-deleteDialogSelector';
import {editDatesDialogSelector} from '@/domains/catalog/album-edit-dates/selector-editDatesDialogSelector';
import {editNameDialogSelector} from '@/domains/catalog/album-edit-name/selector-editNameDialogSelector';
import {newClientAdapterFactory} from '@/domains/catalog/adapters/client-adapter-factory';
import {useSessionAwareThunks} from '@/libs/security/client';
import {ErrorMessage} from '@/components/ErrorMessage';
import {AlbumHeader} from '../AlbumHeader';
import {AlbumRail} from '../AlbumRail';
import {AlbumMediaGrid} from '../AlbumMediaGrid';
import {NeighbourAlbumLink} from '../NeighbourAlbumLink';
import {PreviousAlbumBanner} from '../PreviousAlbumBanner';
import {AlbumActionsFab} from '../AlbumActionsFab';
import {NoMedia} from '../NoMedia';
import {DeleteAlbumModal} from '@/components/DeleteAlbumModal';
import {EditDatesModal} from '../EditDatesModal';
import {EditNameDialog} from '../EditNameDialog';

export interface AlbumPageContentProps {
    initialState: CatalogViewerState;
}

export function AlbumPageContent({initialState}: AlbumPageContentProps) {
    const [state, dispatch] = useReducer(catalogReducer, initialState);

    const adapter = useMemo(() => newClientAdapterFactory(), []);
    const thunks = useSessionAwareThunks(catalogThunks, {dispatch, adapter}, state);

    const selection = catalogViewerPageSelector(state);

    if (!selection.error && selection.albumNotFound) {
        notFound();
    }

    return (
        <>
            <AlbumPageContentView
                {...selection}
                isOwner={displayedAlbumSelector(state).displayedAlbumIdIsOwned}
                onRename={thunks.openEditNameDialog}
                onEditDates={thunks.openEditDatesDialog}
                onDelete={thunks.openDeleteAlbumDialog}
            />

            <DeleteAlbumModal
                {...deleteDialogSelector(state)}
                onClose={thunks.closeDeleteAlbumDialog}
                onDelete={thunks.deleteAlbum}
            />
            <EditDatesModal
                {...editDatesDialogSelector(state)}
                onClose={thunks.closeEditDatesDialog}
                onSave={thunks.updateAlbumDates}
                onStartDateChange={thunks.updateEditDatesDialogStartDate}
                onEndDateChange={thunks.updateEditDatesDialogEndDate}
                onStartAtDayStartChange={thunks.updateEditDatesDialogStartAtDayStart}
                onEndAtDayEndChange={thunks.updateEditDatesDialogEndAtDayEnd}
            />
            <EditNameDialog
                {...editNameDialogSelector(state)}
                onClose={thunks.closeEditNameDialog}
                onSave={thunks.saveAlbumName}
                onAlbumNameChange={thunks.changeAlbumName}
                onFolderNameEnabledChange={thunks.changeFolderNameEnabled}
                onFolderNameChange={thunks.changeFolderName}
            />
        </>
    );
}

export interface AlbumPageContentViewProps extends CatalogViewerPageSelection {
    isOwner: boolean;
    onRename: () => void;
    onEditDates: () => void;
    onDelete: () => void;
}

export function AlbumPageContentView({
                                         albums,
                                         displayedAlbum,
                                         medias,
                                         previousAlbum,
                                         nextAlbum,
                                         error,
                                         isOwner,
                                         onRename,
                                         onEditDates,
                                         onDelete,
                                     }: AlbumPageContentViewProps) {
    if (error) {
        return <ErrorMessage error={error} title="Failed to load the album"/>;
    }

    return (
        <Box sx={{display: 'flex', alignItems: 'flex-start', mx: {xs: -2, sm: -3, md: -4}, mt: {xs: -2, sm: -3, md: -4}}}>
            <AlbumRail albums={albums} displayedAlbumId={displayedAlbum?.albumId}/>

            <Box sx={{flex: 1, minWidth: 0}}>
                <AlbumHeader
                    album={displayedAlbum}
                    isOwner={isOwner}
                    onRename={onRename}
                    onEditDates={onEditDates}
                    onDelete={onDelete}
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
                onRename={onRename}
                onEditDates={onEditDates}
                onDelete={onDelete}
            />
        </Box>
    );
}
