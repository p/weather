var path = require('path');
var webpack = require('webpack');
  const UglifyJSPlugin = require('uglifyjs-webpack-plugin');

module.exports = function(env){
  return {
    mode: 'production',
    devtool: 'source-map',
  entry: [
    '@babel/polyfill',
    'whatwg-fetch',
    './src/index'
  ],
  output: {
    path: path.join(__dirname, 'build'),
    filename: 'bundle.js',
    publicPath: '/static/'
  },
  plugins: [
    new webpack.NamedModulesPlugin(),
  new UglifyJSPlugin({sourceMap:true}),
  new webpack.DefinePlugin({
    API_URL:env.API_URL,
    NODE_ENV:'"production"',
  }),
  ],
  module: {
    rules: [{
      test: /\.js$/,
      use: ['babel-loader'],
      include: path.join(__dirname, 'src')
    },
    ]
  }
}
};
