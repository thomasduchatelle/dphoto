// This file must be loaded BEFORE setupTests.ts to polyfill localStorage for MSW
if (typeof localStorage === 'undefined') {
    class LocalStorageMock {
        private store: Map<string, string> = new Map();

        getItem(key: string): string | null {
            return this.store.get(key) ?? null;
        }

        setItem(key: string, value: string): void {
            this.store.set(key, value);
        }

        removeItem(key: string): void {
            this.store.delete(key);
        }

        clear(): void {
            this.store.clear();
        }

        key(index: number): string | null {
            const keys = Array.from(this.store.keys());
            return keys[index] ?? null;
        }

        get length(): number {
            return this.store.size;
        }
    }

    global.localStorage = new LocalStorageMock() as Storage;
}

