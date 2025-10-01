import {defineConfig} from 'waku/config';
import tailwindcss from '@tailwindcss/vite';

export default defineConfig({
    basePath: '/waku',
    vite: {
        plugins: [tailwindcss()],
        server: {
            port: 3001, // Use different port to avoid conflict with CRA dev server (port 3000)
        },
    },
});
