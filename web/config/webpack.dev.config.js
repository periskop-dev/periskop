const webpack = require("webpack");
const webpackMerge = require('webpack-merge');
const baseConfig = require('./webpack.base.config.js');
const METADATA = require('./metadata.js');

module.exports = function (env) {
    return webpackMerge(baseConfig(), {
        devtool: 'inline-source-map',

        // Webpack Development Server config
        devServer: {
            port: METADATA.port,
            host: METADATA.host,
            historyApiFallback: true,
            watchOptions: {
                aggregateTimeout: 300,
                poll: 1000,
                ignored: /node_modules/
            }
        },

        plugins: [
            // separate all libs from node modules in a vendor file
            new webpack.optimize.CommonsChunkPlugin({
                name: 'vendor',
                minChunks: function (module) {
                    // this assumes your vendor imports exist in the node_modules directory
                    return module.context && module.context.indexOf('node_modules') !== -1;
                }
            })
        ]
    })
};