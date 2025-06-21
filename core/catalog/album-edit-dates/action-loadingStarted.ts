export interface LoadingStarted {
    type: "LoadingStarted";
}

export function loadingStarted(): LoadingStarted {
    return {
        type: "LoadingStarted",
    };
}

export function reduceLoadingStarted(
    current: CatalogViewerState,
    action: LoadingStarted,
): CatalogViewerState {
    if (!current.editDatesDialog) {
        return current;
    }

    return {
        ...current,
        editDatesDialog: {
            ...current.editDatesDialog,
            isLoading: true,
        },
    };
}

export function loadingStartedReducerRegistration(handlers: any) {
    handlers["LoadingStarted"] = reduceLoadingStarted as (
        state: CatalogViewerState,
        action: LoadingStarted
    ) => CatalogViewerState;
}
