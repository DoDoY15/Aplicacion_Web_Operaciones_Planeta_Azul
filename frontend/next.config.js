/** @type {import('next').NextConfig} */
const { withNextOnPages } = require('@cloudflare/next-on-pages/next-config')

const nextConfig = {
  async rewrites() {
    if (process.env.NODE_ENV !== 'development') return []
    return [
      { source: '/api/:path*', destination: 'http://localhost:8080/api/:path*' },
      { source: '/auth/:path*', destination: 'http://localhost:8080/auth/:path*' },
    ]
  },
}

module.exports = withNextOnPages(nextConfig)
