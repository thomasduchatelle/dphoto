'use client';

import {useEffect, useMemo, useReducer} from 'react';
import {notFound, useParams, useRouter} from 'next/navigation';
import {Box} from '@mui/material';
import {catalogReducer, catalogThunks, CatalogViewerState} from '@/domains/catalog';
import {CatalogViewerPageSelection, catalogViewerPageSelector} from '@/domains/catalog/navigation/selector-catalog-viewer-page';
import {albumUrl} from '@/domains/catalog/navigation/album-url';
import {displayedAlbumSelector} from '@/domains/catalog/language/selector-displayedAlbum';
import {createDialogSelector} from '@/domains/catalog/album-create';
import {deleteDialogSelector} from '@/domains/catalog/album-delete/selector-deleteDialogSelector';
import {editDatesDialogSelector} from '@/domains/catalog/album-edit-dates/selector-editDatesDialogSelector';
import {editNameDialogSelector} from '@/domains/catalog/album-edit-name/selector-editNameDialogSelector';
import {sharingDialogSelector} from '@/domains/catalog/sharing/selector-sharingDialogSelector';
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
import {CreateAlbumModal} from '../../../../../_components/CreateAlbumModal';
import {DeleteAlbumModal} from '@/components/DeleteAlbumModal';
import {EditDatesModal} from '../EditDatesModal';
import {EditNameDialog} from '../EditNameDialog';
import {ShareAlbumModal} from '@/components/ShareAlbumModal';

export interface AlbumPageContentProps {
    initialState: CatalogViewerState;
}

export function AlbumPageContent({initialState}: AlbumPageContentProps) {
    const [state, dispatch] = useReducer(catalogReducer, initialState);
    const router = useRouter();
    const params = useParams<{ ownerId: string; albumId: string }>();

    const adapter = useMemo(() => newClientAdapterFactory(), []);
    const thunks = useSessionAwareThunks(catalogThunks, {dispatch, adapter}, state);

    const selection = catalogViewerPageSelector(state);

    const displayedAlbumId = selection.displayedAlbum?.albumId;
    const routeAlbumId = {owner: decodeURIComponent(params.ownerId), folderName: decodeURIComponent(params.albumId)};
    useEffect(() => {
        if (displayedAlbumId && (displayedAlbumId.owner !== routeAlbumId.owner || displayedAlbumId.folderName !== routeAlbumId.folderName)) {
            router.replace(albumUrl(displayedAlbumId));
        }
    }, [displayedAlbumId, routeAlbumId.owner, routeAlbumId.folderName, router]);

    if (!selection.error && selection.albumNotFound) {
        notFound();
    }

    const firstMediaHref = selection.firstMediaHref;

    return (
        <>
            <AlbumPageContentView
                {...selection}
                isOwner={displayedAlbumSelector(state).displayedAlbumIdIsOwned}
                onPlay={firstMediaHref ? () => router.push(firstMediaHref) : undefined}
                onShare={() => displayedAlbumId && thunks.openSharingModal(displayedAlbumId)}
                onRename={thunks.openEditNameDialog}
                onEditDates={thunks.openEditDatesDialog}
                onDelete={thunks.openDeleteAlbumDialog}
                onCreateAlbum={thunks.openCreateDialog}
            />

            <CreateAlbumModal
                {...createDialogSelector(state)}
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
            <ShareAlbumModal
                {...sharingDialogSelector(state)}
                onClose={thunks.closeSharingModal}
                onGrant={thunks.grantAlbumAccess}
                onRevoke={thunks.revokeAlbumAccess}
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
    onPlay: (() => void) | undefined;
    onShare: () => void;
    onRename: () => void;
    onEditDates: () => void;
    onDelete: () => void;
    onCreateAlbum: () => void;
}

export function AlbumPageContentView({
                                         railAlbums,
                                         displayedAlbum,
                                         medias,
                                         firstMediaHref,
                                         previousAlbum,
                                         nextAlbum,
                                         error,
                                         isOwner,
                                         onPlay,
                                         onShare,
                                         onRename,
                                         onEditDates,
                                         onDelete,
                                         onCreateAlbum,
                                     }: AlbumPageContentViewProps) {
    if (error) {
        return <ErrorMessage error={error} title="Failed to load the album"/>;
    }

    return (
        <Box sx={{display: 'flex', alignItems: 'flex-start', mx: {xs: -2, sm: -3, md: -4}, mt: {xs: -2, sm: -3, md: -4}}}>
            <AlbumRail albums={railAlbums} displayedAlbumId={displayedAlbum?.albumId}/>

            <Box sx={{flex: 1, minWidth: 0}}>
                <AlbumHeader
                    album={displayedAlbum}
                    isOwner={isOwner}
                    playHref={firstMediaHref}
                    onShare={onShare}
                    onRename={onRename}
                    onEditDates={onEditDates}
                    onDelete={onDelete}
                    onCreateAlbum={onCreateAlbum}
                />

                {medias.length === 0 || !displayedAlbum ? (
                    <NoMedia nextAlbum={nextAlbum} previousAlbum={previousAlbum}/>
                ) : (
                    <>
                        <AlbumMediaGrid albumId={displayedAlbum.albumId} medias={medias}/>
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
                onPlay={onPlay}
                onShare={onShare}
                onRename={onRename}
                onEditDates={onEditDates}
                onDelete={onDelete}
                onCreateAlbum={onCreateAlbum}
            />
        </Box>
    );
}
