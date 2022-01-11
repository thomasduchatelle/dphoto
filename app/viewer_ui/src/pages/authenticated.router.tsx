import {GoogleLogout} from "react-google-login";
import {Route, Routes} from "react-router-dom"
import googleConfig from "../config/google.config";
import Home from "../Home";

const SomeContent = () => {
  return (
    <>
      <GoogleLogout clientId={googleConfig.clientId}/>
      <Home/>
    </>
  )
}

export default () => {

  console.log("authenticated router")
  return (
    <Routes>
      <Route path='/*' element={<SomeContent/>}/>
    </Routes>
  )
}