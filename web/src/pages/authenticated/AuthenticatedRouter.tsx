import React, {useMemo} from "react";
import {Navigate, Route, Routes, useLocation} from "react-router-dom"
import AlbumRouterPage from "./albums/AlbumRouterPage";
import MediaPage from "./media";

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
      <Route path='/albums' element={<AlbumRouterPage/>}/>
      <Route path='/albums/:owner/:album' element={<AlbumRouterPage/>}/>
      <Route path='/albums/:owner/:album/:encodedId/:filename' element={<MediaPage/>}/>
      <Route path='*' element={<RedirectToDefaultOrPrevious/>}/>
    </Routes>
  )
}

export default AuthenticatedRouter
