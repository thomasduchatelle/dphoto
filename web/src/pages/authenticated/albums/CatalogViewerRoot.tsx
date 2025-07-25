import {useMatch, useNavigate} from "react-router-dom";
import {AlbumId} from "../../../core/catalog";
import React, {ReactElement, useCallback} from "react";
import {CatalogViewerProvider} from "../../../components/catalog-react";
import {useAuthenticatedUser} from "../../../core/application";

export default function CatalogViewerRoot({children}: {
    children: ReactElement;
}) {
    const match = useMatch('/albums/:owner/:folderName/*');
    const navigate = useNavigate()
    const authenticatedUser = useAuthenticatedUser()

    const albumId = match ? {owner: match.params.owner, folderName: match.params.folderName} as AlbumId : undefined
    const onSelectedAlbumIdByDefault = useCallback((albumId: AlbumId) => navigate(`/albums/${albumId.owner}/${albumId.folderName}`), [navigate])

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