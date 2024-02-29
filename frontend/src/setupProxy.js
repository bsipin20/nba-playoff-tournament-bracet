const { createProxyMiddleware } = require('http-proxy-middleware');

module.exports = function (app) {
  app.use(
    '/v1', // Replace with your Go API path
    createProxyMiddleware({
      target: 'http://localhost:8080', // Your Go server URL
      changeOrigin: true,
    })
  );
};
