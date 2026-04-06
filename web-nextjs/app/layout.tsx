import "server-only"

import type {Metadata} from "next";
import {ThemeProvider} from "@/components/theme";
import {AppRouterCacheProvider} from '@mui/material-nextjs/v15-appRouter';
import {NavigationLoadingIndicator} from '@/components/NavigationLoadingIndicator';

export const metadata: Metadata = {
    title: "DPhoto",
    description: "Photo management application",
};

export default async function RootLayout({
                                             children,
                                         }: Readonly<{
    children: React.ReactNode;
}>) {

    return (
        <html lang="en">
        <body>
        <NavigationLoadingIndicator/>
        <AppRouterCacheProvider>
            <ThemeProvider>
                {children}
            </ThemeProvider>
        </AppRouterCacheProvider>
        </body>
        </html>
    );
}
