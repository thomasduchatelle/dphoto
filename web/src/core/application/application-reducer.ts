import {ApplicationContextType} from "./application-context";
import {InternalError} from "./application-errors";
import {ErrorWithPublicMessage} from "./application-model";

export type ConfigLoadedAction = {
    type: 'config-loaded'
    googleClientId: string
}

export type UnrecoverableErrorAction = {
    type: 'unrecoverable-error'
    error: Error
}

export type ApplicationGenericAction = ConfigLoadedAction | UnrecoverableErrorAction

export function applicationGenericReducer(current: ApplicationContextType, action: ApplicationGenericAction): ApplicationContextType {
    switch (action.type) {
        case 'config-loaded':
            return {
                ...current,
                general: {
                    ...current.general,
                    googleClientId: action.googleClientId,
                }
            }

        case "unrecoverable-error":
            let error = action.error as any as ErrorWithPublicMessage
            if (!error || !error.publicMessage) {
                error = new InternalError(action.error.message, action.error)
            }

            return {
                ...current,
                general: {
                    ...current.general,
                    error: error,
                }
            }
    }

    return current
}