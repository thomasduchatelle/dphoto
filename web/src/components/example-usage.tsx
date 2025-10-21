// Example: Using Redirect and Cookies components together
// This demonstrates a typical authentication callback scenario

import { Cookies } from './Cookies';

export default function AuthCallbackExample() {
  // After successful authentication, we have tokens to set
  const tokens = {
    access_token: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...',
    refresh_token: 'rt_abc123xyz789'
  };

  // Redirect to the dashboard after setting cookies
  const redirectUrl = '/dashboard';

  return (
    <html>
      <head>
        <meta httpEquiv="refresh" content={`0;url=${redirectUrl}`} />
        <Cookies cookies={tokens} />
      </head>
      <body>
        <p>Authentication successful. Redirecting...</p>
      </body>
    </html>
  );
}

// This is similar to the existing implementation in:
// /home/runner/work/dphoto/dphoto/web/src/pages/auth/callback/index.tsx
// Lines 119-131
