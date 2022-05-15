import {Alert, Container, Paper, Toolbar} from "@mui/material";
import {useState} from "react";
import {GoogleLogin, GoogleLoginResponse, GoogleLoginResponseOffline} from "react-google-login";
import AppNavComponent from "../../components/AppNav";
import BackdropLoading from "../../components/BackdropLoading";
import {useConfigContext} from "../../core/application/app-config.context";


function isGoogleLoginResponse(value: GoogleLoginResponse | GoogleLoginResponseOffline): value is GoogleLoginResponse {
  return value.hasOwnProperty('profileObj');
}

const Login = ({googleSignIn, authenticationError}: {
  authenticationError?: string
  googleSignIn(identityToken: string): Promise<void>
}) => {
  const [ready, setReady] = useState(false)
  const [failureMessage, setFailureMessage] = useState("")
  const appConfig = useConfigContext();

  const errorToDisplay = authenticationError ?? failureMessage

  const handleFailure = (error: any) => {
    setFailureMessage(JSON.stringify(error))
  }

  const handleSuccess = (response: GoogleLoginResponse | GoogleLoginResponseOffline): void => {
    if (isGoogleLoginResponse(response)) {
      googleSignIn(response.tokenId).then()
    }
  }

  const handleAutoLoadFinished = (successLogin: boolean): void => {
    // note: component will be unmounted in case of a successful authentication
    setReady(true)
  }

  return (
    <div>
      <Toolbar/>
      <BackdropLoading loading={!ready}/>
      <AppNavComponent
        rightContent={<GoogleLogin
          clientId={appConfig.googleClientId}
          uxMode={appConfig.googleLoginUxMode}
          onFailure={handleFailure}
          onSuccess={handleSuccess}
          onAutoLoadFinished={handleAutoLoadFinished}
          isSignedIn={true}
        />}
      />
      <Container maxWidth='md'>
        {errorToDisplay ? (
          <Paper>
            <Alert severity="error" sx={{mt: 3}}>
              {errorToDisplay}
            </Alert>
          </Paper>
        ) : (
          <Paper sx={{mt: 3}}>
            <Alert severity='info'>
              This is an invitation only application. Sign in with your Google account.
            </Alert>
          </Paper>
        )}
      </Container>
    </div>
  )
}

export default Login
