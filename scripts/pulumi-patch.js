// Patch for @pulumi/pulumi/cmd/run/error.js
// Run with `node scripts/pulumi-patch.js` once SST is installed in order to have usable error messages.
const fs = require('fs');
const path = require('path');

const errorFile = path.join(
    __dirname,
    '../web-nextjs/.sst/platform/node_modules/@pulumi/pulumi/cmd/run/error.js'
);

const needle = "return util.inspect(err, { colors: true });";
const replacement = `
try {
    return util.inspect(err, { colors: true });
} catch (inspectError) {
    const max = 20000;
    const truncate = (s) => (typeof s === "string" && s.length > max
        ? s.slice(0, max) + \`\\n... [truncated \${s.length - max} chars]\`
        : s);
    const msg = typeof err?.message === "string" ? err.message : undefined;
    const stack = typeof err?.stack === "string" ? err.stack : undefined;
    return truncate(msg) || truncate(stack) || \`error inspection failed: \${inspectError.message}\`;
}`;

const content = fs.readFileSync(errorFile, 'utf-8');
if (!content.includes('inspectError')) {
    fs.writeFileSync(errorFile, content.replace(needle, replacement));
    console.log('Patched Pulumi error handler');
}