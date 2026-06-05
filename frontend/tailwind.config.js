/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './src/pages/**/*.{js,ts,jsx,tsx,mdx}',
    './src/components/**/*.{js,ts,jsx,tsx,mdx}',
    './src/app/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      colors: {
        // Planeta Azul brand palette
        base: {
          DEFAULT: '#080D12',
          50:  '#0D1620',
          100: '#111E2B',
          200: '#162536',
          300: '#1A2D42',
        },
        brand: {
          50:  '#003D5C',
          100: '#004F78',
          200: '#006296',
          300: '#0076B3',
          400: '#008BD0',
          500: '#009FEE',
          600: '#00BAFF',
          cyan: '#00D4FF',  // accent
        },
        status: {
          draft:       '#4B5563',
          pending:     '#D97706',
          in_progress: '#2563EB',
          waiting:     '#7C3AED',
          done:        '#059669',
          rejected:    '#DC2626',
        },
        priority: {
          low:    '#6B7280',
          medium: '#2563EB',
          high:   '#D97706',
          urgent: '#DC2626',
        },
      },
      fontFamily: {
        display: ['Syne', 'sans-serif'],
        body:    ['DM Sans', 'sans-serif'],
        mono:    ['JetBrains Mono', 'monospace'],
      },
      backgroundImage: {
        'grid-subtle': 'radial-gradient(circle, #162536 1px, transparent 1px)',
        'glow-cyan': 'radial-gradient(ellipse at center, rgba(0,212,255,0.15) 0%, transparent 70%)',
      },
      boxShadow: {
        'card':      '0 1px 3px rgba(0,0,0,0.4), inset 0 1px 0 rgba(255,255,255,0.05)',
        'accent':    '0 0 20px rgba(0,212,255,0.15)',
        'accent-lg': '0 0 40px rgba(0,212,255,0.25)',
      },
      animation: {
        'pulse-slow': 'pulse 3s cubic-bezier(0.4, 0, 0.6, 1) infinite',
        'spin-slow':  'spin 8s linear infinite',
        'fade-in':    'fadeIn 0.3s ease-out',
        'slide-up':   'slideUp 0.3s ease-out',
      },
      keyframes: {
        fadeIn:  { from: { opacity: '0' }, to: { opacity: '1' } },
        slideUp: { from: { opacity: '0', transform: 'translateY(8px)' }, to: { opacity: '1', transform: 'translateY(0)' } },
      },
    },
  },
  plugins: [],
}
