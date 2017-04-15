const path = require('path');
const webpack = require('webpack');

module.exports = {

  context: path.resolve(__dirname, './src'),

  entry: {
    app: './app.js',
  },

  output: {
    path: path.resolve(__dirname, './dist'),
    filename: 'weblabel.bundle.js',
    publicPath: '/',
  },

  devServer: {
    contentBase: path.resolve(__dirname, './src'),
  },

  module: {
    loaders: [
      { test: /.js?$/, loader: 'babel-loader', exclude: /node_modules/, query: { presets: ['es2015'] } },
      { test: /\.css$/, loader: "style-loader!css-loader" },
      { test: /\.png$/, loader: "url-loader?limit=100000" },
      { test: /\.jpg$/, loader: "file-loader" },
      { test: /\.(woff|woff2)(\?v=\d+\.\d+\.\d+)?$/, loader: 'url?limit=10000&mimetype=application/font-woff' },
      { test: /\.ttf(\?v=\d+\.\d+\.\d+)?$/, loader: 'url?limit=10000&mimetype=application/octet-stream' },
      { test: /\.eot(\?v=\d+\.\d+\.\d+)?$/, loader: 'file' },
      { test: /\.svg(\?v=\d+\.\d+\.\d+)?$/, loader: 'url?limit=10000&mimetype=image/svg+xml' },
      { test: /\.styl$/, loader: 'style-loader!css-loader!stylus-loader?paths=node_modules/bootstrap-stylus/stylus/' },
    ]
  },

  plugins: [
    new webpack.ProvidePlugin({
      $: "jquery",
    }),
  ],

  cache: true,
};
