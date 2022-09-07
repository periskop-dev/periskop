const METADATA = {
    title: 'Periskop',
    description: 'Exception Aggregator for micro-service environments',
    host: '0.0.0.0',
    port: 3000,
    api_host: process.env.API_HOST || 'localhost',
    api_port: process.env.API_PORT || 7777,
};

module.exports = METADATA;
