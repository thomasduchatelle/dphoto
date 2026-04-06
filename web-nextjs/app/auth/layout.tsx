import "server-only"

import type {Metadata} from "next";
import {AppBackground} from "@/components/AppLayout/AppBackground";

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
        <AppBackground>
            {children}
        </AppBackground>
    );
}
