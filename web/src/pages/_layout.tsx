import { Providers } from '../components/Providers';
import { AuthProvider } from '../components/AuthProvider';
import GeneralRouter from '../pages-old/GeneralRouter';

export default async function RootLayout({ children }: { children: React.ReactNode }) {
  // In Waku, we don't have direct access to request headers in RSC
  // The cookie header will be undefined here, but cookies will work via browser
  return (
    <AuthProvider>
      <div className="App">
        <Providers>
          <GeneralRouter />
        </Providers>
      </div>
    </AuthProvider>
  );
}
