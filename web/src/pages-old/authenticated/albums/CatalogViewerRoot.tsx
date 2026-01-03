'use client';

import {AlbumId} from "../../../core/catalog";
import React, {ReactElement, useCallback} from "react";
import {CatalogViewerProvider} from "../../../components/catalog-react";
import {useClientRouter} from "../../../components/ClientRouter";
import {useAuthenticatedUser} from "../../../core/security";

export default function CatalogViewerRoot({children}: {
    children: ReactElement;
}) {
    const {path, params, navigate} = useClientRouter();
    const authenticatedUser = useAuthenticatedUser()

    // Check if we're on an album-specific page
    const pathParts = path.split('/').filter(p => p);
    const albumId = pathParts[0] === 'albums' && pathParts.length >= 3
        ? {owner: params.owner, folderName: params.album} as AlbumId 
        : undefined;
    
    const onSelectedAlbumIdByDefault = useCallback((albumId: AlbumId) => 
        navigate(`/albums/${albumId.owner}/${albumId.folderName}`), 
        [navigate]
    );

    return (
        <CatalogViewerProvider
            key={authenticatedUser.email} // Force unmount and remount when authenticated user changes (required for the catalog state to be reset)
            albumId={albumId}
            redirectToAlbumId={onSelectedAlbumIdByDefault}
            authenticatedUser={authenticatedUser}
        >
            {children}
        </CatalogViewerProvider>
    )
}