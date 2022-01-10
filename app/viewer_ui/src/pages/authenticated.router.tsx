import {Button, Typography} from "@mui/material";
import {useContext} from "react";
import {Route, Routes} from "react-router-dom"
import {SecurityContext} from "./layout/security.context";

const SomeContent = () => {
  const securityContext = useContext(SecurityContext);
  return (
    <>
      <Typography component='h1'>You are authenticated!</Typography>
      <Typography component='p'>Welcome authenticated user !</Typography>
      <Button onClick={securityContext.signOut} variant='contained' color='error'>Sign out</Button>
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