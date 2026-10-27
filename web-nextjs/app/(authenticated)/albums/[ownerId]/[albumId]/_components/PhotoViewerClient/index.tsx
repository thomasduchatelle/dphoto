'use client';

import {useReducer} from 'react';
import {notFound, useRouter, useSearchParams} from 'next/navigation';
import {catalogReducer, CatalogViewerState} from '@/domains/catalog';
import {selectFullScreenMediaViewer, photoUrl, albumUrl} from '@/domains/catalog/navigation';
import {FullScreenMediaViewer} from '@/components/FullScreenMediaViewer';

export type PhotoViewerCloseBehavior = 'back' | 'album';

export interface PhotoViewerClientProps {
    initialState: CatalogViewerState;
    mediaId: string;
    closeBehavior: PhotoViewerCloseBehavior;
}

export function PhotoViewerClient({initialState, mediaId, closeBehavior}: PhotoViewerClientProps) {
    const [state] = useReducer(catalogReducer, initialState);
    const router = useRouter();
    const slideshow = useSearchParams().has('play');

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
            onClose={handleClose}
            onNavigate={handleNavigate}
        />
    );
}
