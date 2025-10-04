// This page shows a specific album based on owner and album params
// The logic is the same as the albums index page - the CatalogViewerProvider
// in the layout will handle routing to the correct album based on URL params

import AlbumsIndexPage from "../../index";

export default function AlbumPage() {
    return <AlbumsIndexPage />;
}
