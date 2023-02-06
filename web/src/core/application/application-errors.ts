import {ErrorWithPublicMessage} from "./application-model";

export class AccessForbiddenError extends ErrorWithPublicMessage {
    public readonly publicMessage: string = "You're not allowed to access this page."
}

export class InternalError extends ErrorWithPublicMessage {
    public readonly publicMessage: string = "An error occurred, please report it to the maintainer."
}
