import {Alert, Container, CssBaseline} from "@mui/material";
import {MouseEvent, useContext} from "react";
import {useGoogleLogin} from "react-google-login";
import googleConfig from "../../config/google.config";
import {SecurityContext} from "../layout/security.context";

export default () => {
  const securityContext = useContext(SecurityContext);

  const {signIn, loaded} = useGoogleLogin({clientId: googleConfig.clientId})

  if (!loaded) {
    return null
  }

  const handleLogin = (event: MouseEvent<HTMLDivElement>) => {
    event.preventDefault()
    // securityContext.signInWithGoogle()
    // signIn()
  };

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
      {!securityContext.authenticationError ? (
        <Alert severity="info" sx={{mt: 3, mb: 10}}>
          This is an invitation only application. Sign in with your Google account.
        </Alert>
      ) : (
        <Alert severity="error" sx={{mt: 3, mb: 10}}>
          {securityContext.authenticationError}
        </Alert>
      )}
      {/*<script src="https://accounts.google.com/gsi/client" async defer />*/}
      <div id="g_id_onload"
           data-client_id="YOUR_GOOGLE_CLIENT_ID"
           data-login_uri="https://your.domain/your_login_endpoint"
           data-auto_prompt="false">
      </div>
      <div className="g_id_signin"
           data-type="standard"
           data-size="large"
           data-theme="outline"
           data-text="sign_in_with"
           data-shape="rectangular"
           data-logo_alignment="left">
      </div>
      {/*<GoogleButton type='light' onClick={handleLogin} style={{margin: 'auto'}}/>*/}
      {/*<GoogleLogin clientId={googleConfig.clientId} />*/}
    </Container>
  )
}
