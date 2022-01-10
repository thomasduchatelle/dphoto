This is a fork of https://github.com/anthonyjgrove/react-google-login to keep the sign-in button without any logic. Use directly the library for anything non-business.

No update should be required but just in case:

1. clone and copy the content of the GIT repo `src/` directory in this one
2. copy `index.d.ts` fro th ereop root to this folder
3. remove the use of `useGoogleLogin()` from GoogleLogin component and make `clientId` optional in `index.d.ts`
