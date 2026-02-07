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
    async rewrites() {
        if (process.env.NODE_ENV === 'development') {
            return [
                {
                    source: '/api/:path*',
                    destination: 'http://127.0.0.1:8080/api/:path*',
                },
            ];
        }

        return [];
    },
};

export default nextConfig;
