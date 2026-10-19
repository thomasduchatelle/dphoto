export class SessionExpiredException extends Error {
    constructor(message = "Session expired") {
        super(message);
        this.name = "SessionExpiredException";
    }
}

export function isSessionExpiredException(err: unknown): err is SessionExpiredException {
    return err instanceof Error && err.name === "SessionExpiredException";
}
