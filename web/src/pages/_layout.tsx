import {Providers} from '../components/Providers';
import GeneralRouter from '../pages-old/GeneralRouter';
import {Provider} from "jotai";
import {AuthProvider} from "../components/AuthProvider";

export default async function RootLayout({children}: { children: React.ReactNode }) {
    return (
        <div className="App">
            <Provider>
                <AuthProvider>
                    <Providers>
                        <GeneralRouter/>
                    </Providers>
                </AuthProvider>
            </Provider>
        </div>
    );
}
