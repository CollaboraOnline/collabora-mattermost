const exec = require('child_process').exec;

const path = require('path');

const PLUGIN_ID = require('../plugin.json').id;

const NPM_TARGET = process.env.npm_lifecycle_event; //eslint-disable-line no-process-env
const isDev = NPM_TARGET === 'debug' || NPM_TARGET === 'debug:watch';

const STANDARD_EXCLUDE = [path.join(__dirname, 'node_modules')];

const plugins = [];
if (NPM_TARGET === 'build:watch' || NPM_TARGET === 'debug:watch') {
  plugins.push({
    apply: (compiler) => {
      compiler.hooks.watchRun.tap('WatchStartPlugin', () => {
        // eslint-disable-next-line no-console
        console.log('Change detected. Rebuilding webapp.');
      });
      compiler.hooks.afterEmit.tap('AfterEmitPlugin', () => {
        exec('cd .. && make deploy-from-watch', (err, stdout, stderr) => {
          if (stdout) {
            process.stdout.write(stdout);
          }
          if (stderr) {
            process.stderr.write(stderr);
          }
        });
      });
    },
  });
}

const config = {
  entry: ['./src/index.tsx'],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, 'src'),
      'mattermost-redux': path.resolve(__dirname, 'node_modules/@mattermost/webapp/packages/mattermost-redux/src'),
      '@mattermost/client': path.resolve(__dirname, 'node_modules/@mattermost/webapp/packages/client/src'),
      '@mattermost/types': path.resolve(__dirname, 'node_modules/@mattermost/types/lib'),
      reselect: path.resolve(__dirname, 'node_modules/@mattermost/webapp/packages/reselect/src'),
    },
    modules: ['node_modules', path.resolve(__dirname), 'src'],
    extensions: ['.js', '.jsx', '.ts', '.tsx'],
  },
  module: {
    rules: [
      {
        test: /\.(js|jsx|ts|tsx)$/,
        exclude: /node_modules\/(?!(@mattermost\/webapp)\/).*/,
        use: {
          loader: 'babel-loader',
          options: {
            cacheDirectory: true,

            // Babel configuration is in babel.config.js because jest requires it to be there.
          },
        },
      },
      {
        test: /\.(png|eot|tiff|svg|woff2|woff|ttf|gif|mp3|jpg)$/,
        type: 'asset/inline', // consider 'asset' when URL resource chunks are supported
      },
      {
        test: /\.(scss|css)$/,
        use: [
          'style-loader',
          {
            loader: 'css-loader',
          },
          {
            loader: 'sass-loader',
            options: {
              sassOptions: {
                includePaths: ['node_modules/compass-mixins/lib', 'sass'],
              },
            },
          },
        ],
      },
    ],
  },
  externals: {
    react: 'React',
    'react-dom': 'ReactDOM',
    redux: 'Redux',
    'react-redux': 'ReactRedux',
    'post-utils': 'PostUtils',
    'prop-types': 'PropTypes',
    'react-bootstrap': 'ReactBootstrap',
    'react-router-dom': 'ReactRouterDom',
  },
  output: {
    devtoolNamespace: PLUGIN_ID,
    path: path.join(__dirname, '/dist'),
    publicPath: '/',
    filename: 'main.js',
  },
  mode: isDev ? 'eval-source-map' : 'production',
  plugins,
};

if (isDev) {
  Object.assign(config, { devtool: 'eval-source-map' });
}

module.exports = config;
