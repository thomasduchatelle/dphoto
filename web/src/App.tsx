import {CssBaseline} from "@mui/material";
import {BrowserRouter} from "react-router-dom";
import DPhotoTheme from "./components/DPhotoTheme";
import GeneralRouter from "./pages/GeneralRouter";
import {ApplicationContextComponent} from "./core/application";
import {CatalogContextComponent} from "./core/catalog";

const App = () => {
    // TODO - add React error boundary
    return (
        <div className="App">
            <DPhotoTheme>
                <CssBaseline/>
                <BrowserRouter>
                    <ApplicationContextComponent>
                        <CatalogContextComponent>
                            <GeneralRouter/>
                        </CatalogContextComponent>
                    </ApplicationContextComponent>
                </BrowserRouter>
            </DPhotoTheme>
        </div>
    )
}

export default App;
