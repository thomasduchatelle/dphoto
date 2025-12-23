import {Providers} from '../components/Providers';

export default async function RootLayout({children}: { children: React.ReactNode }) {
    return (
        <div className="App">
            <Providers>
                {children}
            </Providers>
        </div>
    );
}
