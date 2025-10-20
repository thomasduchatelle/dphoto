'use client';

import {CssBaseline} from "@mui/material";
import DPhotoTheme from "./DPhotoTheme";
import {ApplicationContextComponent} from "../core/application";
import {AdapterDayjs} from '@mui/x-date-pickers/AdapterDayjs';
import {LocalizationProvider} from "@mui/x-date-pickers";
import {RouterProvider} from "./ClientRouter";
import dayjs from "dayjs";
import fr from "dayjs/locale/fr";
import {ReactNode} from "react";

dayjs.locale(fr)

export const Providers = ({children}: {children: ReactNode}) => {
    return (
        <DPhotoTheme>
            <CssBaseline/>
            <LocalizationProvider dateAdapter={AdapterDayjs} adapterLocale='fr'>
                <ApplicationContextComponent>
                    <RouterProvider>
                        {children}
                    </RouterProvider>
                </ApplicationContextComponent>
            </LocalizationProvider>
        </DPhotoTheme>
    );
};
