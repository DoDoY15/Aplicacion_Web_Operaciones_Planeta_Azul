/** @type {import('next').NextConfig} */
const nextConfig = {
  async rewrites() {
    if (process.env.NODE_ENV !== 'development') return []
    return [
      { source: '/api/:path*', destination: 'http://localhost:8080/api/:path*' },
      { source: '/auth/:path*', destination: 'http://localhost:8080/auth/:path*' },
    ]
  },
}

module.exports = nextConfig
