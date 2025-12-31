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
        const distributionId = process.env.SST_DISTRIBUTION_ID;

        console.log(`SST_DISTRIBUTION_ID=${distributionId}`)
        console.log(`SST_COGNITO_ISSUER=${process.env.SST_COGNITO_ISSUER}`)
        console.log(`SST_COGNITO_CLIENT_ID=${process.env.SST_COGNITO_CLIENT_ID}`)
        console.log(`SST_COGNITO_CLIENT_SECRET=${process.env.SST_COGNITO_CLIENT_SECRET}`)

        const distribution = sst.aws.Router.get("CDN", distributionId)
        const dphotoWeb = new sst.aws.Nextjs("DPhotoWEB", {
            // domain: "nextjs.next.duchatelle.me",
            buildCommand: "exit 0", // skip the build (done is previous workflow step)
            server: {
                memory: "512 MB",
            },
            router: {
                instance: distribution,
            },
            environment: {
                NEXT_PUBLIC_COGNITO_ISSUER: process.env.SST_COGNITO_ISSUER!,
                NEXT_PUBLIC_COGNITO_CLIENT_ID: process.env.SST_COGNITO_CLIENT_ID!,
                COGNITO_CLIENT_SECRET: process.env.SST_COGNITO_CLIENT_SECRET!,
            }
        });
    },
});
