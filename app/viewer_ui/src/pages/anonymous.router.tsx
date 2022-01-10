import {useContext} from "react";
import {Route, Routes} from 'react-router';
import AuthenticatedRouter from "./authenticated.router";
import {SecurityContext} from "./layout/security.context";
import LoginPage from "./login/login.page";

export default () => {
  const securityContext = useContext(SecurityContext);
  return (
    <Routes>
      <Route path='/login' element={<LoginPage/>}/>
      {securityContext.user ? (
        <Route path='/*' element={<AuthenticatedRouter/>}/>
      ) : (
        <Route>Not authorised.</Route>
      )}
    </Routes>
  )
}