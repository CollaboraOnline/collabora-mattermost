require('core-js/stable');
require('regenerator-runtime/runtime');

const config = {
    presets: [
        ['@babel/preset-env', {
            targets: {
                chrome: 66,
                firefox: 60,
                edge: 42,
                safari: 12,
            },
            modules: false,
            corejs: 3,
            debug: false,
            useBuiltIns: 'usage',
            shippedProposals: true,
        }],
        ['@babel/preset-react', {
            useBuiltIns: true,
        }],
        ['@babel/typescript', {
            allExtensions: true,
            isTSX: true,
        }],
        ['@emotion/babel-preset-css-prop'],
    ],
    plugins: [
        '@babel/plugin-syntax-class-properties',
        '@babel/plugin-syntax-dynamic-import',
        '@babel/plugin-syntax-object-rest-spread',
        '@babel/plugin-syntax-optional-chaining',
        'babel-plugin-typescript-to-proptypes',
    ],
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
