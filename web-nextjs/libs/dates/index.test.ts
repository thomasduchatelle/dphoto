import {describe, expect, it} from 'vitest';
import {toLocaleTimeOfDay} from './index';

describe('toLocaleTimeOfDay', () => {
    it('should format an afternoon time using the 24-hour clock', () => {
        expect(toLocaleTimeOfDay(new Date('2025-04-25T14:05:00'))).toBe('14:05');
    });

    it('should format midnight as 00 in the 24-hour clock', () => {
        expect(toLocaleTimeOfDay(new Date('2025-04-25T00:30:00'))).toBe('00:30');
    });
});
