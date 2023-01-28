export class AccessForbiddenError extends Error {
    constructor(private readonly details: string) {
        super("Error: you're not allowed to access this page.");
    }
}