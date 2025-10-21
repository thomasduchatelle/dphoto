export interface CookiesProps {
  cookies: Record<string, string>;
}

export const Cookies = ({ cookies }: CookiesProps) => {
  return (
    <>
      {Object.entries(cookies).map(([name, value]) => (
        <meta key={name} httpEquiv="Set-Cookie" content={`${name}=${value}; Path=/`} />
      ))}
    </>
  );
};
