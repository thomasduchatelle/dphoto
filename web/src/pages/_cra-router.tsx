import React, {ReactElement, useState} from "react";
import {Navigate, Route, Routes, useSearchParams} from "react-router-dom";
import {useGlobalError, useIsAuthenticated} from "../core/application";
import IndexPage from "./index";
import LoginPage from "./login";
import AlbumsIndexPage from "./albums/index";
import AlbumPage from "./albums/[owner]/[album]/index";
import MediaPage from "./albums/[owner]/[album]/[encodedId]/[filename]";
import AlbumsLayout from "./albums/_layout";
import ErrorPage from "./error-page";

const RestoreAPIGatewayOriginalPath = ({children}: {
    children: ReactElement;
}) => {
    // note - API Gateway + S3 static will redirect on '/?path=<previously requested url>' when a page is reloaded
    const [search] = useSearchParams();
    const path = search.get('path')
    if (path) {
        return (
            <Navigate to={path}/>
        )
    }

    return children
}

const CRARouter = () => {
    const isAuthenticated = useIsAuthenticated();
    const [displayLoginPage] = useState<boolean>(!isAuthenticated)
    const globalError = useGlobalError()

    if (globalError) {
        return <ErrorPage error={globalError}/>
    }

    if (!isAuthenticated || displayLoginPage) {
        return (
            <RestoreAPIGatewayOriginalPath>
                <Routes>
                    <Route path="*" element={<LoginPage />} />
                </Routes>
            </RestoreAPIGatewayOriginalPath>
        )
    }

    return (
        <Routes>
            <Route path="/" element={<IndexPage />} />
            <Route path="/albums" element={<AlbumsLayout><AlbumsIndexPage /></AlbumsLayout>} />
            <Route path="/albums/:owner/:album" element={<AlbumsLayout><AlbumPage /></AlbumsLayout>} />
            <Route path="/albums/:owner/:album/:encodedId/:filename" element={<AlbumsLayout><MediaPage /></AlbumsLayout>} />
            <Route path="*" element={<Navigate to="/albums" />} />
        </Routes>
    )
}

export default CRARouter
