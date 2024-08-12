import type {Config} from 'jest';

export default {
    setupFiles: ['./jest.polyfills.js'],
    // presets: [
    //     ['@babel/preset-env', {targets: {node: 'current'}}],
    //     '@babel/preset-typescript',
    // ],
    transformIgnorePatterns: [
        '/node_modules/(?!(@bundled-es-modules)/)',
    ],
} as Config