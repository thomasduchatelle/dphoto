'use client';

import {useMemo, useReducer} from 'react';
import {notFound, useRouter, useSearchParams} from 'next/navigation';
import {catalogReducer, catalogThunks, CatalogViewerState} from '@/domains/catalog';
import {selectFullScreenMediaViewer, photoUrl, albumUrl} from '@/domains/catalog/navigation';
import {displayedAlbumSelector} from '@/domains/catalog/language/selector-displayedAlbum';
import {newClientAdapterFactory} from '@/domains/catalog/adapters/client-adapter-factory';
import {useSessionAwareThunks} from '@/libs/security/client';
import {FullScreenMediaViewer} from '@/components/FullScreenMediaViewer';

export type PhotoViewerCloseBehavior = 'back' | 'album';

export interface PhotoViewerClientProps {
    initialState: CatalogViewerState;
    mediaId: string;
    closeBehavior: PhotoViewerCloseBehavior;
}

export function PhotoViewerClient({initialState, mediaId, closeBehavior}: PhotoViewerClientProps) {
    const [state, dispatch] = useReducer(catalogReducer, initialState);
    const router = useRouter();
    const slideshow = useSearchParams().has('play');

    const adapter = useMemo(() => newClientAdapterFactory(), []);
    const thunks = useSessionAwareThunks(catalogThunks, {dispatch, adapter}, state);

    const {albumId, medias, displayedIndex, mediaNotFound} = selectFullScreenMediaViewer(state, mediaId);

    if (mediaNotFound || !albumId) {
        notFound();
    }

    const handleClose = () => {
        if (closeBehavior === 'back') {
            router.back();
        } else {
            router.push(albumUrl(albumId));
        }
    };

    const handleNavigate = (newMediaId: string) => {
        window.history.replaceState(window.history.state, '', photoUrl(albumId, newMediaId, {autoplay: slideshow}));
    };

    return (
        <FullScreenMediaViewer
            medias={medias}
            displayedIndex={displayedIndex}
            slideshow={slideshow}
            isOwner={displayedAlbumSelector(state).displayedAlbumIdIsOwned}
            onClose={handleClose}
            onNavigate={handleNavigate}
            onAddTag={thunks.addTag}
            onRemoveTag={thunks.removeTag}
        />
    );
}
