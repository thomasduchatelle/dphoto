import React, {useMemo} from "react";
import {Navigate, Route, Routes, useLocation} from "react-router-dom"
import MediaPage from "../albums/[owner]/[album]/[encodedId]/[filename]";
import CatalogViewerRoot from "./albums/CatalogViewerRoot";
import {CatalogViewerPage} from "../albums/[owner]/[album]";

const RedirectToDefaultOrPrevious = () => {
    // note - API Gateway + S3 static will redirect on '/?path=<previously requested url>' when a page is reloaded
    const {search} = useLocation();

    const query = useMemo(() => new URLSearchParams(search), [search]);
    const redirectTo = query.get("path") ?? '/albums'
    return (
        <Navigate to={redirectTo}/>
    )
}

const AuthenticatedRouter = () => {
    return (
        <Routes>
            <Route path='/albums' element={<CatalogViewerRoot><CatalogViewerPage/></CatalogViewerRoot>}/>
            <Route path='/albums/:owner/:album' element={<CatalogViewerRoot><CatalogViewerPage/></CatalogViewerRoot>}/>
            <Route path='/albums/:owner/:album/:encodedId/:filename' element={<CatalogViewerRoot><MediaPage/></CatalogViewerRoot>}/>
            <Route path='*' element={<RedirectToDefaultOrPrevious/>}/>
        </Routes>
    )
}

export default AuthenticatedRouter
