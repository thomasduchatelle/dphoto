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
import {ErrorBoundary} from "./ErrorBoundary";
import {Provider} from 'jotai'

dayjs.locale(fr)

export const Providers = ({children}: { children: ReactNode }) => {
    return (
        <ErrorBoundary>
            <Provider>
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
            </Provider>
        </ErrorBoundary>
    );
};
