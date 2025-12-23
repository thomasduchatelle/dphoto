import DPhotoTheme from "./DPhotoTheme";
import {RouterProvider} from "./ClientRouter";
import {ReactNode} from "react";
import {ErrorBoundary} from "./ErrorBoundary";
import JotaiProvider from "./JotaiProvider";
import {ApplicationContextComponent} from "../core/application";


export const Providers = ({children}: { children: ReactNode }) => {
    return (
        <JotaiProvider>
            <ApplicationContextComponent>
                <DPhotoTheme>
                    <RouterProvider>
                        <ErrorBoundary> {/* Error Boundaries is using AppNav which requires RouterProvider */}
                            {children}
                        </ErrorBoundary>
                    </RouterProvider>
                </DPhotoTheme>
            </ApplicationContextComponent>
        </JotaiProvider>
    );
};
