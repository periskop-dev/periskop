const METADATA = {
    title: 'Periskop',
    description: 'Exception Aggregator for micro-service environments',
    host: '0.0.0.0',
    port: 3000,
    backend_host: process.env.SERVER_URL,
    backend_port: process.env.SERVER_PORT,
};

module.exports = METADATA;
