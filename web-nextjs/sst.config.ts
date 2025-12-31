// eslint-disable-next-line @typescript-eslint/triple-slash-reference
/// <reference path="./.sst/platform/config.d.ts" />

export default $config({
    app(input) {
        return {
            name: "web-nextjs",
            removal: "remove",
            protect: false,
            home: "aws",
        };
    },
    async run() {
        const dphotoWeb = new sst.aws.Nextjs("DPhotoWEB", {
            domain: "nextjs.next.duchatelle.me",
            buildCommand: "exit 0", // skip the build (done is previous workflow step)
            server: {
                memory: "512 MB",
            },
        });
    },
});
