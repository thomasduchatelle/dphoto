/**
 * Ports might throw a CatalogError. Usage example:
 *
 *      try {
 *          // do something
 *      } catch (err) {
 *          if (isCatalogError(err) && err.code === GoodReasonErrorCode) {
 *              dispatch(serverDidNotProcessTheRequest(err.code, err.message));
 *          } else {
 *              dispatch(somethingWhenWrong(getErrorMessage(err) ?? "Something went wrong."));
 *          }
 *      }
 */
export class CatalogError extends Error {
    constructor(public readonly code: string, message: string) {
        super(message);
        this.name = "ApplicationAPIError";
    }
}

export function isCatalogError(err: any): err is CatalogError {
    return err && err.name === "ApplicationAPIError" && typeof err.code === "string" && err.code;
}

export function getErrorMessage(err: any): string | undefined {
    if (err && typeof err.message === "string") {
        return err.message;
    }
}