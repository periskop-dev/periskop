function buildConfig(env) {
    return require('./config/webpack.' + env + '.config.js')(env)
}

module.exports = buildConfig;