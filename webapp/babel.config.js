require('regenerator-runtime/runtime');

const config = {
  presets: [
    [
      '@babel/preset-env',
      {
        // Set the corejs version we are using to avoid warnings in console.
        // and polyfill every proposal supported by core-js@3
        corejs: 3,
        // Uses browserslist from package.json to select polyfills
        // Adds specific imports for polyfills when they are used in each file.
        // We take advantage of the fact that a bundler will load the same polyfill only once.
        useBuiltIns: 'usage',
        // This will enable polyfills and transforms for proposal which have already been shipped in browsers for a while.
        shippedProposals: true,
        // Setting this to false will preserve ES modules, to ship native ES Modules to browsers.
        modules: false,
        // Set to true to output the polyfills and transform plugins enabled by preset-env to console.log
        // and, if applicable, which one of your targets that needed it.
        debug: false,
      },
    ],
    [
      '@babel/preset-react',
      {
        // Will use the native built-in instead of trying to polyfill
        // behavior for any plugins that require one.
        useBuiltIns: true,
      },
    ],
    [
      '@babel/typescript',
      {
        // Enable jsx parsing for angle brackets instead of typescript's legacy type assertion.
        allExtensions: true,
        isTSX: true,
      },
    ],
  ],
  plugins: ['babel-plugin-typescript-to-proptypes'],
};

// Jest needs module transformation
config.env = {
  test: {
    presets: config.presets,
    plugins: config.plugins,
  },
};
config.env.test.presets[0][1].modules = 'auto';

module.exports = config;
