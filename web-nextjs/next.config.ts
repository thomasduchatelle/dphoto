import type {NextConfig} from "next";

const nextConfig: NextConfig = {
    output: "standalone",
    images: {
        deviceSizes: [640, 750, 828, 1080, 1200, 1440, 1920, 2048, 2400, 3840],
        remotePatterns: [
            {
                protocol: 'https',
                hostname: '**',
            }
        ]
    },
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
