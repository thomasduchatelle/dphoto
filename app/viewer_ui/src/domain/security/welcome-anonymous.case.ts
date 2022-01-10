import {SecurityDependencies} from "./security.domain";

export function welcomeAnonymous(requestedPath: string): Promise<void> {
  if (!SecurityDependencies.navigationManager) {
    return Promise.reject("'navigationManager' has not been initialised.")
  }

  if (!requestedPath.startsWith("/login")) {
    const params = requestedPath && requestedPath != "/" ? `?redirect=${encodeURIComponent(requestedPath)}` : ""
    SecurityDependencies.navigationManager.gotoLoginPage(`/login${params}`);
  }
  return Promise.resolve()
}