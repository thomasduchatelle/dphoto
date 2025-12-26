import {defineConfig} from 'vitest/config';
import react from '@vitejs/plugin-react';
import path from 'path';

export default defineConfig({
  plugins: [react()],
  test: {
    globals: true,
    environment: 'happy-dom', // note - JSDOM has an issue with https://github.com/vitest-dev/vitest/pull/1727 making middleware tests failing.
    setupFiles: ['./src/setupTests.ts'],
    include: [
      'src/**/?(*.)+(spec|test).+(ts|tsx|js)',
    ],
    exclude: ['**/node_modules/**', '**/dist/**', '**/*.disabled.*'],
    coverage: {
      provider: 'v8',
      include: ['src/**/*.{ts,tsx}'],
      exclude: ['src/**/*.d.ts', 'src/**/*.stories.tsx'],
    },
  },
  resolve: {
    alias: {
      src: path.resolve(__dirname, './src'),
    },
  },
});
