import {defineConfig} from 'waku/config';
import path from 'path';
import {fileURLToPath} from 'url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));

export default defineConfig({
    vite: {
        resolve: {
            alias: {
                'src': path.resolve(__dirname, 'src'),
            },
        },
        server: {
            port: 3000,
            proxy: {
                '/oauth': {
                    target: 'http://127.0.0.1:8080',
                    changeOrigin: true,
                },
                '/api': {
                    target: 'http://127.0.0.1:8080',
                    changeOrigin: true,
                },
            },
        },
    },
    unstable_honoEnhancer: './src/hono-enhancer',
    middleware: [
        'waku/middleware/context',
        './src/middleware/cookie.js',
        // './src/middleware/noop.js',
        'waku/middleware/dev-server',
        'waku/middleware/handler',
    ],
});
