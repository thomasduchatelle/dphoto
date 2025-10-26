'use client';

import React from "react";
import MediaPage from "./media";
import CatalogViewerRoot from "./albums/CatalogViewerRoot";
import {CatalogViewerPage} from "./albums/CatalogViewerPage";
import {useClientRouter} from "../../components/ClientRouter";

const RedirectToDefaultOrPrevious = () => {
    // note - API Gateway + S3 static will redirect on '/?path=<previously requested url>' when a page is reloaded
    const {query, navigate} = useClientRouter();

    const redirectTo = query.get("path") ?? '/albums'

    React.useEffect(() => {
        navigate(redirectTo);
    }, [redirectTo, navigate]);

    return null;
}

const AuthenticatedRouter = () => {
    const {path} = useClientRouter();

    // Parse the path to determine which component to render
    const pathParts = path.split('/').filter(p => p);

    console.log("Rendering MediaPage for path:", path);

    // /albums/:owner/:album/:encodedId/:filename
    if (pathParts[0] === 'albums' && pathParts.length >= 5) {
        console.log("rendering path '/albums/:owner/:album/:encodedId/:filename'");
        return <CatalogViewerRoot><MediaPage/></CatalogViewerRoot>;
    }

    // /albums/:owner/:album
    if (pathParts[0] === 'albums' && pathParts.length >= 3) {
        console.log("rendering path '/albums/:owner/:album'");
        return <CatalogViewerRoot><CatalogViewerPage/></CatalogViewerRoot>;
    }

    // /albums
    if (pathParts[0] === 'albums') {
        console.log("rendering path '/albums'");
        return <CatalogViewerRoot><CatalogViewerPage/></CatalogViewerRoot>
    }

    // Default: redirect to /albums
    console.log("Default: redirect to /albums");
    return <RedirectToDefaultOrPrevious/>;
}

export default AuthenticatedRouter
