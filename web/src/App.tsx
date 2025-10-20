'use client';

import GeneralRouter from "./pages-old/GeneralRouter";
import {Providers} from "./components/Providers";

const App = () => {
    // TODO - add React error boundary
    return (
        <div className="App">
            <Providers>
                <GeneralRouter/>
            </Providers>
        </div>
    )
}

export default App;
