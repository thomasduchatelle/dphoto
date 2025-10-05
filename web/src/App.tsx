'use client';

import {CssBaseline} from "@mui/material";
import {BrowserRouter} from "react-router-dom";
import DPhotoTheme from "./components/DPhotoTheme";
import GeneralRouter from "./pages-old/GeneralRouter";
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
                            <GeneralRouter/>
                        </ApplicationContextComponent>
                    </LocalizationProvider>
                </BrowserRouter>
            </DPhotoTheme>
        </div>
    )
}

export default App;
