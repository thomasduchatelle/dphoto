const awsConfig = {

  // REQUIRED only for Federated Authentication - Amazon Cognito Identity Pool ID
  // identityPoolId: 'XX-XXXX-X:XXXXXXXX-XXXX-1234-abcd-1234567890ab',

  // REQUIRED - Amazon Cognito Region
  region: 'eu-west-1',

  // OPTIONAL - Amazon Cognito Federated Identity Pool Region
  // Required only if it's different from Amazon Cognito Region
  // identityPoolRegion: 'XX-XXXX-X',

  // OPTIONAL - Amazon Cognito User Pool ID
  userPoolId: 'eu-west-1_Lko0PucSr',

  // OPTIONAL - Amazon Cognito Web Client ID (26-char alphanumeric string)
  userPoolWebClientId: '72011a6soump7p5pveu5182obh',

  // OPTIONAL - Enforce user authentication prior to accessing AWS resources or not
  mandatorySignIn: true,

  // OPTIONAL - Configuration for cookie storage
  // Note: if the secure flag is set to true, then the cookie transmission requires a secure protocol
  // cookieStorage: {
  //   // REQUIRED - Cookie domain (only required if cookieStorage is provided)
  //   domain: '.yourdomain.com',
  //   // OPTIONAL - Cookie path
  //   path: '/',
  //   // OPTIONAL - Cookie expiration in days
  //   expires: 365,
  //   // OPTIONAL - See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Set-Cookie/SameSite
  //   sameSite: "strict" | "lax",
  //   // OPTIONAL - Cookie secure flag
  //   // Either true or false, indicating if the cookie transmission requires a secure protocol (https).
  //   secure: true
  // },

  // OPTIONAL - customized storage object
  // storage: MyStorage,

  // OPTIONAL - Manually set the authentication flow type. Default is 'USER_SRP_AUTH'
  // authenticationFlowType: 'USER_PASSWORD_AUTH',

  // OPTIONAL - Manually set key value pairs that can be passed to Cognito Lambda Triggers
  // clientMetadata: { myCustomKey: 'myCustomValue' },

  // OPTIONAL - Hosted UI configuration
  oauth: {
    domain: 'dphoto-spike.auth.eu-west-1.amazoncognito.com',
    scope: ['phone', 'email', 'openid'],
    redirectSignIn: 'http://localhost:3000/',
    redirectSignOut: 'http://localhost:3000/',
    responseType: 'code' // or 'token', note that REFRESH token will only be generated when the responseType is code
  }
}
// aws_app_analytics: 'enable',
//
// aws_user_pools: 'enable',
// aws_user_pools_id: 'eu-west-1_x',
// aws_user_pools_mfa_type: 'OFF',
// aws_user_pools_web_client_id: '72011a6soump7p5pveu5182obh',
// aws_user_settings: 'enable',
// };

export default awsConfig