'use client';

import {CssBaseline} from "@mui/material";
import DPhotoTheme from "./DPhotoTheme";
import {AdapterDayjs} from '@mui/x-date-pickers/AdapterDayjs';
import {LocalizationProvider} from "@mui/x-date-pickers";
import {RouterProvider} from "./ClientRouter";
import dayjs from "dayjs";
import fr from "dayjs/locale/fr";
import {ReactNode} from "react";
import {ErrorBoundary} from "./ErrorBoundary";

dayjs.locale(fr)

export const Providers = ({children}: { children: ReactNode }) => {
    return (
        <ErrorBoundary>
            <DPhotoTheme>
                <CssBaseline/>
                <LocalizationProvider dateAdapter={AdapterDayjs} adapterLocale='fr'>
                    <RouterProvider>
                        {children}
                    </RouterProvider>
                </LocalizationProvider>
            </DPhotoTheme>
        </ErrorBoundary>
    );
};
