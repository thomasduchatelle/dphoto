import {ApplicationContextType} from "./application-context";
import {InternalError} from "./application-errors";
import {ErrorWithPublicMessage} from "./application-model";

export type UnrecoverableErrorAction = {
    type: 'unrecoverable-error'
    error: Error
}

export type ApplicationGenericAction = UnrecoverableErrorAction

export function applicationGenericReducer(current: ApplicationContextType, action: ApplicationGenericAction): ApplicationContextType {
    switch (action.type) {
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