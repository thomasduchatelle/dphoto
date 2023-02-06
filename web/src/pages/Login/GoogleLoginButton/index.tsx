import React, {memo, useEffect, useRef} from "react";
import {Box} from "@mui/material";
import useWindowDimensions from "../../../core/utils/window-utils";
import {IdentityProviderError} from "../domain";
import {LogoutListener} from "../../../core/security";
import {googleLogout} from "./google-logout";
import {useConfigContext} from "../../../core/application";

export const loadScript = (src: string) => {
    return new Promise<void>((resolve, reject) => {
        if (document.querySelector(`script[src="${src}"]`)) return resolve()

        const script = document.createElement('script')
        script.src = src
        script.onload = () => resolve()
        script.onerror = (err) => reject(err)
        document.body.appendChild(script)
    })
}

const GoogleLoginIntegration = memo(function ({onError, onIdentitySuccess, onWaitingUserInput}: {
    onError(error: Error): void
    onIdentitySuccess(identityToken: string, logoutListener: LogoutListener): void
    onWaitingUserInput(): void
}) {
    const buttonRef = useRef<HTMLDivElement>(null);
    const {googleClientId} = useConfigContext();
    const windowDimension = useWindowDimensions();

    const buttonSize = windowDimension.width <= 400 ? 200 : 400;

    useEffect(() => {
        if (!googleClientId) {
            return;
        }
        loadScript('https://accounts.google.com/gsi/client')
            .then(() => {
                if (typeof window === "undefined" || !window.google || !buttonRef.current) {
                    onError(new IdentityProviderError(`window.google not ready [google=${window.google} ; buttonRef=${buttonRef.current}]`))
                    return
                }
                try {
                    window.google.accounts.id.initialize({
                        auto_select: true,
                        client_id: googleClientId,
                        cancel_on_tap_outside: false,
                        prompt_parent_id: 'google-login-prompt',
                        callback: (res) => {
                            if (res.credential) {
                                onIdentitySuccess(res.credential, {
                                    onLogout: () => {
                                        googleLogout().catch(err => {
                                            console.log(`WARN: failed to logout: ${err}`)
                                        })
                                    }
                                })
                            } else {
                                onError(new IdentityProviderError(`no credentials in Google response ${JSON.stringify(res)}`))
                            }
                        },
                    });
                    window.google.accounts.id.prompt(res => {
                        if (res.isDismissedMoment() && res.getDismissedReason() !== "credential_returned") {
                            // do nothing - callback handler will do the rest

                        } else if (res.isDisplayMoment() && res.isDisplayed()) {
                            onWaitingUserInput(); // loading signal and messages are built-in

                        } else if (res.isNotDisplayed() || res.isSkippedMoment()) {
                            // fallback on native button
                            onWaitingUserInput();
                            if (window.google && buttonRef.current) {
                                window.google.accounts.id.renderButton(buttonRef.current, {
                                    type: 'standard',
                                    width: `${buttonSize}px`, // this is not reactive without refresh
                                    text: 'continue_with',
                                });
                            }
                        }
                    });
                } catch (error) {
                    onError(new IdentityProviderError("Google Login button cannot be generated"))
                }
            })

        const currentButtonRef = buttonRef.current
        return () => {
            if (currentButtonRef) {
                currentButtonRef.innerText = ''
            }
        }
    }, [buttonSize, googleClientId, onError, onWaitingUserInput, onIdentitySuccess])

    return <Box ref={buttonRef} id='google-login-prompt' sx={{
        width: `${buttonSize}px`,
        margin: 'auto',
    }}></Box>
})

export default GoogleLoginIntegration
