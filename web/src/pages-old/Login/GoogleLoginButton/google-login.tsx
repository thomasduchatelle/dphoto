import React, {memo, useEffect, useRef} from "react";
import {LogoutListener} from "../../../core/security";
import {useConfigContext} from "../../../core/application";
import useWindowDimensions from "../../../core/utils/window-utils";
import {loadScript} from "./loadScript";
import {IdentityProviderError} from "../domain";
import {googleLogout} from "./google-logout";
import {Box} from "@mui/material";

export const GoogleLoginIntegration = memo(function ({onError, onIdentitySuccess}: {
    onError(error: Error): void
    onIdentitySuccess(identityToken: string, logoutListener: LogoutListener): void
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

                    if (window.google && buttonRef.current) {
                        window.google.accounts.id.renderButton(buttonRef.current, {
                            type: 'standard',
                            width: `${buttonSize}px`, // this is not reactive without refresh
                            text: 'continue_with',
                        });
                    }
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
    }, [buttonSize, googleClientId, onError, onIdentitySuccess])

    return <Box ref={buttonRef} id='google-login-prompt' sx={{
        width: `${buttonSize}px`,
        margin: 'auto',
    }}></Box>
})