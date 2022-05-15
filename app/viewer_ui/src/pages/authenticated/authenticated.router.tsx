import {Box} from "@mui/material";
import React, {useMemo} from "react";
import {Navigate, Route, Routes, useLocation} from "react-router-dom"
import AppNavComponent from "../../components/app-nav.component";
import UserMenu from "../../components/user.menu";
import {useMustBeAuthenticated} from "../../core/application";
import AlbumRouterPage from "./albums/AlbumRouterPage";

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
  const {loggedUser, signOutCase} = useMustBeAuthenticated()

  return (
    <Box sx={{display: 'flex'}}>
      <AppNavComponent
        rightContent={<UserMenu user={loggedUser} onLogout={signOutCase.logout}/>}
      />
      <Routes>
        <Route path='/albums' element={<AlbumRouterPage/>}/>
        <Route path='/albums/:owner/:album' element={<AlbumRouterPage/>}/>
        <Route path='*' element={<RedirectToDefaultOrPrevious/>}/>
      </Routes>
    </Box>
  )
}

export default AuthenticatedRouter
