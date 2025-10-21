export interface RedirectProps {
  url: string;
}

export const Redirect = ({ url }: RedirectProps) => {
  return (
    <html>
      <head>
        <meta httpEquiv="refresh" content={`0;url=${url}`} />
      </head>
      <body>
        <p>Redirecting...</p>
      </body>
    </html>
  );
};
