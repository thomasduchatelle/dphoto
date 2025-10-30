import * as fs from 'fs';
import * as path from 'path';
import * as crypto from 'crypto';

export interface FilePattern {
    test: (filename: string) => boolean;
}

export const GO_FILES: FilePattern = {
    test: (filename: string) => {
        return filename.endsWith('.go') ||
            filename === 'go.mod' ||
            filename === 'go.sum';
    }
};

export interface ComputeHashOptions {
    directories: string[];
    filePatterns: FilePattern[];
}

export function computeHash(options: ComputeHashOptions): string {
    const hash = crypto.createHash('sha256');
    const filesToHash: string[] = [];

    const collectFiles = (dir: string) => {
        if (!fs.existsSync(dir)) {
            return;
        }
        const stats = fs.statSync(dir);

        if (stats.isFile()) {
            const filename = path.basename(dir);
            const matchesPattern = options.filePatterns.some(pattern =>
                pattern.test(filename)
            );
            if (matchesPattern) {
                filesToHash.push(dir);
            }
            return;
        }

        if (stats.isDirectory()) {
            const entries = fs.readdirSync(dir, {withFileTypes: true});
            for (const entry of entries) {
                const fullPath = path.join(dir, entry.name);
                if (entry.isDirectory()) {
                    collectFiles(fullPath);
                } else if (entry.isFile()) {
                    const matchesPattern = options.filePatterns.some(pattern =>
                        pattern.test(entry.name)
                    );
                    if (matchesPattern) {
                        filesToHash.push(fullPath);
                    }
                }
            }
        }
    };

    for (const directory of options.directories) {
        collectFiles(directory);
    }

    filesToHash.sort();

    for (const file of filesToHash) {
        if (fs.existsSync(file)) {
            const content = fs.readFileSync(file);
            hash.update(file);
            hash.update(new Uint8Array(content));
        }
    }

    return hash.digest('hex');
}

