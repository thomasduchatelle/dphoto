const CracoAlias = require('craco-alias');

module.exports = {
  plugins: [
    {
      plugin: CracoAlias,
      options: {
        source: 'tsconfig',
        // baseUrl should be the same as the baseUrl in your tsconfig.json
        // If your tsconfig.json is in 'web/' relative to your project root,
        // and craco.config.js is in the project root, then baseUrl should be 'web'.
        baseUrl: './web', // This assumes your tsconfig.json is located at ./web/tsconfig.json
        tsConfigPath: './web/tsconfig.json', // Explicitly point to your tsconfig.json
      },
    },
  ],
};
