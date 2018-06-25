/*
Copyright (c) 2018 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the 'License');
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

const fs = require('fs')
const path = require('path')
const webpack = require('webpack')

const CopyWebpackPlugin = require('copy-webpack-plugin')

const modDir = path.resolve(__dirname, 'node_modules')
const srcDir = path.resolve(__dirname, 'src')
const rscDir = path.resolve(__dirname, 'src')
const outDir = path.resolve(__dirname, 'build')

module.exports = {
  entry: {
    prolog: path.resolve(srcDir, 'prolog.js'),
    main: path.resolve(srcDir, 'main.js'),
  },

  output: {
    path: outDir,
    filename: 'bundle.js'
  },

  module: {
    rules: [
      {
        test: /\.js$/,
        exclude: [
          /node_modules/
        ],
        use: {
          loader: 'babel-loader',
          options: {
            'presets': [
              path.join(__dirname, 'node_modules/babel-preset-react'),
              path.join(__dirname, 'node_modules/babel-preset-env')
            ],
            'plugins': [
              path.join(__dirname, 'node_modules/babel-plugin-transform-class-properties'),
              path.join(__dirname, 'node_modules/babel-plugin-transform-object-rest-spread'),
              path.join(__dirname, 'node_modules/babel-plugin-transform-object-assign'),
            ]
          }
        }
      },
      {
        test: /\.css$/,
        use: [
          'style-loader',
          'css-loader'
        ]
      },
      {
        test: /\.(eot|ttf|woff|woff2)$/,
        loader: 'file-loader',
        options: {
          name: 'fonts/[name].[ext]'
        }
      },
      {
        test: /\.(gif|jpg|png|svg)$/,
        loader: 'url-loader',
        options: {
          name: 'images/[name].[ext]'
        }
      }
    ]
  },

  plugins: [
	// Some dependencies need to be defined before anything else. For example,
	// jQuery needs to be defined before loading the Bootstrap or Patternfly
	// scripts, as they assume that it is already available. But Webpack doesn't
	// guarantee the order of modules in the generated bundle, so to force that
	// order we generate this separate chunk, and we make sure that it is loaded
	// first in the HTML page.
    new webpack.optimize.CommonsChunkPlugin({
      name: 'prolog',
      filename: 'prolog.js',
    }),

    // Copy the static files to the output directory:
    new CopyWebpackPlugin([
      { from: rscDir, to: outDir }
    ])
  ],

  resolve: {
    modules: [
      srcDir,
      modDir
    ]
  },

  devServer: {
    contentBase: outDir,
    outputPath: outDir,
    hot: true,
    inline: true,
    port: 8001,
    proxy: [{
      context: [
        '/api',
      ],
      target: 'http://localhost:8000',
    }]
  }
}
