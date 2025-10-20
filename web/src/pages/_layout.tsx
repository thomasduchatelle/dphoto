import { Providers } from '../components/Providers';
import GeneralRouter from '../pages-old/GeneralRouter';

export default async function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="App">
      <Providers>
        <GeneralRouter />
      </Providers>
    </div>
  );
}
