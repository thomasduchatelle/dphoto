import {loadScript} from "./index";

export const googleLogout = (email?: string): Promise<void> => {
    return loadScript('https://accounts.google.com/gsi/client')
        .then(() => {
            if (typeof window === "undefined" || !window.google) {
                return
            }

            window.google.accounts.id.disableAutoSelect();

            if (email) {
                return new Promise<void>((resolve, reject) => {
                    if (window.google) {
                        window.google.accounts.id.revoke(email, _ => resolve());
                    }
                })
            }
        })
}
