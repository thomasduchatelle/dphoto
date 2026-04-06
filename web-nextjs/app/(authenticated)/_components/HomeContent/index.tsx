'use client';

import {AlbumGrid} from "../AlbumGrid";
import {AlbumId, CatalogViewerState} from "@/domains/catalog";

export const HomeContent = ({initialState}: { initialState: CatalogViewerState }) => (
    <AlbumGrid albums={initialState.albums} onShare={(id: AlbumId) => console.log('onShare', id)}/>
)