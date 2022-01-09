import {Alert, Container, CssBaseline} from "@mui/material";
import {MouseEvent} from "react";
import {GoogleLogin} from "../../components/google-signin-button";

export default () => {
  const handleLogin = (event: MouseEvent<HTMLButtonElement>) => {
    event.preventDefault()
    console.log("> CLICK")
  }
  return (
    <Container component="div"
               sx={{
                 marginTop: 8,
                 width: '650px',
                 textAlign: "center",
                 margin: '8 auto 0 auto'
               }}>
      <CssBaseline/>
      <img src="/dphoto-fulllogo-large.png" alt="dphoto-logo"/>
      <Alert severity="info" sx={{mt: 3, mb: 10}}>This is an invitation only application. Sign in with your Google
        account.</Alert>
      <GoogleLogin onClick={handleLogin}/>
    </Container>
  )
}
