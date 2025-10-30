import * as fs from 'fs';
import * as path from 'path';
import * as os from 'os';
import {computeHash, FilePattern, GO_FILES} from './hash-utils';

describe('hash-utils', () => {
    let tempDir: string;

    beforeEach(() => {
        tempDir = fs.mkdtempSync(path.join(os.tmpdir(), 'hash-utils-test-'));
    });

    afterEach(() => {
        fs.rmSync(tempDir, {recursive: true, force: true});
    });

    const createFile = (relativePath: string, content: string) => {
        const fullPath = path.join(tempDir, relativePath);
        const dir = path.dirname(fullPath);
        if (!fs.existsSync(dir)) {
            fs.mkdirSync(dir, {recursive: true});
        }
        fs.writeFileSync(fullPath, content, 'utf-8');
    };

    describe('GO_FILES pattern', () => {
        it('should match .go files', () => {
            expect(GO_FILES.test('main.go')).toBe(true);
            expect(GO_FILES.test('handler.go')).toBe(true);
        });

        it('should match go.mod', () => {
            expect(GO_FILES.test('go.mod')).toBe(true);
        });

        it('should match go.sum', () => {
            expect(GO_FILES.test('go.sum')).toBe(true);
        });

        it('should not match non-Go files', () => {
            expect(GO_FILES.test('main.ts')).toBe(false);
            expect(GO_FILES.test('README.md')).toBe(false);
            expect(GO_FILES.test('package.json')).toBe(false);
        });
    });

    describe('computeHash', () => {
        it('should compute hash of Go files in a directory', () => {
            createFile('main.go', 'package main\n\nfunc main() {}\n');
            createFile('go.mod', 'module example.com/test\n\ngo 1.23\n');
            createFile('README.md', 'This is a readme');

            const hash = computeHash({
                directories: [tempDir],
                filePatterns: [GO_FILES]
            });

            expect(hash).toMatch(/^[a-f0-9]{64}$/);
        });

        it('should compute hash recursively', () => {
            createFile('main.go', 'package main');
            createFile('pkg/utils/helper.go', 'package utils');
            createFile('pkg/other/other.go', 'package other');

            const hash = computeHash({
                directories: [tempDir],
                filePatterns: [GO_FILES]
            });

            expect(hash).toMatch(/^[a-f0-9]{64}$/);
        });

        it('should produce different hashes for different content', () => {
            createFile('main.go', 'package main\n\nfunc main() {}\n');

            const hash1 = computeHash({
                directories: [tempDir],
                filePatterns: [GO_FILES]
            });

            createFile('main.go', 'package main\n\nfunc main() {\n  println("hello")\n}\n');

            const hash2 = computeHash({
                directories: [tempDir],
                filePatterns: [GO_FILES]
            });

            expect(hash1).not.toBe(hash2);
        });

        it('should produce the same hash for identical content', () => {
            createFile('main.go', 'package main\n\nfunc main() {}\n');
            createFile('go.mod', 'module test\n\ngo 1.23\n');

            const hash1 = computeHash({
                directories: [tempDir],
                filePatterns: [GO_FILES]
            });

            const hash2 = computeHash({
                directories: [tempDir],
                filePatterns: [GO_FILES]
            });

            expect(hash1).toBe(hash2);
        });

        it('should handle multiple directories', () => {
            const tempDir2 = fs.mkdtempSync(path.join(os.tmpdir(), 'hash-utils-test-2-'));

            try {
                createFile('main.go', 'package main');
                fs.writeFileSync(path.join(tempDir2, 'other.go'), 'package other', 'utf-8');

                const hash = computeHash({
                    directories: [tempDir, tempDir2],
                    filePatterns: [GO_FILES]
                });

                expect(hash).toMatch(/^[a-f0-9]{64}$/);
            } finally {
                fs.rmSync(tempDir2, {recursive: true, force: true});
            }
        });

        it('should filter files by custom pattern', () => {
            createFile('file.txt', 'text file');
            createFile('file.md', 'markdown file');
            createFile('main.go', 'package main');

            const txtPattern: FilePattern = {
                test: (filename: string) => filename.endsWith('.txt')
            };

            const hash = computeHash({
                directories: [tempDir],
                filePatterns: [txtPattern]
            });

            expect(hash).toMatch(/^[a-f0-9]{64}$/);
        });

        it('should handle multiple file patterns', () => {
            createFile('file.txt', 'text file');
            createFile('file.md', 'markdown file');
            createFile('main.go', 'package main');

            const txtPattern: FilePattern = {
                test: (filename: string) => filename.endsWith('.txt')
            };
            const mdPattern: FilePattern = {
                test: (filename: string) => filename.endsWith('.md')
            };

            const hash = computeHash({
                directories: [tempDir],
                filePatterns: [txtPattern, mdPattern]
            });

            expect(hash).toMatch(/^[a-f0-9]{64}$/);
        });

        it('should handle non-existent directories gracefully', () => {
            const hash = computeHash({
                directories: ['/non/existent/path'],
                filePatterns: [GO_FILES]
            });

            expect(hash).toMatch(/^[a-f0-9]{64}$/);
        });

        it('should return consistent hash for empty directories', () => {
            const hash1 = computeHash({
                directories: [tempDir],
                filePatterns: [GO_FILES]
            });

            const hash2 = computeHash({
                directories: [tempDir],
                filePatterns: [GO_FILES]
            });

            expect(hash1).toBe(hash2);
        });

        it('should handle individual files in directories array', () => {
            createFile('main.go', 'package main');
            createFile('go.mod', 'module test\n\ngo 1.23\n');
            createFile('pkg/utils/helper.go', 'package utils');

            const hash = computeHash({
                directories: [
                    path.join(tempDir, 'go.mod'),
                    path.join(tempDir, 'pkg')
                ],
                filePatterns: [GO_FILES]
            });

            expect(hash).toMatch(/^[a-f0-9]{64}$/);
        });
    });
});

