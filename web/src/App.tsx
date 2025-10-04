import {CssBaseline} from "@mui/material";
import {BrowserRouter} from "react-router-dom";
import DPhotoTheme from "./pages/DPhotoTheme";
import CRARouter from "./pages/_cra-router";
import {ApplicationContextComponent} from "./core/application";
import {AdapterDayjs} from '@mui/x-date-pickers/AdapterDayjs';
import {LocalizationProvider} from "@mui/x-date-pickers";
import dayjs from "dayjs";
import fr from "dayjs/locale/fr";

dayjs.locale(fr)

const App = () => {
    // TODO - add React error boundary
    return (
        <div className="App">
            <DPhotoTheme>
                <CssBaseline/>
                <BrowserRouter>
                    <LocalizationProvider dateAdapter={AdapterDayjs} adapterLocale='fr'>
                        <ApplicationContextComponent>
                            <CRARouter/>
                        </ApplicationContextComponent>
                    </LocalizationProvider>
                </BrowserRouter>
            </DPhotoTheme>
        </div>
    )
}

export default App;
