import type {NextConfig} from "next";

const nextConfig: NextConfig = {
    output: "standalone",
    images: {
        remotePatterns: [
            {
                protocol: 'https',
                hostname: '**',
            }
        ]
    },
    basePath: '/nextjs', // when removed, each image src must be updated.
};

export default nextConfig;
