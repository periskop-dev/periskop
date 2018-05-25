const webpackMerge = require('webpack-merge');
const baseConfig = require('./webpack.base.config.js');
const METADATA = require('./metadata.js');

module.exports = function (env) {
    return webpackMerge(baseConfig(), {
        // chunk plugin doesn't work here
    })
};