import React from "react";
import {GoogleLogin, GoogleLoginResponse, GoogleLoginResponseOffline} from "react-google-login";

const LoginPage = (props: {
  clientId: string
  isSignedIn: boolean
  authenticateWithGoogle: (googleAnswer: (GoogleLoginResponse | GoogleLoginResponseOffline)) => void
  onFailure: (error: any) => void
  errorMessage: string
}) => {
  return (
    <header className="App-header">
      <p>Use your google account to access your pictures</p>
      {props.errorMessage && <p>{props.errorMessage}</p>}
      <GoogleLogin
        clientId={props.clientId}
        buttonText="Login"
        onSuccess={props.authenticateWithGoogle}
        onFailure={props.onFailure}
        cookiePolicy={'single_host_origin'}
        isSignedIn={props.isSignedIn}
      />

    </header>
  );
}

export default LoginPage