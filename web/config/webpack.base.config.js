const webpack = require("webpack");
const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const ExtractTextPlugin = require('mini-css-extract-plugin');
const METADATA = require('./metadata.js');

module.exports = function (env) {
    return {
        entry: {
            app: './src/index.tsx'
        },

        output: {
            path: path.resolve(__dirname, '../dist'),
            filename: '[name].[chunkhash].js'
        },

        module: {
            rules: [
                {
                    enforce: 'pre',
                    test: /\.tsx?$/,
                    loader: 'tslint-loader',
                    options: {
                        fileOutput: {
                            // The directory where each file's report is saved
                            dir: './lint-reports/',

                            // The extension to use for each report's filename. Defaults to 'txt'
                            ext: 'xml',

                            // If true, all files are removed from the report directory at the beginning of run
                            clean: true,

                            // A string to include at the top of every report file.
                            // Useful for some report formats.
                            header: '<?xml version="1.0" encoding="utf-8"?>\n<checkstyle version="5.7">',

                            // A string to include at the bottom of every report file.
                            // Useful for some report formats.
                            footer: '</checkstyle>'
                        }
                    }
                },
                {
                    enforce: 'pre',
                    test: /\.js$/,
                    loader: 'source-map-loader'
                },
                {
                    enforce: 'pre',
                    test: /\.tsx?$/,
                    use: 'source-map-loader'
                },
                {
                    test: /\.tsx?$/,
                    loader: 'ts-loader',
                    options: {
                        transpileOnly: true
                    },
                    exclude: /node_modules/
                },
                {
                    test: /\.less$/,
                    use: [
                      {
                        loader: 'style-loader', // creates style nodes from JS strings
                      },
                      {
                        loader: 'css-loader', // translates CSS into CommonJS
                      },
                      {
                        loader: 'less-loader', // compiles Less to CSS
                      },
                    ],
                },
                {
                    test: /\.css$/,
                    use: {
                        loader: 'css-loader'
                    }
                },
                {
                    test: /\.(gif|png|jpg|jpeg|svg|ico)($|\?)/,
                    // embed images and fonts smaller than 5kb
                    loader: 'url-loader?limit=5000&hash=sha512&digest=hex&size=16&name=images/[name]-[hash].[ext]'
                },
                {
                    test: /\.(woff|woff2|eot|ttf)($|\?)/,
                    // embed images and fonts smaller than 5kb
                    loader: 'url-loader?limit=5000&hash=sha512&digest=hex&size=16&name=fonts/[name]-[hash].[ext]'
                }
            ]
        },

        resolve: {
            extensions: ['.tsx', '.ts', '.js', '.less'],
            modules: ['./', path.resolve(__dirname, 'src'), 'node_modules']
        },

        plugins: [
            // providing the lib dependencies so that they are present in the global scope
            new webpack.ProvidePlugin({
                React: 'react',
                ReactDOM: 'react-dom'
            }),

            // insert bundled script and metadata into index.html
            new HtmlWebpackPlugin({
                template: 'src/index.html',
                filename: 'index.html',
                metadata: METADATA,
                favicon: 'src/favicon.ico'
            }),

        ]
    }
};
